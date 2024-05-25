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

func (uc *beneficiaryUsecase) CreateBeneficiary(ctx context.Context, userID int, req dto.CreateBeneficiaryRequest) error {
	// retrieve beneficiary id by mobile number
	beneficiaryID, err := uc.beneficiaryRepository.GetUserIDByMobileNumber(ctx, req.MobileCountryCode, req.MobileNumber)
	if err != nil {
		log.Println("failed to get user id by mobile number with error:", err)
		return err
	}

	if beneficiaryID == userID {
		return exception.ErrUserIDEqualBeneficiaryID
	}

	// link user to beneficiary
	err = uc.beneficiaryRepository.CreateBeneficiary(ctx, userID, beneficiaryID)
	if err != nil {
		log.Println("failed to create beneficiary with error:", err)
		return err
	}

	return nil
}

func (uc *beneficiaryUsecase) UpdateBeneficiary(ctx context.Context, userID int, req dto.UpdateBeneficiaryRequest) error {
	// check if beneficiary id equal to user id
	if req.BeneficiaryID == userID {
		return exception.ErrUserIDEqualBeneficiaryID
	}

	// check if beneficiary id is related to user id
	err := uc.beneficiaryRepository.IsLinkedByUserIDAndBeneficiaryID(ctx, userID, req.BeneficiaryID)
	if err != nil {
		return err
	}

	// update beneficiary soft delete feature flag (is_deleted)
	err = uc.beneficiaryRepository.UpdateBeneficiaryIsDeleted(ctx, userID, req.BeneficiaryID, req.IsDeleted)
	if err != nil {
		return err
	}

	return nil
}

func (uc *beneficiaryUsecase) GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*dto.GetBeneficiaryResponse, error) {
	// check if beneficiary id equal to user id
	if beneficiaryID == userID {
		return nil, exception.ErrUserIDEqualBeneficiaryID
	}

	// get one beneficiary
	beneficiary, err := uc.beneficiaryRepository.GetBeneficiary(ctx, beneficiaryID, userID)
	if err != nil {
		return nil, err
	}

	resp := &dto.GetBeneficiaryResponse{
		BeneficiaryID:                beneficiary.BeneficiaryID,
		IsDeleted:                    beneficiary.IsDeleted,
		BeneficiaryFirstName:         beneficiary.BeneficiaryFirstName,
		BeneficiaryLastName:          beneficiary.BeneficiaryLastName,
		BeneficiaryEmail:             beneficiary.BeneficiaryEmail,
		BeneficiaryUsername:          beneficiary.BeneficiaryUsername,
		IsActive:                     beneficiary.IsActive,
		BeneficiaryMobileCountryCode: beneficiary.BeneficiaryMobileCountryCode,
		BeneficiaryMobileNumber:      beneficiary.BeneficiaryMobileNumber,
	}

	return resp, nil
}

func (uc *beneficiaryUsecase) GetBeneficiaries(ctx context.Context, userID int) (*dto.GetBeneficiariesResponse, error) {
	beneficiaries, err := uc.beneficiaryRepository.GetBeneficiaries(ctx, userID)
	if err != nil {
		log.Println("error getting beneficiaries", err)
		return nil, err
	}

	// TODO: might remove this unnecessary for loop and just return domain.Beneficiary
	resp := &dto.GetBeneficiariesResponse{}
	for _, b := range *beneficiaries {
		beneficiary := dto.GetBeneficiaryResponse{
			BeneficiaryID:                b.BeneficiaryID,
			IsDeleted:                    b.IsDeleted,
			BeneficiaryFirstName:         b.BeneficiaryFirstName,
			BeneficiaryLastName:          b.BeneficiaryLastName,
			BeneficiaryEmail:             b.BeneficiaryEmail,
			BeneficiaryUsername:          b.BeneficiaryUsername,
			IsActive:                     b.IsActive,
			BeneficiaryMobileCountryCode: b.BeneficiaryMobileCountryCode,
			BeneficiaryMobileNumber:      b.BeneficiaryMobileNumber,
		}
		resp.Beneficiaries = append(resp.Beneficiaries, beneficiary)
	}

	return resp, nil
}
