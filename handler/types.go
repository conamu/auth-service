package handler

type UserRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
