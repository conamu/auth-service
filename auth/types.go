package auth

type UserRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}
