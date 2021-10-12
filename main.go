package main

import (
	"auth-service/sender"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type server struct {
	sender sender.ISender
}

func NewServer() *server {
	send, err := sender.NewSender()
	if err != nil {
		fmt.Println(err.Error())
	}
	return &server{
		sender: send,
	}
}

func main() {
	godotenv.Load()

	authServer := NewServer()

}
