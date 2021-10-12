package sender

type ISender interface {
	SendWelcome(username, email, subject string) error
	SendSignup(username string) error
	SendPasswordReset(username, resetUrl string) error
}

type Sender struct {
	BaseUrl           string
	WelcomeMail       string
	SignupMail        string
	PasswordResetMail string
}
