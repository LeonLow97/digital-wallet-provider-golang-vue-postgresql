package auth

import (
	"errors"
	"time"

	"github.com/LeonLow97/internal/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtTokenExpiry = time.Minute * 15
var refreshTokenExpiry = time.Hour * 24
var jwtSecretKey = "putthisinenvfile!"
var issuer = "mobilewallet"

// generateToken gives a secure token of exactly 26 characters in length and returns it
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
	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()

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
