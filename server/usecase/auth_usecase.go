package usecase

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	userRepository domain.UserRepository
}

func NewAuthUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &loginUsecase{
		userRepository: userRepository,
	}
}

func (uc *loginUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *dto.Token, error) {
	user, err := uc.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, bcrypt.ErrHashTooShort):
		return nil, nil, exception.ErrInvalidCredentials
	case err != nil:
		return nil, nil, err
	}

	if !user.Active {
		return nil, nil, exception.ErrInactiveUser
	}

	token, err := generateJwtAccessTokenAndRefreshToken(user, jwtTokenExpiry)
	if err != nil {
		return nil, nil, err
	}

	resp := dto.LoginResponse{
		Email:    user.Email,
		Username: user.Username,
	}

	return &resp, token, nil
}

func (uc *loginUsecase) SignUp(ctx context.Context, req dto.SignUpRequest) error {
	user, err := uc.userRepository.GetUserByEmailOrMobileNumber(ctx, req.Email, req.MobileNumber)
	if err != nil {
		if err != exception.ErrUserNotFound {
			return err
		}
	}

	// user already exist
	if user != nil {
		return exception.ErrUserFound
	}

	if !utils.IsValidPassword(req.Password) {
		return exception.ErrInvalidPassword
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}
	req.Password = string(hashedPasswordBytes)

	insertUser := domain.User{
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		MobileNumber: req.MobileNumber,
	}

	if req.FirstName != nil {
		insertUser.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		insertUser.LastName = *req.LastName
	}

	// create one user
	if err = uc.userRepository.InsertUser(ctx, &insertUser); err != nil {
		return err
	}

	return nil
}

var (
	jwtTokenExpiry     = time.Minute * 20
	refreshTokenExpiry = time.Hour * 24
	jwtSecretKey       = os.Getenv("JWT_SECRET_KEY")
	issuer             = os.Getenv("API_DOMAIN")
)

// generateToken gives a secure token and returns it with claims
func generateJwtAccessTokenAndRefreshToken(user *domain.User, ttl time.Duration) (*dto.Token, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Username
	claims["sub"] = user.ID
	claims["aud"] = issuer // audience
	claims["iss"] = issuer // issuer (assigned to claims.Issuer)
	claims["admin"] = 0
	if user.Admin {
		claims["admin"] = true
	}

	// set token expiry
	claims["exp"] = time.Now().Add(ttl).Unix()

	// generate signed token
	signedAccessToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	var tokens = dto.Token{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return &tokens, nil
}
