package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Beneficiary struct {
	BeneficiaryID           int    `db:"beneficiary_id"`
	IsDeleted               int    `db:"is_deleted"`
	BeneficiaryFirstName    string `db:"first_name"`
	BeneficiaryLastName     string `db:"last_name"`
	BeneficiaryEmail        string `db:"email"`
	BeneficiaryUsername     string `db:"username"`
	IsActive                int    `db:"active"`
	BeneficiaryMobileNumber string `db:"mobile_number"`
}

type BeneficiaryUsecase interface {
	CreateBeneficiary(ctx context.Context, req dto.CreateBeneficiaryRequest) error
	UpdateBeneficiary(ctx context.Context, req dto.UpdateBeneficiaryRequest) error
	GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*dto.GetBeneficiaryResponse, error)
	GetBeneficiaries(ctx context.Context, userID int) (*dto.GetBeneficiariesResponse, error)
}

type BeneficiaryRepository interface {
	GetUserIDByMobileNumber(ctx context.Context, mobileNumber string) (int, error)
	CreateBeneficiary(ctx context.Context, userID int, beneficiaryID int) error

	IsLinkedByUserIDAndBeneficiaryID(ctx context.Context, userID int, beneficiaryID int) error
	UpdateBeneficiaryIsDeleted(ctx context.Context, userID int, beneficiaryID int, isDeleted int) error

	GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*Beneficiary, error)
	GetBeneficiaries(ctx context.Context, userID int) (*[]Beneficiary, error)
}
