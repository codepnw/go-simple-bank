package auth

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
