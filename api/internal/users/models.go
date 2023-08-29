package users

type GetUser struct {
	Username     string  `json:"username"`
	MobileNumber string  `json:"mobile_number"`
	Currency     string  `json:"currency"`
	Balance      float64 `json:"balance"`
}
