package main

import (
	"auth-service/handler"
	"auth-service/sender"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var userDbInitQuery = `CREATE TABLE IF NOT EXISTS USERS (USERNAME VARCHAR(30) NOT NULL, PASSWORD VARCHAR(120) NOT NULL, EMAIL VARCHAR(30) NOT NULL PRIMARY KEY, PERMISSION VARCHAR(10) NOT NULL);`
var passwordResetDBInitQuery = `CREATE TABLE IF NOT EXISTS PWRESETS(EMAIL VARCHAR(30) NOT NULL, RESETID VARCHAR(30) NOT NULL PRIMARY KEY);`

func main() {
	log.Println("Waiting for DB to be up...")
	time.Sleep(time.Second * 3)
	godotenv.Load()

	dbConfig := mysql.Config{
		User:   "kb-auth",
		Passwd: "kb-auth",
		Net:    "tcp",
		Addr:   "auth-service-db:3306",
		DBName: "auth",
	}

	router := mux.NewRouter()

	router.Use(corsMiddleware)

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	_, err = db.Exec(userDbInitQuery)
	_, err = db.Exec(passwordResetDBInitQuery)
	if err != nil {
		fmt.Println(err.Error())
	}

	mailSender, err := sender.NewSender()
	if err != nil {
		log.Println("Error creating sender: " + err.Error())
	}

	// User Registration - Create user in DB
	router.HandleFunc("/register", checkAuthHeader(handler.SignUpHandlerFunc(db, mailSender)))
	// User Authentication - Create PASETO Token and send back if email and password match
	router.HandleFunc("/login", checkAuthHeader(handler.LogInHandlerFunc(db)))
	// Token validation for frontend Return 200 OK if PASETO valid
	router.HandleFunc("/validate", checkAuthHeader(handler.ValidateHandlerFunc()))
	// Submit a password reset for an email
	router.HandleFunc("/resetpassword", checkAuthHeader(handler.ResetPasswordFunc(db, mailSender)))
	// Password reset execution
	router.HandleFunc("/reset", checkAuthHeader(handler.PerformPasswordResetFunc(db, mailSender)))
	// Ping - sends 200 OK
	router.HandleFunc("/ping", handler.PingHandlerFunc())

	log.Println("KB-Auth-Service listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func checkAuthHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			next(res, req)
			return
		}
		header := req.Header.Get("X-KBU-Auth")
		if header != "abcdefghijklmnopqrstuvwxyz" {
			res.WriteHeader(403)
			log.Println("Auth header does not match!")
			return
		}
		next(res, req)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Access-Control-Allow-Headers",
				"content-type,x-kbu-auth,content-length,x-kbu-login")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
