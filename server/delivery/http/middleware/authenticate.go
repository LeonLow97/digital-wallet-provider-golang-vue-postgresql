package middleware

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/LeonLow97/go-clean-architecture/domain"
	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type AuthenticationMiddleware struct {
	cfg         infrastructure.Config
	skipperFunc SkipperFunc
	redisClient infrastructure.RedisClient
	authUsecase domain.UserUsecase
}

func NewAuthenticationMiddleware(cfg infrastructure.Config, skipperFunc SkipperFunc, redisClient infrastructure.RedisClient, authUsecase domain.UserUsecase) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		cfg:         cfg,
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

		ctx := r.Context()

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
			return []byte(m.cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			log.Println("JWT Token is invalid", err)
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
		ctx = utils.UserIDWithContext(ctx, userID)

		// Retrieve sessionID from claims
		sessionID, ok := claims["sessionID"].(string)
		if !ok {
			log.Println("Unable to retrieve session id from token claims")
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		ctx = utils.SessionIDWithContext(ctx, sessionID)

		// check if session exists in redis string and redis set.
		// If session exist, extend the session in redis. If session does not exist, unauthorized
		if _, err := m.redisClient.GetEx(ctx, sessionID, utils.SESSION_EXPIRY); err != nil {
			if errors.Is(err, redis.Nil) {
				// clean up stale sessionID from Redis Set for the specified userID
				// if sessionID has expired in string, it might still be present in Redis Set
				userIDString := strconv.Itoa(userID)
				_ = m.redisClient.SRem(ctx, userIDString, sessionID)

				// sessionID missing, user is unauthorized (session expired or invalid sessionID provided)
				log.Printf("failed to get key from redis for sessionID: %s and userID: %d\n", sessionID, userID)
				utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
				return
			} else {
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
