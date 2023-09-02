package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LeonLow97/internal/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtTokenExpiry = time.Minute * 15
var refreshTokenExpiry = time.Hour * 24
var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
var issuer = os.Getenv("API_DOMAIN")

// generateToken gives a secure token and returns it with claims
func generateJwtAccessTokenAndRefreshToken(user *User, ttl time.Duration) (*Token, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Username
	claims["sub"] = user.ID
	claims["aud"] = issuer // audience
	claims["iss"] = issuer // issuer (assigned to claims.Issuer)
	claims["admin"] = 0
	if user.Admin == 1 {
		claims["admin"] = true
	}

	// set token expiry
	claims["exp"] = time.Now().Add(ttl).Unix()

	// generate signed token
	signedAccessToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return &Token{}, utils.InternalServerError{Message: err.Error()}
	}

	// generate refresh token (users might not use) - less claims as compared to jwt token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = user.ID
	// set expiry, must be longer than access token
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// generate signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return &Token{}, utils.InternalServerError{Message: err.Error()}
	}

	var tokens = Token{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return &tokens, nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func passwordMatchers(hashedPassword, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, utils.InternalServerError{Message: err.Error()}
		}
	}

	return true, nil
}

func ValidateToken(r *http.Request) (int, error) {
	// retrieve the authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return 0, utils.UnauthorizedError{Message: "no authorization header received"}
	}

	// retrieve the jwt token from header
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return 0, utils.UnauthorizedError{Message: "no valid authorization header received"}
	}

	// checking jwt token expiry
	jwtToken := headerParts[1]
	accessToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !accessToken.Valid {
		return 0, utils.UnauthorizedError{Message: err.Error()}
	}

	accessClaims := accessToken.Claims.(jwt.MapClaims)

	// Check access token expiration
	accessExp := time.Unix(int64(accessClaims["exp"].(float64)), 0)
	if time.Until(accessExp) > jwtTokenExpiry {
		return 0, utils.UnauthorizedError{Message: "token has expired"}
	}

	// TODO: check whether user is active

	// retrieve userId from token claims
	sub := accessClaims["sub"]

	// perform type assertion to convert interface{} to int
	floatValue, ok := sub.(float64)
	if !ok {
		return 0, utils.UnauthorizedError{Message: "userId is not of type float64"}
	}

	return int(floatValue), nil
}
