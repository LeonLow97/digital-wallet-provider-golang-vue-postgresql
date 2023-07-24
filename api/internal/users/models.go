package users

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUser struct {
	Username     string  `json:"username"`
	MobileNumber string  `json:"mobile_number"`
	Currency     string `json:"currency"`
	Balance      float64  `json:"balance"`
}
