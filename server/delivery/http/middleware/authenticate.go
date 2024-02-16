package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/golang-jwt/jwt/v4"
)

var (
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the jwt token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			log.Println("Missing Authorization in request header")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// split the header on spaces
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			log.Println("Authorization header does not contain 2 parts")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// check to see if we have the word "Bearer"
		if headerParts[0] != "Bearer" {
			log.Println("Bearer is missing in the authentication header")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		jwtTokenString := headerParts[1]

		// validate the JWT Token
		token, err := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(JWT_SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			log.Println("Token is invalid", err)
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// retrieve claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Unable to retrieve token claims")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Retrieve user id from claims
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			log.Println("Unable to retrieve user id from token claims")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		// set user id in context
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)

		// Retrieve sessionID from claims
		sessionID, ok := claims["sessionID"].(string)
		if !ok {
			log.Println("Unable to retrieve session id from token claims")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, utils.SessionIDKey, sessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
