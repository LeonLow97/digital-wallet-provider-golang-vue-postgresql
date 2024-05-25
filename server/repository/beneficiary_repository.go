package repository

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/jmoiron/sqlx"
)

type beneficiaryRepository struct {
	db *sqlx.DB
}

func NewBeneficiaryRepository(db *sqlx.DB) domain.BeneficiaryRepository {
	return &beneficiaryRepository{
		db: db,
	}
}

func (r *beneficiaryRepository) GetUserIDByMobileNumber(ctx context.Context, mobileCountryCode, mobileNumber string) (int, error) {
	query := `SELECT id FROM users WHERE mobile_country_code = $1 AND mobile_number = $2`

	var id int
	err := r.db.QueryRowContext(ctx, query, mobileCountryCode, mobileNumber).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, exception.ErrUserNotFound
		}
		return 0, err
	}
	return id, nil
}

func (r *beneficiaryRepository) CreateBeneficiary(ctx context.Context, userID int, beneficiaryID int) error {
	query := `
		INSERT INTO user_beneficiary (user_id, beneficiary_id)
		SELECT $1, $2
		WHERE NOT EXISTS
			(SELECT 1 FROM user_beneficiary WHERE user_id = $1 AND beneficiary_id = $2);
	`

	result, err := r.db.ExecContext(ctx, query, userID, beneficiaryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return exception.ErrBeneficiaryAlreadyExists
	}

	return nil
}

func (r *beneficiaryRepository) IsLinkedByUserIDAndBeneficiaryID(ctx context.Context, userID int, beneficiaryID int) error {
	query := `
		SELECT 1
		FROM user_beneficiary
		WHERE user_id = $1 AND beneficiary_id = $2;
	`

	var exists int
	err := r.db.QueryRowContext(ctx, query, userID, beneficiaryID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return exception.ErrUserNotLinkedToBeneficiary
		}
		return err
	}

	return nil
}

func (r *beneficiaryRepository) UpdateBeneficiaryIsDeleted(ctx context.Context, userID int, beneficiaryID int, isDeleted int) error {
	query := `
		UPDATE user_beneficiary
		SET is_deleted = $1
		WHERE user_id = $2 AND beneficiary_id = $3;
	`

	_, err := r.db.ExecContext(ctx, query, isDeleted, userID, beneficiaryID)
	if err != nil {
		return err
	}

	return nil
}

func (r *beneficiaryRepository) GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*domain.Beneficiary, error) {
	query := `
		SELECT 
			ub.beneficiary_id, ub.is_deleted, u.first_name, u.last_name, u.email,
			u.username, u.active, u.mobile_country_code, u.mobile_number
		FROM user_beneficiary ub
		JOIN users u
			ON u.id = ub.beneficiary_id
		WHERE ub.user_id = $1 AND ub.beneficiary_id = $2;
	`

	var beneficiary domain.Beneficiary
	err := r.db.QueryRowContext(ctx, query, userID, beneficiaryID).Scan(
		&beneficiary.BeneficiaryID,
		&beneficiary.IsDeleted,
		&beneficiary.BeneficiaryFirstName,
		&beneficiary.BeneficiaryLastName,
		&beneficiary.BeneficiaryEmail,
		&beneficiary.BeneficiaryUsername,
		&beneficiary.IsActive,
		&beneficiary.BeneficiaryMobileCountryCode,
		&beneficiary.BeneficiaryMobileNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrUserNotLinkedToBeneficiary
		}
		return nil, err
	}

	return &beneficiary, nil
}

func (r *beneficiaryRepository) GetBeneficiaries(ctx context.Context, userID int) (*[]domain.Beneficiary, error) {
	query := `
		SELECT 
			ub.beneficiary_id, ub.is_deleted, u.first_name, u.last_name, u.email,
			u.username, u.active, u.mobile_country_code, u.mobile_number
		FROM user_beneficiary ub
		JOIN users u
			ON u.id = ub.beneficiary_id
		WHERE ub.user_id = $1;
	`

	var beneficiaries []domain.Beneficiary
	if err := r.db.SelectContext(ctx, &beneficiaries, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrUserHasNoBeneficiary
		}
		return nil, err
	}

	return &beneficiaries, nil
}
