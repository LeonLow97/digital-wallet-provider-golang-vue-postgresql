package users

import (
	"context"
	"log"
	"time"

	"github.com/LeonLow97/internal/utils"
	"github.com/golang-jwt/jwt"
)

var jwtTokenExpiry = time.Minute * 15
var refreshTokenExpiry = time.Hour * 24
var jwtSecretKey = "putthisinenvfile!"
var issuer = "mobilewallet"

type Service interface {
	Login(ctx context.Context, creds *Credentials) (*User, *Token, error)
	GetUser(ctx context.Context, username string) (*GetUser, error)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

func (s *service) Login(ctx context.Context, creds *Credentials) (*User, *Token, error) {
	// look up the user by username
	user, err := s.repo.GetByUsername(ctx, creds.Username)
	if err != nil {
		log.Println(err)
		return nil, nil, utils.UnauthorizedError{Message: "Incorrect username/password. Please try again."}
	}

	// validate the user's password
	validPassword, err := utils.PasswordMatchers(user.Password, creds.Password)
	if err != nil {
		return nil, nil, utils.InternalServerError{Message: err.Error()}
	}
	if !validPassword {
		log.Println("wrong password")
		return nil, nil, utils.UnauthorizedError{Message: "Incorrect username/password. Please try again."}
	}

	// ensure the user is active
	if user.Active == 0 {
		return nil, nil, utils.UnauthorizedError{Message: "This user account has been disabled. Please contact the system administrator."}
	}

	// we have a valid user, generate a JWT Token
	token, err := generateJwtAccessTokenAndRefreshToken(user, 15*time.Minute)

	return user, token, nil
}

func (s *service) GetUser(ctx context.Context, username string) (*GetUser, error) {
	user, err := s.repo.GetUserCurrencyAndBalanceByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

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
		return &Token{}, err
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
		return &Token{}, err
	}

	var tokens = Token{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return &tokens, nil
}
