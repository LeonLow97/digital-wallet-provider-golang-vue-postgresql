package transactions

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func Test_GetCountByUserId_Repository(t *testing.T) {
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

	r, err := NewRepo(sqlxDB)
	require.NoError(t, err, "creating the shared repo")

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			tc.QueryExpect(mock)

			returnedCount, err := r.GetCountByUserId(context.Background(), tc.UserId)

			if !tc.ExpectErr {
				require.NoError(t, err, "running GetCountByUserId on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			} else {
				require.Error(t, err, "running GetCountByUserId on repository layer")
				assert.Equal(t, returnedCount, tc.ExpectedCount)
			}
		})
	}
}

// var (
// 	host     = "localhost"
// 	user     = "postgres"
// 	password = "postgres"
// 	dbName   = "db_test"
// 	port     = "5437"
// 	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=30"
// )

// var resource *dockertest.Resource
// var pool *dockertest.Pool
// var testDB *sqlx.DB
// var testRepo Repo

// func TestMain(m *testing.M) {
// 	// connect to docker; fail if docker not running
// 	p, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("could not connect to docker; is it running? %s", err)
// 	}
// 	pool = p

// 	opts := dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "14.5",
// 		Env: []string{
// 			"POSTGRES_USER=" + user,
// 			"POSTGRES_PASSWORD=" + password,
// 			"POSTGRES_DB=" + dbName,
// 		},
// 		ExposedPorts: []string{"5432"},
// 		PortBindings: map[docker.Port][]docker.PortBinding{
// 			"5432": {
// 				{HostIP: "0.0.0.0", HostPort: port},
// 			},
// 		},
// 	}

// 	resource, err = pool.RunWithOptions(&opts)
// 	if err != nil {
// 		// _ = pool.Purge(resource)
// 		log.Fatalf("could not start resource: %s", err)
// 	}

// 	if err := pool.Retry(func() error {
// 		var err error
// 		testDB, err = sqlx.Connect("postgres", fmt.Sprintf(dsn, host, port, user, password, dbName))
// 		pool.MaxWait = 20 * time.Minute
// 		if err != nil {
// 			return err
// 		}
// 		return testDB.Ping()
// 	}); err != nil {
// 		_ = pool.Purge(resource)
// 		log.Fatalf("could not connect to database: %s", err)
// 	}

// 	err = createTables()
// 	if err != nil {
// 		log.Fatalf("error creating tables: %s", err)
// 	}

// 	code := m.Run()

// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("could not purge resource: %s", err)
// 	}

// 	testRepo = &repo{db: testDB}

// 	os.Exit(code)
// }

// func createTables() error {
// 	tableSQL, err := os.ReadFile("./testdata/tables.sql")
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	_, err = testDB.Exec(string(tableSQL))
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	return nil
// }

// func Test_pingDB(t *testing.T) {
// 	err := testDB.Ping()
// 	if err != nil {
// 		t.Error("can't ping database")
// 	}
// }

// func insertTestUser(username, password, mobileNumber string) error {
// 	query := `INSERT INTO users (username, password, mobile_number) VALUES ($1, $2, $3);`
// 	query = testDB.Rebind(query)

// 	_, err := testDB.Exec(query, username, password, mobileNumber)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func insertTestBeneficiary(beneficiaryName, beneficiaryNumber, currency string, isInternal int) error {
// 	query := `INSERT INTO beneficiaries (beneficiary_name, mobile_number, currency, is_internal)
// 				VALUES ($1, $2, $3, $4);`

// 	query = testDB.Rebind(query)

// 	_, err := testDB.ExecContext(context.Background(), query, beneficiaryName, beneficiaryNumber, currency, isInternal)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func insertTestUserBeneficiaryMapping(username, beneficiaryNumber string) error {
// 	query := `INSERT INTO user_beneficiary (user_id, beneficiary_id)
// 				VALUES (
// 					(SELECT id FROM users WHERE username = $1),
// 					(SELECT beneficiary_id FROM beneficiaries WHERE mobile_number = $2)
// 				);`

