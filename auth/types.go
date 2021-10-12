package auth

type UserRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
