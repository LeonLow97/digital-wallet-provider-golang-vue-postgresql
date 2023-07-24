package beneficiaries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetByUsername(ctx context.Context, username string) (*Beneficiaries, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) GetByUsername(ctx context.Context, username string) (*Beneficiaries, error) {

	query := `SELECT ub.beneficiary_id, b.beneficiary_name, b.mobile_number, b.currency
				FROM user_beneficiary ub
				LEFT JOIN beneficiaries b
				ON ub.beneficiary_id = b.beneficiary_id
				WHERE ub.user_id = (SELECT id FROM users WHERE username = ?);`

	var BeneficiaryId sql.NullInt64
	var BeneficiaryName, MobileNumber, Currency sql.NullString

	query = r.db.Rebind(query)

	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("error with QueryContext: %s", err)
	}

	defer rows.Close()

	var beneficiaries Beneficiaries

	for rows.Next() {
		var beneficiary Beneficiary
		if err := rows.Scan(&BeneficiaryId, &BeneficiaryName, &MobileNumber, &Currency); err != nil {
			return nil, fmt.Errorf("error scanning results: %s", err)
		}

		beneficiary = Beneficiary{
			BeneficiaryId:   int(BeneficiaryId.Int64),
			BeneficiaryName: BeneficiaryName.String,
			MobileNumber:    MobileNumber.String,
			Currency:        Currency.String,
		}

		beneficiaries.Beneficiaries = append(beneficiaries.Beneficiaries, beneficiary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over result rows: %s", err)
	}

	return &beneficiaries, nil
}
