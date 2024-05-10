package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	cfg            infrastructure.Config
	redisClient    infrastructure.RedisClient
	smtpClient     infrastructure.SMTPClient
	userRepository domain.UserRepository
}

func NewUserUsecase(cfg infrastructure.Config, userRepository domain.UserRepository, redisClient infrastructure.RedisClient, smtpClient infrastructure.SMTPClient) domain.UserUsecase {
	return &loginUsecase{
		cfg:            cfg,
		redisClient:    redisClient,
		smtpClient:     smtpClient,
		userRepository: userRepository,
	}
}

var (
	accessTokenExpiry  = time.Minute * 15
	refreshTokenExpiry = time.Hour * 24
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
	accessToken, err := uc.GenerateJWTAccessToken(user.ID, accessTokenExpiry, sessionID)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := uc.generateJWTRefreshToken(user.ID, refreshTokenExpiry, sessionID)
	if err != nil {
		return nil, nil, err
	}
	token := &dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// storing sessionID => userID mapping
	if err := uc.redisClient.SetEx(ctx, sessionID, user.ID, utils.SESSION_EXPIRY); err != nil {
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
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
	}

	return &resp, token, nil
}

func (uc *loginUsecase) SignUp(ctx context.Context, req dto.SignUpRequest) error {
	user, err := uc.userRepository.GetUserByEmailOrMobileNumber(ctx, req.Email, req.MobileNumber)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		log.Println("failed to get user by email or mobile number")
		return err
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

func (uc *loginUsecase) ChangePassword(ctx context.Context, userID int, req dto.ChangePasswordRequest) error {
	// get user by user id
	user, err := uc.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("failed to get user by user id %d with error %v\n", userID, err)
		return err
	}

	// ensure current password is same as db password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		log.Printf("current password is not the same as the password in db for user id %d\n", userID)
		return exception.ErrInvalidCredentials
	}

	// ensure new password is different from current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword))
	if err == nil {
		log.Printf("new password is the same as the current password for user id %d\n", userID)
		return exception.ErrSamePassword
	}

	// hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		log.Printf("failed to hash new password with error %v\n", err)
		return err
	}

	// set user password to new password
	if err := uc.userRepository.ChangePassword(ctx, string(hashedNewPassword), userID); err != nil {
		log.Printf("failed to update user password for user id %d with error %v\n", userID, err)
		return err
	}

	return nil
}

func (uc *loginUsecase) SendPasswordResetEmail(ctx context.Context, req dto.SendPasswordResetEmailRequest) error {
	// check if user is valid by email
	user, err := uc.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("failed to retrieve user %s with error: %v\n", req.Email, err)
		return err
	}

	// generate authentication token
	authToken, err := utils.GenerateAuthenticationToken(32)
	if err != nil {
		log.Println("error generating authentication token")
		return err
	}

	// construct reset url with authentication token
	resetUrl := fmt.Sprintf("http://localhost:3000/#/password-reset/%s", authToken)
	emailBody := strings.Replace(utils.PasswordResetEmailBody, "{{.ResetURL}}", resetUrl, 1)

	// using 2 goroutines to send email and store authentication token in redis
	// Channels to receive errors from go routines
	emailErrChan := make(chan error)
	redisErrChan := make(chan error)

	// start a new go routine to send email
	go func() {
		err := uc.smtpClient.SendEmail(ctx, "digital-wallet@email.com", "Digital Wallet", []string{req.Email}, "Reset Your Password", emailBody)
		emailErrChan <- err
	}()

	// start a new go routine to store the token in redis
	go func() {
		key := fmt.Sprintf("password-reset-%s", authToken)

		values := map[string]interface{}{
			"email": user.Email,
			"id":    user.ID,
		}

		if err := uc.redisClient.HSet(ctx, key, values); err != nil {
			redisErrChan <- err
		}

		// set expiration time for the hash table
		if err := uc.redisClient.Expire(ctx, key, utils.PASSWORD_RESET_AUTH_TOKEN_EXPIRY); err != nil {
			redisErrChan <- err
		}
	}()

	// wait for both go routines to complete
	var emailErr, redisErr error
	select {
	case emailErr = <-emailErrChan:
	case redisErr = <-redisErrChan:
	}

	if emailErr != nil {
		log.Println("failed to send email", emailErr)
		return emailErr
	}
	if redisErr != nil {
		log.Println("failed to store authentication token in redis", redisErr)
		return redisErr
	}

	return nil
}

func (uc loginUsecase) PasswordReset(ctx context.Context, req dto.PasswordResetRequest) error {
	// retrieve user email by auth token
	key := fmt.Sprintf("password-reset-%s", req.Token)
	values, err := uc.redisClient.HGetAll(ctx, key)
	if err != nil {
		log.Println("failed to get user email in redis client with error:", err)
		return err
	}
	email := values["email"]

	// retrieve user details from db by email
	user, err := uc.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("failed to get user by email in db with error:", err)
		return err
	}

	// check if new password is same as the previous password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
		return exception.ErrSamePassword
	}

	// update user password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Println("failed to hash new password with error:", err)
		return err
	}

	if err := uc.userRepository.ChangePassword(ctx, string(hashedNewPassword), user.ID); err != nil {
		log.Printf("failed to update user password for user id %d with error %v\n", user.ID, err)
		return err
	}

	// delete the authentication token key in redis
	if err := uc.redisClient.Del(ctx, req.Token); err != nil {
		log.Println("failed to delete auth token key in redis with error:", err)
		return err
	}

	return nil
}

func (uc *loginUsecase) UpdateUser(ctx context.Context, userID int, req dto.UpdateUserRequest) error {
	updatedUser := domain.User{
		ID:           userID,
		Username:     req.Username,
		MobileNumber: req.MobileNumber,
		Email:        req.Email,
	}

	if req.FirstName != nil {
		updatedUser.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		updatedUser.LastName = *req.LastName
	}

	// update one user
	if err := uc.userRepository.UpdateUser(ctx, &updatedUser); err != nil {
		log.Printf("error updating one user with user id: %d\n", userID)
		return err
	}

	return nil
}

func (uc *loginUsecase) ExtendUserSessionInRedis(ctx context.Context, sessionID string, sessionExpiryInMinutes time.Duration) error {
	if err := uc.redisClient.Expire(ctx, sessionID, sessionExpiryInMinutes); err != nil {
		log.Println("error extending user session")
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
func (uc *loginUsecase) GenerateJWTAccessToken(userID int, ttl time.Duration, sessionID string) (string, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["aud"] = uc.cfg.JWT.Issuer // audience
	claims["iss"] = uc.cfg.JWT.Issuer // issuer (assigned to claims.Issuer)
	claims["admin"] = 0
	claims["sessionID"] = sessionID

	// TODO: Add to redis sessionID => sessionObject
	// if user.Admin {
	// 	claims["admin"] = true
	// }

	// set token expiry
	claims["exp"] = time.Now().Add(ttl).Unix()

	// generate signed access token
	signedAccessToken, err := token.SignedString([]byte(uc.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

// generateJWTRefreshToken returns the JWT Refresh Token for retrieving subsequent fresh JWT Access Token
func (uc *loginUsecase) generateJWTRefreshToken(userID int, ttl time.Duration, sessionID string) (string, error) {
	// generate refresh token (users might not use) - less claims as compared to jwt token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = userID
	// set expiry, must be longer than access token
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// generate signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(uc.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return signedRefreshToken, nil
}
