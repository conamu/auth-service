package auth

import (
	"auth-service/sender"
	"database/sql"
	"errors"
	"gitlab.ho-me.zone/conamu/base-tools/v2/hashing"
	"log"
	"time"
)

func RegisterUser(user *UserRequest, db *sql.DB, sender sender.ISender) error {

	hash, err := hashing.BcryptHash([]byte(user.Password))
	if err != nil {
		return err
	}

	query := `INSERT INTO USERS (USERNAME, PASSWORD, EMAIL, PERMISSION) VALUES (?,?,?,?);`
	_, err = db.Exec(query, user.User, hash, user.Email, "user")
	if err != nil {
		return err
	}
	err = sender.SendWelcome(user.User, user.Email, "Welcome")
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(user *UserRequest, db *sql.DB) (string, error) {
	var (
		username   string
		password   string
		email      string
		permission string
	)
	query := `SELECT * FROM USERS WHERE EMAIL=?`

	rows, err := db.Query(query, user.Email)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&username, &password, &email, &permission)
		if err != nil {
			return "", err
		}
	}
	log.Println("User found: " + username + email + permission)
	readUser := &UserRequest{
		User:     username,
		Password: password,
		Email:    email,
	}
	if readUser.Email != user.Email {
		return "", errors.New("email does not match to entry")
	}
	err = hashing.BcryptComparePassword([]byte(user.Password), []byte(readUser.Password))
	if err != nil {
		return "", err
	}

	pasetoGen, err := NewPasetoMaker("afik==hgb24sdfeoufcafik==hgb24sd")
	if err != nil {
		return "", err
	}

	token, err := pasetoGen.CreateToken(readUser.User, time.Hour*1)
	if err != nil {
		return "", err
	}
	log.Println("User logged in, generating PASETO Token")
	return token, nil
}

func ValidateToken(token string) error {
	pasetoChecker, err := NewPasetoMaker("afik==hgb24sdfeoufcafik==hgb24sd")
	if err != nil {
		return err
	}
	payload, err := pasetoChecker.VerifyToken(token)
	if err != nil {
		return err
	}
	log.Println("Token valid: " + payload.Username + " " + payload.ID.String())
	return nil
}

func EditPassword(db *sql.DB, sender sender.ISender) error {

	return nil
}

func ResetPassword(user *UserRequest, db *sql.DB, sender sender.ISender) error {

	return nil
}

func PerformPasswordReset(user *UserRequest, db *sql.DB, sender sender.ISender) error {
	return nil
}
