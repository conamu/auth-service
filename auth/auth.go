package auth

import (
	"auth-service/sender"
	"database/sql"
	"gitlab.ho-me.zone/conamu/base-tools/v2/hashing"
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

func LoginUser(user *UserRequest, db *sql.DB) error {

	return nil
}

func EditPassword(db *sql.DB, sender sender.ISender) error {

	return nil
}

func ResetPassword(user *UserRequest, db *sql.DB, sender sender.ISender) error {

	return nil
}
