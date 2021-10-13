package main

import (
	"auth-service/handler"
	"auth-service/sender"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var userDbInitQuery = `CREATE TABLE IF NOT EXISTS USERS (USERNAME VARCHAR(30) NOT NULL, PASSWORD VARCHAR(120) NOT NULL, EMAIL VARCHAR(30) NOT NULL PRIMARY KEY, PERMISSION VARCHAR(10) NOT NULL);`
var passwordResetDBInitQuery = `CREATE TABLE IF NOT EXISTS PWRESETS(EMAIL VARCHAR(30) NOT NULL, RESETID VARCHAR(30) NOT NULL PRIMARY KEY);`

type Server struct {
	sender sender.ISender
	router *mux.Router
	db     *sql.DB
}

func NewServer(sender sender.ISender, router *mux.Router, db *sql.DB) *Server {
	return &Server{
		sender: sender,
		router: router,
		db:     db,
	}
}

func main() {
	godotenv.Load()
	send, err := sender.NewSender()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(err)
	}

	dbConfig := mysql.Config{
		User:   "kb-auth",
		Passwd: "kb-auth",
		Net:    "tcp",
		Addr:   "auth-service-db:3306",
		DBName: "auth",
	}

	router := mux.NewRouter()
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	_, err = db.Exec(userDbInitQuery)
	_, err = db.Exec(passwordResetDBInitQuery)
	if err != nil {
		fmt.Println(err.Error())
	}
	authServer := NewServer(send, router, db)

	// User Registration - Create user in DB
	authServer.router.HandleFunc("/register", handler.SignUpHandlerFunc(authServer.db, authServer.sender))
	// User Authentication - Create PASETO Token and send back if email and password match
	authServer.router.HandleFunc("/login", handler.LogInHandlerFunc(authServer.db))
	// Token validation for frontend Return 200 OK if PASETO valid
	authServer.router.HandleFunc("/validate", handler.ValidateHandlerFunc())
	// Submit a password reset for an email
	authServer.router.HandleFunc("/resetpassword", handler.ResetPasswordFunc(authServer.db, authServer.sender))
	// Password reset execution
	authServer.router.HandleFunc("/reset", handler.PerformPasswordResetFunc(authServer.db, authServer.sender))
	// Ping - sends 200 OK
	authServer.router.HandleFunc("/ping", handler.PingHandlerFunc())

	log.Println("KB-Auth-Service listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", authServer.router))
	defer db.Close()
}
