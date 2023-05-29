package auth

import (
	"auth-service/sender"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gitlab.ho-me.zone/conamu/base-tools/util"
	"gitlab.ho-me.zone/conamu/base-tools/v2/hashing"
)

func RegisterUser(user *UserRequest, db *sql.DB, sender sender.ISender) error {
	if err := ValidateUser(user); err != nil {
		return err
	}
	hash, err := hashing.BcryptHash([]byte(user.Password))
	if err != nil {
		return errors.Join(errors.New("error in hashing"), err)
	}

	query := `INSERT INTO USERS (USERNAME, PASSWORD, EMAIL, PERMISSION) VALUES (?,?,?,?);`
	_, err = db.Exec(query, user.User, hash, user.Email, user.Permission)
	if err != nil {
		return errors.Join(errors.New("error in database client"), err)
	}
	err = sender.SendWelcome(user.User, user.Email, "Welcome")
	if err != nil {
		return errors.Join(errors.New("error in mail sender"), err)
	}

	return nil
}

func LoginUser(user *UserRequest, db *sql.DB) (*LoginResponse, error) {
	var (
		username   string
		password   string
		email      string
		permission string
	)
	query := `SELECT * FROM USERS WHERE EMAIL=?`

	rows, err := db.Query(query, user.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&username, &password, &email, &permission)
		if err != nil {
			return nil, err
		}
	}
	log.Println("User found: " + username + email + permission)
	readUser := &UserRequest{
		User:     username,
		Password: password,
		Email:    email,
	}
	if readUser.Email != user.Email {
		return nil, errors.New("email does not match to entry")
	}
	err = hashing.BcryptComparePassword([]byte(user.Password), []byte(readUser.Password))
	if err != nil {
		return nil, err
	}

	pasetoGen, err := NewPasetoMaker("afik==hgb24sdfeoufcafik==hgb24sd")
	if err != nil {
		return nil, err
	}

	token, err := pasetoGen.CreateToken(readUser.User, permission, time.Hour*1)
	if err != nil {
		return nil, err
	}
	log.Println("User logged in, generating PASETO Token")
	return &LoginResponse{
		Token: token,
		User:  username,
		Role:  permission,
	}, nil
}

func ValidateToken(token string) (*Payload, error) {
	pasetoChecker, err := NewPasetoMaker("afik==hgb24sdfeoufcafik==hgb24sd")
	if err != nil {
		return nil, err
	}
	payload, err := pasetoChecker.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	log.Println("Token valid: " + payload.Username + " " + payload.ID.String())
	return payload, nil
}

func ResetPasswordRequest(pwrq *PasswordResetRequest, db *sql.DB, sender sender.ISender) (string, error) {
	var (
		user  string
		email string
	)
	getUserQuery := `SELECT USERNAME,EMAIL FROM USERS WHERE EMAIL=?;`
	createPwResetEntryQuery := `INSERT INTO PWRESETS (EMAIL,RESETID) VALUES (?,?)`
	rows, err := db.Query(getUserQuery, pwrq.Email)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		err = rows.Scan(&user, &email)
		if err != nil {
			return "", err
		}
	}
	resetId := strconv.Itoa(util.GetRandom(100000000000))
	_, err = db.Exec(createPwResetEntryQuery, email, resetId)
	if err != nil {
		return "", err
	}
	resetUrl := fmt.Sprintf("https://murat.karl-bock.academy/reset?resetId=%s", resetId)
	err = sender.SendPasswordReset(user, resetUrl, email, "Password Reset")
	if err != nil {
		return "", err
	}
	return resetId, nil
}

func PerformPasswordReset(pwReset *PasswordReset, db *sql.DB, sender sender.ISender) error {
	var (
		resetId string
		email   string
		user    string
	)
	getResetEntryQuery := `SELECT * FROM PWRESETS WHERE RESETID=?`
	getUserEntryQuery := `SELECT USERNAME FROM USERS WHERE EMAIL=?`
	deleteResetIDQuery := `DELETE FROM PWRESETS WHERE RESETID=?`
	setPasswordQuery := `UPDATE USERS SET PASSWORD=? WHERE EMAIL=?;`
	hashedPw, err := hashing.BcryptHash([]byte(pwReset.Password))
	if err != nil {
		return err
	}
	rows, err := db.Query(getResetEntryQuery, pwReset.ResetId)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&email, &resetId)
		if err != nil {
			return err
		}
	}
	_, err = db.Exec(deleteResetIDQuery, resetId)
	if err != nil {
		return err
	}
	rows, err = db.Query(getUserEntryQuery, email)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&user)
		if err != nil {
			return err
		}
	}
	log.Println("Update user password with email " + email)
	_, err = db.Exec(setPasswordQuery, hashedPw, email)
	if err != nil {
		return err
	}
	err = sender.SendPasswordWasReset(user, email, "Password Reset")
	if err != nil {
		return err
	}
	return nil
}

func ValidateUser(user *UserRequest) error {
	if user.User == "" {
		return errors.New("Username is empty!")
	}
	if user.Email == "" {
		return errors.New("Email is empty!")
	}
	if user.Password == "" {
		return errors.New("Password is empty!")
	}
	if user.Permission == "" {
		return errors.New("Role is empty!")
	}
	if user.Permission != "admin" && user.Permission != "user" {
		return errors.New("Wrong role specified on Signup!")
	}
	return nil
}
