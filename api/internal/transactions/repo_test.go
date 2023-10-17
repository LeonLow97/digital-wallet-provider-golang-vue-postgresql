package transactions

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDB struct {
	mock.Mock
}

func TestGetDB(t *testing.T) {
	// Create a mock database connection
	mockDB := &sqlx.DB{} // Replace this with your actual mock or real db

	// Create a repo instance with the mock database
	r := &repo{db: mockDB}

	// Call the GetDB method
	db := r.GetDB()

	// Assert that the returned database connection is the same as the mockDB
	assert.Equal(t, mockDB, db, "GetDB should return the mock database connection")
}

func Test_GetUserCountByUserId(t *testing.T) {
	testCases := []struct {
		Test          string
		UserId        int
		ExpectedCount int
		ExpectErr     bool
		QueryExpect   func(mock sqlmock.Sqlmock)
	}{
		{
			Test:          "Successfully returned count by userId",
			UserId:        1,
			ExpectedCount: 1,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

				mock.ExpectQuery("SELECT COUNT(*) FROM users WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			Test:          "Returned no rows",
			UserId:        513,
			ExpectedCount: 0,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*) FROM users WHERE id = $1;").
					WithArgs(513).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Test:          "Returned internal server error",
			UserId:        514,
			ExpectedCount: 0,
			ExpectErr:     true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*) FROM users WHERE id = $1;").
					WithArgs(514).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewRepo(sqlxDB)

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			tc.QueryExpect(mock)

			returnedCount, err := r.GetUserCountByUserId(context.Background(), tc.UserId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetUserCountByUserId on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			} else {
				require.Error(t, err, "running GetUserCountByUserId on repository layer expected error")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			}
		})
	}
}

func Test_GetUserIdByMobileNumber_Repo(t *testing.T) {
	testCases := []struct {
		Test          string
		MobileNumber  string
		ExpectedCount int
		ExpectErr     bool
		QueryExpect   func(mock sqlmock.Sqlmock)
	}{
		{
			Test:          "Successfully returned count by userId and beneficiaryId",
			MobileNumber:  "+555",
			ExpectedCount: 1,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

				mock.ExpectQuery("SELECT id FROM users where mobile_number = $1;").
					WithArgs("+555").WillReturnRows(rows)
			},
		},
		{
			Test:          "Returned no rows",
			MobileNumber:  "+556",
			ExpectedCount: 0,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id FROM users where mobile_number = $1;").
					WithArgs("+556").WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Test:          "Returned internal server error",
			MobileNumber:  "+557",
			ExpectedCount: 0,
			ExpectErr:     true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id FROM users where mobile_number = $1;").
					WithArgs("+557").WillReturnError(sql.ErrConnDone)
			},
		},
	}

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewRepo(sqlxDB)

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			tc.QueryExpect(mock)

			returnedCount, err := r.GetUserIdByMobileNumber(context.Background(), tc.MobileNumber)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetUserIdByMobileNumber on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			} else {
				require.Error(t, err, "running GetUserIdByMobileNumber on repository layer expected error")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			}
		})
	}
}

func Test_GetCountByUserIdAndBeneficiaryId_Repo(t *testing.T) {
	testCases := []struct {
		Test          string
		UserId        int
		BeneficiaryId int
		ExpectedCount int
		ExpectErr     bool
		QueryExpect   func(mock sqlmock.Sqlmock)
	}{
		{
			Test:          "Successfully returned count by userId and beneficiaryId",
			UserId:        1,
			BeneficiaryId: 2,
			ExpectedCount: 1,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

				mock.ExpectQuery("SELECT COUNT(*) FROM user_beneficiary WHERE user_id = $1 AND beneficiary_id = $2;").
					WithArgs(1, 2).WillReturnRows(rows)
			},
		},
		{
			Test:          "Returned no rows",
			UserId:        513,
			BeneficiaryId: 2,
			ExpectedCount: 0,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*) FROM user_beneficiary WHERE user_id = $1 AND beneficiary_id = $2;").
					WithArgs(513, 2).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Test:          "Returned internal server error",
			UserId:        514,
			BeneficiaryId: 2,
			ExpectedCount: 0,
			ExpectErr:     true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*) FROM user_beneficiary WHERE user_id = $1 AND beneficiary_id = $2;").
					WithArgs(514, 2).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewRepo(sqlxDB)

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			tc.QueryExpect(mock)

			returnedCount, err := r.GetCountByUserIdAndBeneficiaryId(context.Background(), tc.UserId, tc.BeneficiaryId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetCountByUserIdAndBeneficiaryId on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			} else {
				require.Error(t, err, "running GetCountByUserIdAndBeneficiaryId on repository layer expected error")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			}
		})
	}
}

