package auth

type UserRequest struct {
	User       string `json:"user,omitempty"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Permission string `json:"role,omitempty"`
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
	User string `json:"user"`
	Role string `json:"role"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  string `json:"user"`
	Role  string `json:"role"`
}
