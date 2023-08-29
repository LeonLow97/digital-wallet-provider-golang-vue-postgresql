package users

import "context"

type Service interface {
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

func (s *service) GetUser(ctx context.Context, username string) (*GetUser, error) {
	user, err := s.repo.GetUserCurrencyAndBalanceByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
