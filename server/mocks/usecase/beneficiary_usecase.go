package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/stretchr/testify/mock"
)

type BeneficiaryUsecase struct {
	mock.Mock
}

type BeneficiaryUsecaseReturnValues struct {
	CreateBeneficiary []interface{}
	UpdateBeneficiary []interface{}
	GetBeneficiary    []interface{}
	GetBeneficiaries  []interface{}
}

func (m *BeneficiaryUsecase) CreateBeneficiary(ctx context.Context, userID int, req dto.CreateBeneficiaryRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *BeneficiaryUsecase) UpdateBeneficiary(ctx context.Context, userID int, req dto.UpdateBeneficiaryRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *BeneficiaryUsecase) GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*dto.GetBeneficiaryResponse, error) {
	args := m.Called(ctx, beneficiaryID, userID)

	var beneficiary *dto.GetBeneficiaryResponse
	if v, ok := args.Get(0).(*dto.GetBeneficiaryResponse); ok {
		beneficiary = v
	}

	return beneficiary, args.Error(1)
}

func (m *BeneficiaryUsecase) GetBeneficiaries(ctx context.Context, userID int) (*dto.GetBeneficiariesResponse, error) {
	args := m.Called(ctx, userID)

	var beneficiaries *dto.GetBeneficiariesResponse
	if v, ok := args.Get(0).(*dto.GetBeneficiariesResponse); ok {
		beneficiaries = v
	}

	return beneficiaries, args.Error(1)
}
