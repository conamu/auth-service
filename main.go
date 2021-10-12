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
	to := "constantin.amundsen@gmail.com"
	subject := "Testing"

	server := NewServer()
	if err := server.sender.SendWelcome("Erdem", to, subject); err != nil {
		fmt.Println(err.Error())
	}

}