func Test_GetCountByUserIdAndCurrency_Repo(t *testing.T) {
	testCases := []struct {
		Test          string
		UserId        int
		Currency      string
		ExpectedCount int
		ExpectedId    int
		ExpectErr     bool
		QueryExpect   func(mock sqlmock.Sqlmock)
	}{
		{
			Test:          "Successfully returned count by userId and beneficiaryId",
			UserId:        1,
			Currency:      "SGD",
			ExpectedCount: 1,
			ExpectedId:    1,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)", "id"}).AddRow(1, 1)

				mock.ExpectQuery("SELECT COUNT(*), id FROM user_balance WHERE user_id = $1 AND currency = $2 GROUP BY (id);").
					WithArgs(1, "SGD").WillReturnRows(rows)
			},
		},
		{
			Test:          "Returned no rows",
			UserId:        513,
			Currency:      "SGD",
			ExpectedCount: 0,
			ExpectedId:    0,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*), id FROM user_balance WHERE user_id = $1 AND currency = $2 GROUP BY (id);").
					WithArgs(513, "SGD").WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Test:          "Returned internal server error",
			UserId:        514,
			Currency:      "SGD",
			ExpectedCount: 0,
			ExpectedId:    0,
			ExpectErr:     true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*), id FROM user_balance WHERE user_id = $1 AND currency = $2 GROUP BY (id);").
					WithArgs(514, "SGD").WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// Create a new database connection mock for each test case
			mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %s", err)
			}
			defer mockDB.Close()

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			r := NewRepo(sqlxDB)

			tc.QueryExpect(mock)

			returnedCount, returnedId, err := r.GetCountByUserIdAndCurrency(context.Background(), tc.UserId, tc.Currency)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetCountByUserIdAndCurrency on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
				assert.Equal(t, returnedId, tc.ExpectedId)
			} else {
				require.Error(t, err, "running GetCountByUserIdAndCurrency on repository layer expected error")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
				assert.Equal(t, returnedId, tc.ExpectedId)
			}
		})
	}
}

func Test_GetBalanceIdByUserIdAndPrimary_Repo(t *testing.T) {
	testCases := []struct {
		Test             string
		UserId           int
		ExpectedCount    int
		ExpectedCurrency string
		ExpectErr        bool
		QueryExpect      func(mock sqlmock.Sqlmock)
	}{
		{
			Test:             "Successfully returned count by userId and beneficiaryId",
			UserId:           1,
			ExpectedCount:    1,
			ExpectedCurrency: "SGD",
			ExpectErr:        false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "currency"}).AddRow(1, "SGD")

				mock.ExpectQuery("SELECT id, currency FROM user_balance WHERE user_id = $1 AND is_primary = 1;").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			Test:             "Returned no rows",
			UserId:           513,
			ExpectedCount:    0,
			ExpectedCurrency: "",
			ExpectErr:        false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, currency FROM user_balance WHERE user_id = $1 AND is_primary = 1;").
					WithArgs(513).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Test:             "Returned internal server error",
			UserId:           514,
			ExpectedCount:    0,
			ExpectedCurrency: "",
			ExpectErr:        true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, currency FROM user_balance WHERE user_id = $1 AND is_primary = 1;").
					WithArgs(514).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// Create a new database connection mock for each test case
			mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %s", err)
			}
			defer mockDB.Close()

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			r := NewRepo(sqlxDB)

			tc.QueryExpect(mock)

			returnedCount, returnedCurrency, err := r.GetBalanceIdByUserIdAndPrimary(context.Background(), tc.UserId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetBalanceIdByUserIdAndPrimary on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
				assert.Equal(t, returnedCurrency, tc.ExpectedCurrency)
			} else {
				require.Error(t, err, "running GetBalanceIdByUserIdAndPrimary on repository layer expected error")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
				assert.Equal(t, returnedCurrency, tc.ExpectedCurrency)
			}
		})
	}
}

