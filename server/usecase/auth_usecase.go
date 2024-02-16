package usecase

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	redisClient    infrastructure.RedisClient
	userRepository domain.UserRepository
}

func NewAuthUsecase(userRepository domain.UserRepository, redisClient infrastructure.RedisClient) domain.UserUsecase {
	return &loginUsecase{
		redisClient:    redisClient,
		userRepository: userRepository,
	}
}

var (
	accessTokenExpiry  = time.Hour * 99999 // TODO: for development, reset to 15 minutes for PROD
	refreshTokenExpiry = time.Hour * 24
	jwtSecretKey       = os.Getenv("JWT_SECRET_KEY")
	issuer             = os.Getenv("API_DOMAIN")
)

func (uc *loginUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *dto.Token, error) {
	// retrieve user details from db
	user, err := uc.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}

	// authenticating user via password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, bcrypt.ErrHashTooShort):
		return nil, nil, exception.ErrInvalidCredentials
	case err != nil:
		return nil, nil, err
	}

	// checking if user is active
	if !user.Active {
		return nil, nil, exception.ErrInactiveUser
	}

	// generate session token with uuid
	sessionID := uuid.New().String()

	// generate access token and refresh token
	accessToken, err := generateJWTAccessToken(user, accessTokenExpiry, sessionID)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := generateJWTRefreshToken(user, refreshTokenExpiry, sessionID)
	if err != nil {
		return nil, nil, err
	}
	token := &dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// storing sessionID => userID mapping
	if err := uc.redisClient.Set(ctx, sessionID, user.ID); err != nil {
		return nil, nil, err
	}

	// storing userID => sessionID mapping in Redis Set to keep track of users with multiple devices logged on
	userIDString := strconv.Itoa(user.ID)
	if err := uc.redisClient.SAdd(ctx, userIDString, sessionID); err != nil {
		return nil, nil, err
	}

	// storing sessionID => sessionObject mapping
	// TODO: Store roles, permissions, emails in sessionObject
	// sessionObjectBytes, _ := json.Marshal(user)
	// if err := uc.redisClient.Set(ctx, sessionID, sessionObjectBytes); err != nil {
	// 	return nil, nil, err
	// }

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

func (uc *loginUsecase) RemoveSessionFromRedis(ctx context.Context, sessionID string) error {
	// retrieve userID from redis
	userID, err := uc.redisClient.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	// remove sessionID from redis
	if err := uc.redisClient.Del(ctx, sessionID); err != nil {
		return err
	}

	// remove sessionID from redis set, key is userID
	if err := uc.redisClient.SRem(ctx, userID, sessionID); err != nil {
		return err
	}

	return nil
}

// generateJWTAccessToken returns the JWT Access Token with the stores session ID
func generateJWTAccessToken(user *domain.User, ttl time.Duration, sessionID string) (string, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Username
	claims["sub"] = user.ID
	claims["aud"] = issuer // audience
	claims["iss"] = issuer // issuer (assigned to claims.Issuer)
	claims["admin"] = 0
	claims["sessionID"] = sessionID
	if user.Admin {
		claims["admin"] = true
	}

	// set token expiry
	claims["exp"] = time.Now().Add(ttl).Unix()

	// generate signed access token
	signedAccessToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

// generateJWTRefreshToken returns the JWT Refresh Token for retrieving subsequent fresh JWT Access Token
func generateJWTRefreshToken(user *domain.User, ttl time.Duration, sessionID string) (string, error) {
	// generate refresh token (users might not use) - less claims as compared to jwt token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = user.ID
	// set expiry, must be longer than access token
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// generate signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return signedRefreshToken, nil
}
