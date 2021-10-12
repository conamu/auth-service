package sender

type ISender interface {
	SendWelcome(username, email, subject string) error
	SendSignup(username, email, subject string) error
	SendPasswordReset(username, resetUrl, email, subject string) error
}

type Sender struct {
	BaseUrl           string
	WelcomeMail       string
	SignupMail        string
	PasswordResetMail string
}