func Test_GetBalanceAmountById_Repo(t *testing.T) {
	testCases := []struct {
		Test            string
		BalanceId       int
		ExpectedBalance float64
		ExpectErr       bool
		QueryExpect     func(mock sqlmock.Sqlmock)
	}{
		{
			Test:            "Successfully returned count by userId and beneficiaryId",
			BalanceId:       1,
			ExpectedBalance: 100.0,
			ExpectErr:       false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(100.0)

				mock.ExpectQuery("SELECT balance FROM user_balance WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			Test:            "Returned no rows",
			BalanceId:       513,
			ExpectedBalance: 0,
			ExpectErr:       false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT balance FROM user_balance WHERE id = $1;").
					WithArgs(513).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Test:            "Returned internal server error",
			BalanceId:       514,
			ExpectedBalance: 0,
			ExpectErr:       true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT balance FROM user_balance WHERE id = $1;").
					WithArgs(514).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// Create a new database connection mock for each test case
			mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %s", err)
			}
			defer mockDB.Close()

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			r := NewRepo(sqlxDB)

			mock.ExpectBegin()
			tc.QueryExpect(mock)

			if tc.ExpectErr {
				// Expect a rollback if an error is expected
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			tx, err := mockDB.Begin()
			if err != nil {
				t.Fatalf("Error creating mock transaction: %v", err)
			}
			defer tx.Rollback()

			returnedBalance, err := r.GetBalanceAmountById(tx, context.Background(), tc.BalanceId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetBalanceAmountById on repository layer")
				assert.Equal(t, returnedBalance, tc.ExpectedBalance)
			} else {
				require.Error(t, err, "running GetBalanceAmountById on repository layer expected error")
				assert.Equal(t, returnedBalance, tc.ExpectedBalance)
			}
		})
	}
}

func Test_UpdateBalanceAmountById_Repo(t *testing.T) {
	testCases := []struct {
		Test        string
		Balance     float64
		BalanceId   int
		ExpectErr   bool
		QueryExpect func(mock sqlmock.Sqlmock)
	}{
		{
			Test:      "Successfully update balance by ID",
			Balance:   100.0,
			BalanceId: 1,
			ExpectErr: false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE user_balance SET balance = $1 WHERE id = $2;").
					WithArgs(100.0, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			Test:      "Returned internal server error",
			Balance:   100.0,
			BalanceId: 514,
			ExpectErr: true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("UPDATE user_balance SET balance = $1 WHERE id = $2;").
					WithArgs(100.0, 514).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// Create a new database connection mock for each test case
			mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %s", err)
			}
			defer mockDB.Close()

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			r := NewRepo(sqlxDB)

			mock.ExpectBegin()
			tc.QueryExpect(mock)

			if tc.ExpectErr {
				// Expect a rollback if an error is expected
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			tx, err := mockDB.Begin()
			if err != nil {
				t.Fatalf("Error creating mock transaction: %v", err)
			}
			defer tx.Rollback()

			err = r.UpdateBalanceAmountById(tx, context.Background(), tc.Balance, tc.BalanceId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running UpdateBalanceAmountById on repository layer")

			} else {
				require.Error(t, err, "running UpdateBalanceAmountById on repository layer expected error")
			}
		})
	}
}

