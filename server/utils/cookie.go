package utils

import (
	"net/http"
	"time"
)

func IssueCookie(writer http.ResponseWriter, token string) {
	cookieExpiration := time.Now().Add(14 * time.Minute).Second()

	cookie := &http.Cookie{
		Name:     JWT_COOKIE,
		Value:    token,
		MaxAge:   cookieExpiration,
		Path:     "/",
		Domain:   "localhost", // TODO: replace with config domain name
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(writer, cookie)
}
