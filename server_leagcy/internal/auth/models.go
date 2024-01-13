package auth

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Active   int    `json:"active"`
	Admin    int    `json:"admin"`
}

type Token struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
