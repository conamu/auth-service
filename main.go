package main

import (
	"auth-service/sender"
	"fmt"
	"github.com/joho/godotenv"
)

type server struct {
	sender sender.ISender
}

func NewServer() *server {
	return &server{
		sender: sender.NewSender(),
	}
}

func main() {
	godotenv.Load()

	server := NewServer()

}
