package utils

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func RetrieveJWTClaimsUsername(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", UnauthorizedError{"authorization header is empty"}
	}

	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", UnauthorizedError{"invalid authorization header format"}
	}

	tokenString := parts[1]

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		return "", UnauthorizedError{"failed to parse or validate the token"}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", UnauthorizedError{"invalid token or claims"}
	}

	return claims.Username, nil
}
