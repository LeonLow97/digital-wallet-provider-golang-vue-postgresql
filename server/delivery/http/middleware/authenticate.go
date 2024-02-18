package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/LeonLow97/go-clean-architecture/domain"
	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

var (
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
)

type AuthenticationMiddleware struct {
	skipperFunc SkipperFunc
	redisClient infrastructure.RedisClient
	authUsecase domain.UserUsecase
}

func NewAuthenticationMiddleware(skipperFunc SkipperFunc, redisClient infrastructure.RedisClient, authUsecase domain.UserUsecase) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		skipperFunc: skipperFunc,
		redisClient: redisClient,
		authUsecase: authUsecase,
	}
}

func (m AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.skipperFunc != nil && m.skipperFunc(r) {
			next.ServeHTTP(w, r)
			return
		}

		// // get the jwt token from the Authorization header
		// authHeader := r.Header.Get("Authorization")
		// if len(authHeader) == 0 {
		// 	log.Println("Missing Authorization in request header")
		// 	utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		// 	return
		// }

		// // split the header on spaces
		// headerParts := strings.Split(authHeader, " ")
		// if len(headerParts) != 2 {
		// 	log.Println("Authorization header does not contain 2 parts")
		// 	utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		// 	return
		// }

		// // check to see if we have the word "Bearer"
		// if headerParts[0] != "Bearer" {
		// 	log.Println("Bearer is missing in the authentication header")
		// 	utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		// 	return
		// }

		// jwtTokenString := headerParts[1]

		cookie, err := r.Cookie("mw-token")
		if err != nil {
			log.Println("Missing token cookie:", err)
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		jwtTokenString := cookie.Value

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

		// check if session exists in redis string and redis set.
		// If session exist, extend the session in redis. If session does not exist, unauthorized
		if _, err := m.redisClient.GetEx(ctx, sessionID, utils.SESSION_EXPIRY); err != nil {
			switch {
			case err == redis.Nil:
				// clean up stale sessionID from Redis Set for the specified userID
				// if sessionID has expired in string, it might still be present in Redis Set
				userIDString := strconv.Itoa(userID)
				_ = m.redisClient.SRem(ctx, userIDString, sessionID)

				// sessionID missing, user is unauthorized (session expired or invalid sessionID provided)
				log.Printf("failed to get key from redis for sessionID: %s and userID: %d\n", sessionID, userID)
				utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
				return
			case err != nil:
				log.Println("failed to get key from redis in authentication middleware", err)
				utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		// issue new JWT Token with the same session id
		jwtToken, err := m.authUsecase.GenerateJWTAccessToken(userID, utils.SESSION_EXPIRY, sessionID)
		if err != nil {
			log.Println("failed to reissue jwt access token", err)
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		// set HTTP Cookie with new JWT Token
		utils.IssueCookie(w, jwtToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
