package auth

import (
	"auth-service/sender"
	"database/sql"
)

func RegisterUser(user *UserRequest, db *sql.DB, sender sender.ISender) error {

	return nil
}

func LoginUser(user *UserRequest, db *sql.DB, sender sender.ISender) error {

	return nil
}

func EditPassword(db *sql.DB, sender sender.ISender) error {

	return nil
}

func ResetPassword(user *UserRequest, db *sql.DB, sender sender.ISender) error {

	return nil
}