// 	query = testDB.Rebind(query)

// 	_, err := testDB.ExecContext(context.Background(), query, username, beneficiaryNumber)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func insertTestUserBalance(username string, balance float64, currency, countryISOCode string) error {
// 	query := `INSERT INTO user_balance (user_id, balance, currency, country_iso_code)
// 				VALUES ((SELECT id FROM users WHERE username = $1), $2, $3, $4);`

// 	query = testDB.Rebind(query)

// 	_, err := testDB.Exec(query, username, balance, currency, countryISOCode)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func getTestUserBalance(username string) (float64, error) {
// 	var balance float64
// 	query := `SELECT balance FROM user_balance WHERE user_id = (SELECT id FROM users WHERE username = $1);`
// 	query = testDB.Rebind(query)

// 	err := testDB.QueryRowContext(context.Background(), query, username).Scan(&balance)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return balance, nil
// }

// // func getTestUserTransactions(username string) ([]Transaction, error) {
// // 	var transactions []Transaction
// // 	query := `SELECT u.username, b.beneficiary_name, t.amount_transferred, t.amount_transferred_currency, t.amount_received,
// // 			t.amount_received_currency, t.status, t.date_transferred, t.date_received
// // 			FROM transactions t
// // 			LEFT JOIN users u ON u.id = t.sender_id
// // 			LEFT JOIN beneficiaries b ON b.beneficiary_id = t.beneficiary_id
// // 			WHERE t.user_id = (SELECT id FROM users WHERE username = $1);`

// // 	query = testDB.Rebind(query)

// // 	rows, err := testDB.QueryContext(context.Background(), query, username)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	for rows.Next() {
// // 		var transaction Transaction
// // 		if err := rows.Scan(&transaction.SenderName, &transaction.BeneficiaryName, &transaction.AmountTransferred, &transaction.AmountTransferredCurrency, &transaction.AmountReceived, &transaction.AmountReceivedCurrency, &transaction.Status, &transaction.DateTransferred, &transaction.DateReceived); err != nil {
// // 			return nil, err
// // 		}
// // 		transactions = append(transactions, transaction)
// // 	}

// // 	if err = rows.Err(); err != nil {
// // 		return nil, err
// // 	}

// // 	return transactions, nil
// // }

// func insertTestTransaction(senderName, beneficiaryName string, amountTransferred float64, amountTransferredCurrency string, amountReceived float64, amountReceivedCurrency string, status string) error {
// 	query := `INSERT INTO transactions (user_id, sender_id, beneficiary_id, amount_transferred, amount_transferred_currency, amount_received, amount_received_currency, status, date_transferred, date_received)
// 				VALUES (
// 					(SELECT id FROM users WHERE username = $1),
// 					(SELECT id FROM users WHERE username = $2),
// 					(SELECT beneficiary_id FROM beneficiaries WHERE beneficiary_name = $3),
// 					$4, $5, $6, $7, $8, $9, $10
// 				);`

// 	query = testDB.Rebind(query)

// 	_, err := testDB.ExecContext(context.Background(), query, senderName, senderName, beneficiaryName, amountTransferred, amountTransferredCurrency, amountReceived, amountReceivedCurrency, status, time.Now(), time.Now())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func cleanupTestData() error {
// 	// Run SQL DELETE statements to remove data from test tables
// 	_, err := testDB.Exec("TRUNCATE TABLE transactions, user_beneficiary, user_balance, beneficiaries, users CASCADE;")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func TestGetCountByUsername(t *testing.T) {
// 	testRepo := &repo{db: testDB}

// 	testUser := "testuser"

// 	initialCount, err := testRepo.GetCountByUsername(context.Background(), testUser)
// 	if err != nil {
// 		t.Fatalf("error getting initial count: %s", err)
// 	}

// 	if initialCount != 0 {
// 		t.Fatalf("expected initial count to be 0, but got %d", initialCount)
// 	}

// 	insertTestUser("testuser", "testpassword", "1234567890")
// 	count, err := testRepo.GetCountByUsername(context.Background(), "testuser")
// 	if err != nil {
// 		t.Fatalf("error getting count by username: %s", err)
// 	}

