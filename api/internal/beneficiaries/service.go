package beneficiaries

import (
	"context"

	"github.com/LeonLow97/internal/utils"
)

type Service interface {
	GetBeneficiaries(ctx context.Context, username string) (*Beneficiaries, error)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

func (s *service) GetBeneficiaries(ctx context.Context, username string) (*Beneficiaries, error) {
	beneficiaries, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, utils.ServiceError{Message: "User does not have any beneficiaries."}
	}

	return beneficiaries, nil

}
