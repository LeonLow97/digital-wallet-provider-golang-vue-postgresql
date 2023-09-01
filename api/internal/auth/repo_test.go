package auth

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func Test_GetByUsername(t *testing.T) {
	testCases := []struct {
		Test         string
		Username     string
		ExpectedUser *User
		ExpectErr    bool
		QueryExpect  func(mock sqlmock.Sqlmock)
	}{
		{
			Test:     "Valid User",
			Username: "validUsername",
			ExpectedUser: &User{
				ID:       1,
				Username: "validUsername",
				Password: "validPassword",
				Active:   1,
				Admin:    1,
			},
			ExpectErr: false,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				// Define the expected rows that will be returned by the query
				rows := sqlmock.NewRows([]string{"id", "username", "password", "active", "admin"}).
					AddRow(1, "validUsername", "validPassword", 1, 1)

				mock.ExpectQuery("SELECT id, username, password, active, admin FROM users WHERE username = $1").
					WithArgs("validUsername").
					WillReturnRows(rows)
			},
		},
		{
			Test:         "Non-Existent User",
			Username:     "invalidUsername",
			ExpectedUser: nil,
			ExpectErr:    true,
			QueryExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, username, password, active, admin FROM users WHERE username = $1").
					WithArgs("invalidUsername").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	// Create a mock database connection using sqlmock
	// mockDB, mock, err := sqlmock.New() // less restricted query string
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	// Convert the mock database connection to an *sqlx.DB instance
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Create a repository instance using the mock database connection
	r, err := NewRepo(sqlxDB)
	require.NoError(t, err, "creating the shared repo")

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// generating the expected query results
			tc.QueryExpect(mock)

			// call the GetByUsername method in the repository layer
			results, err := r.GetByUsername(context.Background(), tc.Username)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetByUsername on repository layer")
				assert.Equal(t, fmt.Sprintf("%+v", results), fmt.Sprintf("%+v", tc.ExpectedUser))
			} else {
				require.Error(t, err, "running GetByUsername on repository layer")
				assert.Equal(t, results, tc.ExpectedUser)
			}
		})
	}
}
