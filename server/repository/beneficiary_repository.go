package repository

import (
	"context"
	"database/sql"
	"time"

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

func (r *beneficiaryRepository) GetBeneficiary(ctx context.Context, beneficiaryID int, userID int) (*domain.Beneficiary, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

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
	if err := r.db.GetContext(ctx, &beneficiary, query, userID, beneficiaryID); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrUserNotLinkedToBeneficiary
		}
		return nil, err
	}

	return &beneficiary, nil
}

func (r *beneficiaryRepository) GetBeneficiaries(ctx context.Context, userID int) (*[]domain.Beneficiary, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

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
		return nil, err
	}

	if len(beneficiaries) == 0 {
		return nil, exception.ErrUserHasNoBeneficiary
	}

	return &beneficiaries, nil
}

func (r *beneficiaryRepository) GetUserIDByMobileNumber(ctx context.Context, mobileCountryCode, mobileNumber string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT id 
		FROM users 
		WHERE 
			mobile_country_code = $1 AND 
			mobile_number = $2;
	`

	var id int
	if err := r.db.QueryRowContext(ctx, query, mobileCountryCode, mobileNumber).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, exception.ErrUserNotFound
		}
		return 0, err
	}

	return id, nil
}

func (r *beneficiaryRepository) CreateBeneficiary(ctx context.Context, userID int, beneficiaryID int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
        INSERT INTO user_beneficiary (user_id, beneficiary_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, beneficiary_id) DO NOTHING;
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
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT 1
		FROM user_beneficiary
		WHERE user_id = $1 AND beneficiary_id = $2;
	`

	var exists int
	if err := r.db.QueryRowContext(ctx, query, userID, beneficiaryID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return exception.ErrUserNotLinkedToBeneficiary
		}
		return err
	}

	return nil
}

func (r *beneficiaryRepository) UpdateBeneficiaryIsDeleted(ctx context.Context, userID int, beneficiaryID int, isDeleted int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		UPDATE user_beneficiary
		SET is_deleted = $1
		WHERE user_id = $2 AND beneficiary_id = $3;
	`

	if _, err := r.db.ExecContext(ctx, query, isDeleted, userID, beneficiaryID); err != nil {
		return err
	}

	return nil
}
