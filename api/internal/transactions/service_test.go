package transactions

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) ValidateFloatPrecision(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func TestCreateTransactions(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockRepo{}
	s := &service{
		repo: mockRepo,
	}

	// Test case 1: Zero amountTransferred
	err := s.CreateTransaction(ctx, "existing_user", "beneficiary_name", "beneficiary_number", "USD", "0.0")
	expectedError := "Please specify an amount to be transferred."
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 2: Negative amountTransferred
	err = s.CreateTransaction(ctx, "existing_user", "beneficiary_name", "beneficiary_number", "USD", "5.0")
	expectedError = "Minimum amount allowed per transfer is 10. Please try again."
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	err = s.CreateTransaction(ctx, "existing_user", "beneficiary_name", "beneficiary_number", "USD", "11000")
	expectedError = "Maximum amount allowed per transfer is 10000. Please try again."
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 2: Negative amountTransferred
	err = s.CreateTransaction(ctx, "existing_user", "beneficiary_name", "beneficiary_number", "USD", "-10.0")
	expectedError = "Minimum amount allowed per transfer is 10. Please try again."
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 3: Non-existing sender, doesn't exist in users table
	err = s.CreateTransaction(ctx, "non_existing_user", "beneficiary_name", "beneficiary_number", "USD", "50.0")
	expectedError = "This sender does not exist."
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 4: Duplicate senders, means there are duplicate usernames in `users` table
	err = s.CreateTransaction(ctx, "another_existing_user", "beneficiary_name", "beneficiary_number", "USD", "50.0")
	expectedError = "Duplicate senders."
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 5: sender is not linked to the specified beneficiary, prevent sender from transferring funds
	err = s.CreateTransaction(ctx, "existing_user", "beneficiary_name", "non_existing_beneficiary_number", "USD", "50.0")
	expectedError = "This user is not linked to the specified beneficiary"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 6: Incorrect beneficiary name provided but beneficiary mobile number exists in db
	err = s.CreateTransaction(ctx, "existing_user", "incorrect_beneficiary_name", "beneficiary_number", "USD", "50.0")
	expectedError = "Do you mean: existing_beneficiary ?"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 7: More than 1 beneficiary with the same name and mobile number
	err = s.CreateTransaction(ctx, "existing_user", "duplicate_beneficiary_name", "beneficiary_number", "USD", "50.0")
	expectedError = "Duplicate beneficiary"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %q but got %q", expectedError, err.Error())
	}

	// Test case 8: Successful transaction
	err = s.CreateTransaction(ctx, "existing_user", "beneficiary_name", "beneficiary_number", "USD", "50.0")
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err.Error())
	}

}

func TestGetTransactions(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockRepo{}
	s := &service{
		repo: mockRepo,
	}

	expectedTransactions := &Transactions{
		Transactions: []Transaction{
			{
				SenderName:                "Alice",
				BeneficiaryName:           "Bob",
				AmountTransferred:         100.0,
				AmountTransferredCurrency: "USD",
				AmountReceived:            90.0,
				AmountReceivedCurrency:    "EUR",
				Status:                    "COMPLETED",
				DateTransferred:           time.Now(),
				DateReceived:              time.Now(),
			},
			{
				SenderName:                "Alice",
				BeneficiaryName:           "David",
				AmountTransferred:         50.0,
				AmountTransferredCurrency: "GBP",
				AmountReceived:            60.0,
				AmountReceivedCurrency:    "CAD",
				Status:                    "COMPLETED",
				DateTransferred:           time.Now(),
				DateReceived:              time.Now(),
			},
		},
	}

	// Stub the repository method
	mockRepo.On("GetByUserId", ctx, mock.AnythingOfType("string")).Return(expectedTransactions, nil)

	actualTransactions, err := s.GetTransactions(ctx, "Alice")

	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, actualTransactions)

	mockRepo.AssertCalled(t, "GetByUserId", ctx, "Alice")
}
