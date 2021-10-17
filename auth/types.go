package auth

type UserRequest struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Permission string `json:"permission"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type PasswordReset struct {
	Password string `json:"password"`
	ResetId  string `json:"resetId"`
}

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetResponse struct {
	ResetId string `json:"resetId"`
}

type ValidationResponse struct {
	Role string `json:"role"`
}