// 	if count != 1 {
// 		t.Fatalf("expected count to be 1, but got %d", count)
// 	}
// }

// func TestGetCountByUsernameAndBeneficiaryNumber(t *testing.T) {
// 	testRepo := &repo{db: testDB}

// 	insertTestUser("testuser", "testpassword", "1234567890")
// 	insertTestBeneficiary("testbeneficiary", "9876543210", "EUR", 1)
// 	insertTestUserBeneficiaryMapping("testuser", "9876543210")

// 	username := "testuser"
// 	beneficiaryNumber := "9876543210"

// 	count, err := testRepo.GetCountByUsernameAndBeneficiaryNumber(context.Background(), username, beneficiaryNumber)
// 	if err != nil {
// 		t.Fatalf("error calling GetCountByUsernameAndBeneficiaryNumber: %s", err)
// 	}

// 	if count != 1 {
// 		t.Fatalf("expected count to be 1, but got %d", count)
// 	}
// }

// func TestGetCountByBeneficiaryNameAndBeneficiaryNumber(t *testing.T) {
// 	testRepo := &repo{db: testDB}

// 	insertTestUser("testuser", "testpassword", "1234567890")

// 	beneficiaryName := "testuser"
// 	beneficiaryNumber := "1234567890"

// 	count, err := testRepo.GetCountByBeneficiaryNameAndBeneficiaryNumber(context.Background(), beneficiaryName, beneficiaryNumber)
// 	if err != nil {
// 		t.Fatalf("error calling GetCountByBeneficiaryNameAndBeneficiaryNumber: %s", err)
// 	}

// 	if count != 1 {
// 		t.Fatalf("expected count to be 1, but got %d", count)
// 	}
// }

// func TestGetUsernameByBeneficiaryNumber(t *testing.T) {
// 	testRepo := &repo{db: testDB}

// 	insertTestUser("testuser", "testpassword", "1234567890")

// 	beneficiaryNumber := "1234567890"
// 	actualUsername, err := testRepo.GetUsernameByBeneficiaryNumber(context.Background(), beneficiaryNumber)
// 	if err != nil {
// 		t.Fatalf("error calling GetUsernameByBeneficiaryNumber: %s", err)
// 	}

// 	expectedUsername := "testuser"
// 	if actualUsername != expectedUsername {
// 		t.Fatalf("expected username to be %s, but got %s", expectedUsername, actualUsername)
// 	}
// }

// func TestGetUserBalanceByUsername(t *testing.T) {
// 	testRepo := &repo{db: testDB}

// 	insertTestUser("testuser", "testpassword", "1234567890") // this user will have a userID of 1
// 	insertTestUserBalance("testuser", 100.00, "USD", "US")   // inserted user balance with userID of 1

// 	insertedUsername := "testuser"

// 	actualUserBalance, err := testRepo.GetUserBalanceByUsername(context.Background(), insertedUsername)
// 	if err != nil {
// 		t.Fatalf("error calling GetUserBalanceByUsername: %s", err)
// 	}

// 	expectedUserBalance := 100.00
// 	if actualUserBalance != expectedUserBalance {
// 		t.Fatalf("expected balance to be %f, but got %f", expectedUserBalance, actualUserBalance)
// 	}
// }

// // func TestSQLTransactionMoneyTransfer(t *testing.T) {
// // 	err := cleanupTestData()
// // 	assert.NoError(t, err)

// // 	// Inserting Test Data
// // 	err = insertTestUser("sender", "senderpassword", "123456789")
// // 	assert.NoError(t, err)
// // 	err = insertTestUser("beneficiary", "beneficiarypassword", "987654321")
// // 	assert.NoError(t, err)

// // 	err = insertTestUserBalance("sender", 100.0, "SGD", "SG")
// // 	assert.NoError(t, err)
// // 	err = insertTestUserBalance("beneficiary", 50.0, "SGD", "SG")
// // 	assert.NoError(t, err)

