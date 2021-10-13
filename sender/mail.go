package sender

import (
	"fmt"
	"github.com/containrrr/shoutrrr"
	"io/ioutil"
	"os"
)

func NewSender() (ISender, error) {
	url := fmt.Sprintf("smtp://%s:%s@%s:%s/?from=%s&to=%s&fromName=%s&title=%s&usehtml=%s",
		os.Getenv("username"),
		os.Getenv("password"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("from"),
		"%s",
		os.Getenv("fromName"),
		"%s",
		os.Getenv("useHtml"))

	welcomeMailfile, err := ioutil.ReadFile("./mail-templates/welcome-email.html")
	if err != nil {
		return nil, err
	}
	signupMailfile, err := ioutil.ReadFile("./mail-templates/signup-email.html")
	if err != nil {
		return nil, err
	}
	passwordMailfile, err := ioutil.ReadFile("./mail-templates/password-reset-email.html")
	if err != nil {
		return nil, err
	}
	passwordWasRestMailfile, err := ioutil.ReadFile("./mail-templates/password-was-reset-email.html")
	if err != nil {
		return nil, err
	}

	return &Sender{
		BaseUrl:              url,
		WelcomeMail:          string(welcomeMailfile),
		SignupMail:           string(signupMailfile),
		PasswordResetMail:    string(passwordMailfile),
		PasswordWasResetMail: string(passwordWasRestMailfile),
	}, nil
}

func (s *Sender) SendWelcome(username, email, subject string) error {
	url := fmt.Sprintf(s.BaseUrl, email, subject)
	message := fmt.Sprintf(s.WelcomeMail, username)
	return send(message, url)
}
func (s *Sender) SendSignup(username, email, subject string) error {
	url := fmt.Sprintf(s.BaseUrl, email, subject)
	message := fmt.Sprintf(s.WelcomeMail, username)
	return send(message, url)
}
func (s *Sender) SendPasswordReset(username, resetUrl, email, subject string) error {
	url := fmt.Sprintf(s.BaseUrl, email, subject)
	message := fmt.Sprintf(s.PasswordResetMail, username, resetUrl)
	return send(message, url)
}
func (s *Sender) SendPasswordWasReset(username, email, subject string) error {
	url := fmt.Sprintf(s.BaseUrl, email, subject)
	message := fmt.Sprintf(s.PasswordWasResetMail, username)
	return send(message, url)
}

func send(message, url string) error {
	err := shoutrrr.Send(url, message)
	if err != nil {
		return err
	}
	return nil
}
