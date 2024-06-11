package utils

import (
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/utils/constants"
)

func IssueCookie(writer http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     constants.JWT_COOKIE,
		Value:    token,
		MaxAge:   3600,
		Path:     "/",
		Domain:   "localhost", // TODO: replace with config domain name
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(writer, cookie)
}
