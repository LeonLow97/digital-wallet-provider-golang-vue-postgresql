package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func Test_GetByUsername(t *testing.T) {
	// Create a mock database connection using sqlmock
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer mockDB.Close()

	// Convert the mock database connection to an *sqlx.DB instance
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Define the expected rows that will be returned by the query
	rows := sqlmock.NewRows([]string{"id", "username", "password", "active", "admin"}).
		AddRow(1, "testUsername", "testPassword", 1, 1)

	// Set up the expectation for the query using the mock
	// mockDB (now sqlxDB) then responds with the rows defined in `NewRows(...).AddRow(...)`
	// Thus, we can simulate the behavior of the database without interacting with a real database
	mock.ExpectQuery("SELECT id, username, password, active, admin FROM users").WithArgs("testUsername").WillReturnRows(rows)

	// Create a repository instance using the mock database connection
	r, err := NewRepo(sqlxDB)
	require.NoError(t, err, "creating the shared repo")

	// Call the GetByUsername method on the repository
	results, err := r.GetByUsername(context.Background(), "testUsername")
	require.NoError(t, err, "running GetByUsername on repository instance")

	expectedResult := &User{
		ID:       1,
		Username: "testUsername",
		Password: "testPassword",
		Active:   1,
		Admin:    1,
	}

	// Compare the expected result with the actual result using assertions
	assert.Equal(t, fmt.Sprintf("%+v", expectedResult), fmt.Sprintf("%+v", results))
}