func Test_InsertIntoTransactions_Repo(t *testing.T) {
	testCases := []struct {
		Test              string
		TransactionEntity *TransactionEntity
		ExpectErr         bool
		QueryExpect       func(mock sqlmock.Sqlmock)
	}{
		{
			Test: "Successfully update balance by ID",
			TransactionEntity: &TransactionEntity{
				UserId:                    1,
				SenderId:                  1,
				BeneficiaryId:             2,
				TransferredAmount:         100.0,
				TransferredAmountCurrency: "SGD",
				ReceivedAmount:            100.0,
				ReceivedAmountCurrency:    "SGD",
				Status:                    "COMPLETED",
				TransferredDate:           time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:              time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO transactions (
					user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency,
					received_amount, received_amount_currency, status, transferred_date, received_date
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
				`).
					WithArgs(1, 1, 2, 100.0, "SGD", 100.0, "SGD", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			Test: "Returned internal server error",
			TransactionEntity: &TransactionEntity{
				UserId:                    2,
				SenderId:                  3,
				BeneficiaryId:             7,
				TransferredAmount:         100.0,
				TransferredAmountCurrency: "SGD",
				ReceivedAmount:            100.0,
				ReceivedAmountCurrency:    "SGD",
				Status:                    "COMPLETED",
				TransferredDate:           time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:              time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO transactions (
					user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency,
					received_amount, received_amount_currency, status, transferred_date, received_date
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
				`).
					WithArgs(2, 3, 7, 100.0, "SGD", 100.0, "SGD", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// Create a new database connection mock for each test case
			mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %s", err)
			}
			defer mockDB.Close()

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			r := NewRepo(sqlxDB)

			mock.ExpectBegin()
			tc.QueryExpect(mock)

			if tc.ExpectErr {
				// Expect a rollback if an error is expected
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			tx, err := mockDB.Begin()
			if err != nil {
				t.Fatalf("Error creating mock transaction: %v", err)
			}
			defer tx.Rollback()

			err = r.InsertIntoTransactions(tx, context.Background(), tc.TransactionEntity)

			if !tc.ExpectErr {
				require.NoError(t, err, "running InsertIntoTransactions on repository layer")
			} else {
				require.Error(t, err, "running InsertIntoTransactions on repository layer expected error")
			}
		})
	}
}

func Test_GetTransactionsCountByUserId_Repo(t *testing.T) {
	testCases := []struct {
		Test          string
		UserId        int
		ExpectedCount int
		ExpectErr     bool
		QueryExpect   func(mock sqlmock.Sqlmock)
	}{
		{
			Test:          "Successfully returned count for transactions",
			UserId:        1,
			ExpectedCount: 1,
			ExpectErr:     false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

				mock.ExpectQuery("SELECT COUNT(*) FROM transactions WHERE user_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			Test:          "Returned no rows",
			UserId:        513,
			ExpectedCount: 0,
			ExpectErr:     true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(*) FROM transactions WHERE user_id = $1;").
					WithArgs(513).WillReturnError(sql.ErrNoRows)
			},
		},
	}

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewRepo(sqlxDB)

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			tc.QueryExpect(mock)

			returnedCount, err := r.GetTransactionsCountByUserId(context.Background(), tc.UserId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetTransactionsCountByUserId on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			} else {
				require.Error(t, err, "running GetTransactionsCountByUserId on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			}
		})
	}
}