// // 	err = insertTestBeneficiary("beneficiary", "987654321", "SGD", 1)
// // 	assert.NoError(t, err)

// // 	testRepo := &repo{db: testDB}

// // 	// Mocking a transfer of $20
// // 	senderName := "sender"
// // 	beneficiaryName := "beneficiary"
// // 	amountTransferredCurrency := "SGD"
// // 	amountReceivedCurrency := "SGD"
// // 	confirmedStatus := TRANSACTION_STATUS.CONFIRMED
// // 	receivedStatus := TRANSACTION_STATUS.RECEIVED
// // 	amountTransferred := 20.0
// // 	amountReceived := 20.0

// // 	err = testRepo.SQLTransactionMoneyTransfer(context.Background(), senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, confirmedStatus, receivedStatus, amountTransferred, amountReceived)
// // 	if err != nil {
// // 		t.Fatalf("error calling SQLTransactionMoneyTransfer: %s", err)
// // 	}

// // 	// Test that data was inserted and updated properly for sender and beneficiary
// // 	actualSenderBalance, err := getTestUserBalance("sender")
// // 	assert.NoError(t, err)
// // 	actualBeneficiaryBalance, err := getTestUserBalance("beneficiary")
// // 	assert.NoError(t, err)

// // 	// Asserting that balance was updated correctly for both sender and beneficiary
// // 	expectedSenderBalance := 80.0
// // 	if actualSenderBalance != expectedSenderBalance {
// // 		t.Fatalf("expected balance to be %f, but got %f", expectedSenderBalance, actualSenderBalance)
// // 	}

// // 	expectedBeneficiaryBalance := 70.0
// // 	if actualBeneficiaryBalance != expectedBeneficiaryBalance {
// // 		t.Fatalf("expected balance to be %f, but got %f", expectedBeneficiaryBalance, actualBeneficiaryBalance)
// // 	}

// // 	// Asserting that the transactions were created
// // 	transactions, err := getTestUserTransactions("sender")
// // 	assert.NoError(t, err)
// // 	assert.Len(t, transactions, 1)

// // 	transactions, err = getTestUserTransactions("beneficiary")
// // 	assert.NoError(t, err)
// // 	assert.Len(t, transactions, 1)

// // }

// func TestGetByUserId(t *testing.T) {
// 	err := cleanupTestData()
// 	assert.NoError(t, err)

// 	// Insert test data
// 	err = insertTestUser("sender", "senderpassword", "123456789")
// 	assert.NoError(t, err)
// 	err = insertTestUser("beneficiary", "beneficiarypassword", "987654321")
// 	assert.NoError(t, err)

// 	err = insertTestUserBalance("sender", 100.0, "SGD", "SG")
// 	assert.NoError(t, err)
// 	err = insertTestUserBalance("beneficiary", 50.0, "SGD", "SG")
// 	assert.NoError(t, err)

// 	err = insertTestBeneficiary("beneficiary", "987654321", "SGD", 1)
// 	assert.NoError(t, err)

// 	// Inserting 3 transaction data, should return 3 transactions upon calling repository
// 	err = insertTestTransaction("sender", "beneficiary", 50.0, "SGD", 50.0, "SGD", "CONFIRMED")
// 	assert.NoError(t, err)
// 	err = insertTestTransaction("sender", "beneficiary", 50.0, "SGD", 50.0, "SGD", "CONFIRMED")
// 	assert.NoError(t, err)
// 	err = insertTestTransaction("sender", "beneficiary", 50.0, "SGD", 50.0, "SGD", "CONFIRMED")
// 	assert.NoError(t, err)

// 	testRepo := &repo{db: testDB}

// 	transactions, err := testRepo.GetByUserId(context.Background(), 1, 0, 0)
// 	assert.NoError(t, err)

// 	// Assert
// 	assert.Len(t, transactions.Transactions, 3)
// 	assert.Equal(t, "sender", transactions.Transactions[0].SenderUsername)
// 	assert.Equal(t, "beneficiary", transactions.Transactions[0].BeneficiaryUsername)
// }
