package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/stretchr/testify/mock"
)

type BeneficiaryRepository struct {
	mock.Mock
}

func (m *BeneficiaryRepository) GetUserIDByMobileNumber(ctx context.Context, mobileCountryCode, mobileNumber string) (int, error) {
	args := m.Called(ctx, mobileCountryCode, mobileNumber)

	return args.Int(0), args.Error(1)
}

func (m *BeneficiaryRepository) CreateBeneficiary(ctx context.Context, userID int, beneficiaryID int) error {
	args := m.Called(ctx, userID, beneficiaryID)
	return args.Error(0)
}

func (m *BeneficiaryRepository) IsLinkedByUserIDAndBeneficiaryID(ctx context.Context, userID int, beneficiaryID int) error {
	args := m.Called(ctx, userID, beneficiaryID)
	return args.Error(0)
}

func (m *BeneficiaryRepository) UpdateBeneficiaryIsDeleted(ctx context.Context, userID int, beneficiaryID int, isDeleted int) error {
	args := m.Called(ctx, userID, beneficiaryID, isDeleted)
	return args.Error(0)
}

func (m *BeneficiaryRepository) GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*domain.Beneficiary, error) {
	args := m.Called(ctx, beneficiaryID, userID)

	var beneficiary *domain.Beneficiary
	if v, ok := args.Get(0).(*domain.Beneficiary); ok {
		beneficiary = v
	}

	return beneficiary, args.Error(1)
}

func (m *BeneficiaryRepository) GetBeneficiaries(ctx context.Context, userID int) (*[]domain.Beneficiary, error) {
	args := m.Called(ctx, userID)

	var beneficiaries *[]domain.Beneficiary
	if v, ok := args.Get(0).(*[]domain.Beneficiary); ok {
		beneficiaries = v
	}

	return beneficiaries, args.Error(1)
}
