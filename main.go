package main

import (
	"auth-service/handler"
	"auth-service/sender"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

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
	}
	router := mux.NewRouter()
	db, err := sql.Open("mysql", "kb-auth:kb-auth@localhost/auth")
	if err != nil {
		fmt.Println(err.Error())
	}
	authServer := NewServer(send, router, db)

	authServer.router.HandleFunc("/register", handler.SignUpHandlerFunc(authServer.db, authServer.sender))
	authServer.router.HandleFunc("/login", handler.LogInHandlerFunc(authServer.db, authServer.sender))
	authServer.router.HandleFunc("/edituser", handler.EditUserHandlerFunc(authServer.db, authServer.sender))
	authServer.router.HandleFunc("/ping", handler.PingHandlerFunc())

	log.Fatal(http.ListenAndServe(":8080", authServer.router))
}
