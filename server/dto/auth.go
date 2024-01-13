package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Active   bool   `json:"active"`
	Admin    bool   `json:"admin"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
