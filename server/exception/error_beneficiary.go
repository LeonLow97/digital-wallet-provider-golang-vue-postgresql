package exception

import "errors"

var (
	ErrUserIDEqualBeneficiaryID = errors.New("user id is equal to beneficiary id")
	ErrBeneficiaryAlreadyExists = errors.New("beneficiary already exists")

	ErrUserNotLinkedToBeneficiary = errors.New("user id is not linked to beneficiary id")
	ErrUserHasNoBeneficiary       = errors.New("user id is not linked to any beneficiaries")
)
