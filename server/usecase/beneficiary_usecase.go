package usecase

import (
	"context"
	"log"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
)

type beneficiaryUsecase struct {
	beneficiaryRepository domain.BeneficiaryRepository
}

func NewBeneficiaryUsecase(beneficiaryRepository domain.BeneficiaryRepository) domain.BeneficiaryUsecase {
	return &beneficiaryUsecase{
		beneficiaryRepository: beneficiaryRepository,
	}
}

func (uc *beneficiaryUsecase) CreateBeneficiary(ctx context.Context, req dto.CreateBeneficiaryRequest) error {
	// retrieve beneficiary id by mobile number
	beneficiaryID, err := uc.beneficiaryRepository.GetUserIDByMobileNumber(ctx, req.MobileNumber)
	if err != nil {
		log.Println("failed to get user id by mobile number", err)
		return err
	}

	if beneficiaryID == req.UserID {
		return exception.ErrUserIDEqualBeneficiaryID
	}

	// link user to beneficiary
	err = uc.beneficiaryRepository.CreateBeneficiary(ctx, req.UserID, beneficiaryID)
	if err != nil {
		log.Println("failed to create beneficiary", err)
		return err
	}

	return nil
}

func (uc *beneficiaryUsecase) UpdateBeneficiary(ctx context.Context, req dto.UpdateBeneficiaryRequest) error {
	// check if beneficiary id equal to user id
	if req.BeneficiaryID == req.UserID {
		return exception.ErrUserIDEqualBeneficiaryID
	}

	// check if beneficiary id is related to user id
	err := uc.beneficiaryRepository.IsLinkedByUserIDAndBeneficiaryID(ctx, req.UserID, req.BeneficiaryID)
	if err != nil {
		return err
	}

	// update beneficiary soft delete feature flag (is_deleted)
	err = uc.beneficiaryRepository.UpdateBeneficiaryIsDeleted(ctx, req.UserID, req.BeneficiaryID, req.IsDeleted)
	if err != nil {
		return err
	}

	return nil
}
