package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Beneficiary struct {
	UserID        int `db:"user_id"`
	BeneficiaryID int `db:"beneficiary_id"`
}

type BeneficiaryUsecase interface {
	CreateBeneficiary(ctx context.Context, req dto.CreateBeneficiaryRequest) error
	UpdateBeneficiary(ctx context.Context, req dto.UpdateBeneficiaryRequest) error
}

type BeneficiaryRepository interface {
	GetUserIDByMobileNumber(ctx context.Context, mobileNumber string) (int, error)
	CreateBeneficiary(ctx context.Context, userID int, beneficiaryID int) error

	IsLinkedByUserIDAndBeneficiaryID(ctx context.Context, userID int, beneficiaryID int) error
	UpdateBeneficiaryIsDeleted(ctx context.Context, userID int, beneficiaryID int, isDeleted int) error
}