func Test_GetTransactionsByUserId_Repo(t *testing.T) {
	testCases := []struct {
		Test                 string
		UserId               int
		PageSize             int
		Offset               int
		ExpectedTransactions *Transactions
		ExpectErr            bool
		QueryExpect          func(mock sqlmock.Sqlmock)
	}{
		{
			Test:     "Successfully retrieve the list of transactions",
			UserId:   1,
			PageSize: 1,
			Offset:   20,
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName:           "John",
						SenderLastName:            "Doe",
						SenderUsername:            "johndoe",
						BeneficiaryFirstName:      "Jane",
						BeneficiaryLastName:       "Smith",
						BeneficiaryUsername:       "janesmith",
						TransferredAmount:         100.0,
						TransferredAmountCurrency: "USD",
						ReceivedAmount:            90.0,
						ReceivedAmountCurrency:    "USD",
						Status:                    "Completed",
						TransferredDate:           time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
						ReceivedDate:              time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			ExpectErr: false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"sender_first_name", "sender_last_name", "sender_username",
					"beneficiary_first_name", "beneficiary_last_name", "beneficiary_username",
					"transferred_amount", "transferred_amount_currency",
					"received_amount", "received_amount_currency",
					"status", "transferred_date", "received_date"}).
					AddRow("John", "Doe", "johndoe", "Jane", "Smith", "janesmith", 100.0, "USD", 90.0, "USD", "Completed", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC))

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT COALESCE(s.first_name, '') AS sender_first_name, COALESCE(s.last_name, '') AS sender_last_name, s.username AS sender_username, COALESCE(b.first_name, '') AS beneficiary_first_name, COALESCE(b.last_name, '') AS beneficiary_last_name, b.username AS beneficiary_username, t.transferred_amount, t.transferred_amount_currency, t.received_amount, t.received_amount_currency, t.status, t.transferred_date, t.received_date FROM transactions t JOIN users s ON s.id = t.sender_id JOIN users b ON b.id = t.beneficiary_id WHERE t.user_id = $1 ORDER BY t.transferred_date DESC LIMIT $2 OFFSET $3;"),
				).WithArgs(1, 1, 20).WillReturnRows(rows)
			},
		},
		{
			Test:                 "Internal Server Error at QueryContext",
			UserId:               2,
			PageSize:             1,
			Offset:               20,
			ExpectedTransactions: nil,
			ExpectErr:            true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT COALESCE(s.first_name, '') AS sender_first_name, COALESCE(s.last_name, '') AS sender_last_name, s.username AS sender_username, COALESCE(b.first_name, '') AS beneficiary_first_name, COALESCE(b.last_name, '') AS beneficiary_last_name, b.username AS beneficiary_username, t.transferred_amount, t.transferred_amount_currency, t.received_amount, t.received_amount_currency, t.status, t.transferred_date, t.received_date FROM transactions t JOIN users s ON s.id = t.sender_id JOIN users b ON b.id = t.beneficiary_id WHERE t.user_id = $1 ORDER BY t.transferred_date DESC LIMIT $2 OFFSET $3;"),
				).WithArgs(2, 1, 20).WillReturnError(sql.ErrNoRows)
			},
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewRepo(sqlxDB)
	require.NoError(t, err, "creating the shared repo")

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			tc.QueryExpect(mock)

			transactions, err := r.GetTransactionsByUserId(context.Background(), tc.UserId, tc.PageSize, tc.Offset)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetTransactionsByUserId on repository layer")
				assert.Equal(t, fmt.Sprintf("%+v", transactions), fmt.Sprintf("%+v", tc.ExpectedTransactions))
			} else {
				require.Error(t, err, "running GetTransactionsByUserId on repository layer")
				assert.Nil(t, transactions)
			}
		})
	}
}

