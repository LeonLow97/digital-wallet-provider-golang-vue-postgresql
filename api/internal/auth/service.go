package auth

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/utils"
)

type Service interface {
	Login(ctx context.Context, creds *Credentials) (*User, *Token, error)
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
		return nil, nil, err
	}

	// validate the user's password
	validPassword, err := passwordMatchers(user.Password, creds.Password)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	if !validPassword {
		return nil, nil, utils.UnauthorizedError{Message: "Incorrect username/password. Please try again."}
	}

	// ensure the user is active
	if user.Active == 0 {
		return nil, nil, utils.UnauthorizedError{Message: "This user account has been disabled. Please contact the system administrator."}
	}

	// we have a valid user, generate a JWT Token
	token, err := generateJwtAccessTokenAndRefreshToken(user, jwtTokenExpiry)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