func Test_CreateTransactionSQLTransaction_Repo(t *testing.T) {
	testCases := []struct {
		Test                                    string
		SenderId                                int
		BeneficiaryId                           int
		SenderTransaction                       *TransactionEntity
		BeneficiaryTransaction                  *TransactionEntity
		ExpectErr                               bool
		QueryExpectForSenderBalanceAmount       func(mock sqlmock.Sqlmock)
		QueryExpectForBeneficiaryBalanceAmount  func(mock sqlmock.Sqlmock)
		QueryUpdateSenderBalanceAmountById      func(mock sqlmock.Sqlmock)
		QueryUpdateBeneficiaryBalanceAmountById func(mock sqlmock.Sqlmock)
		QueryInsertSenderBalanceAmountById      func(mock sqlmock.Sqlmock)
		QueryInsertBeneficiaryBalanceAmountById func(mock sqlmock.Sqlmock)
	}{
		{
			Test:     "Successfully created a transaction",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         2,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: false,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(2).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Internal Server Error in GetBalanceAmountById for user balance",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         510,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         2,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(510).WillReturnError(sql.ErrConnDone)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(2).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Insufficient funds in GetBalanceAmountById for user balance",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         410,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         2,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(4000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(410).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(2).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Internal Server Error in GetBalanceAmountById for beneficiary balance",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         511,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(511).WillReturnError(sql.ErrConnDone)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "beneficiaryTransaction.BeneficiaryHasTransferredCurrency == 0",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         3,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: false,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(3).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Internal Server Error in UpdateBalanceAmountById for sender",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         3,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(3).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnError(sql.ErrConnDone)
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Internal Server Error in UpdateBalanceAmountById for beneficiary",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         3,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(3).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnError(sql.ErrConnDone)
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Internal Server Error in InsertIntoTransactions for sender",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         3,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(3).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnError(sql.ErrConnDone)
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Test:     "Internal Server Error in InsertIntoTransactions for beneficiary",
			SenderId: 12,
			SenderTransaction: &TransactionEntity{
				UserId:                            12,
				SenderId:                          12,
				BalanceId:                         1,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 0,
				Status:                            "COMPLETED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			BeneficiaryTransaction: &TransactionEntity{
				UserId:                            13,
				SenderId:                          12,
				BalanceId:                         3,
				BeneficiaryId:                     13,
				TransferredAmount:                 9000,
				TransferredAmountCurrency:         "YEN",
				ReceivedAmount:                    9000,
				ReceivedAmountCurrency:            "YEN",
				BeneficiaryHasTransferredCurrency: 1,
				Status:                            "RECEIVED",
				TransferredDate:                   time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				ReceivedDate:                      time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
			},
			ExpectErr: true,
			QueryExpectForSenderBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(1).WillReturnRows(rows)
			},
			QueryExpectForBeneficiaryBalanceAmount: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(16000)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM user_balance WHERE id = $1;")).WithArgs(3).WillReturnRows(rows)
			},
			QueryUpdateSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			QueryUpdateBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE user_balance SET balance = $1 WHERE id = $2;")).
					WillReturnResult(sqlmock.NewResult(13, 1))
			},
			QueryInsertSenderBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					12, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "COMPLETED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			QueryInsertBeneficiaryBalanceAmountById: func(mock sqlmock.Sqlmock) {
				regexPattern := regexp.QuoteMeta("INSERT INTO transactions ( user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status, transferred_date, received_date ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")

				mock.ExpectExec(regexPattern).WithArgs(
					13, 12, 13, 9000.0, "YEN", 9000.0, "YEN", "RECEIVED", time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %s", err)
			}
			defer mockDB.Close()

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			r := NewRepo(sqlxDB)

			mock.ExpectBegin()
			tc.QueryExpectForSenderBalanceAmount(mock)
			tc.QueryExpectForBeneficiaryBalanceAmount(mock)
			tc.QueryUpdateSenderBalanceAmountById(mock)
			tc.QueryUpdateBeneficiaryBalanceAmountById(mock)
			tc.QueryInsertSenderBalanceAmountById(mock)
			tc.QueryInsertBeneficiaryBalanceAmountById(mock)

			if tc.ExpectErr {
				// Expect a rollback if an error is expected
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			err = r.CreateTransactionSQLTransaction(context.Background(), tc.SenderTransaction, tc.BeneficiaryTransaction)

			if !tc.ExpectErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
