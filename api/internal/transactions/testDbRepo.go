package transactions

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Stub the Transaction Repo
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetByUserId(ctx context.Context, username string) (*Transactions, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*Transactions), args.Error(1)
}

func (m *MockRepo) GetCountByUsername(ctx context.Context, username string) (int, error) {
	if username == "existing_user" || username == "zero_balance_user" {
		return 1, nil
	} else if username == "another_existing_user" {
		return 2, nil
	} else if username == "non_existing_user" {
		return 0, nil
	}
	return 0, nil
}

func (m *MockRepo) GetCountByUsernameAndBeneficiaryNumber(ctx context.Context, username, beneficiaryNumber string) (int, error) {
	if username == "existing_user" && beneficiaryNumber == "beneficiary_number" {
		return 1, nil
	} else if username == "zero_balance_user" && beneficiaryNumber == "beneficiary_number" {
		return 1, nil
	} else if username == "non_existing_user" && beneficiaryNumber == "beneficiary_number" {
		return 0, nil
	}
	return 0, nil
}

func (m *MockRepo) GetCountByBeneficiaryNameAndBeneficiaryNumber(ctx context.Context, beneficiaryName, beneficiaryNumber string) (int, error) {

	if beneficiaryName == "duplicate_beneficiary_name" && beneficiaryNumber == "beneficiary_number" {
		return 2, nil
	} else if beneficiaryName == "beneficiary_name" && beneficiaryNumber == "beneficiary_number" {
		return 1, nil
	} else if beneficiaryName == "existing_beneficiary" && beneficiaryNumber == "beneficiary_number" {
		return 1, nil
	} else if beneficiaryName == "external_beneficiary" && beneficiaryNumber == "beneficiary_number" {
		return 0, nil
	} else if beneficiaryName == "incorrect_beneficiary_name" && beneficiaryNumber == "beneficiary_number" {
		return 0, nil
	}
	return 0, nil
}

func (m *MockRepo) GetUsernameByBeneficiaryNumber(ctx context.Context, beneficiaryNumber string) (string, error) {
	if beneficiaryNumber == "beneficiary_number" {
		return "existing_beneficiary", nil
	}
	return "", nil
}

func (m *MockRepo) GetCurrencyByBeneficiaryMobileNumber(ctx context.Context, beneficiaryNumber string) (string, error) {
	if beneficiaryNumber == "beneficiary_number" {
		return "USD", nil
	}
	return "", nil
}

func (m *MockRepo) GetUserBalanceByUsername(ctx context.Context, username string) (float64, error) {
	if username == "existing_user" {
		return 100.0, nil
	} 
	return 50.0, nil
}

func (m *MockRepo) UpdateUserBalanceByUsername(ctx context.Context, finalAmount float64, username string) error {
	return nil
}

func (m *MockRepo) InsertIntoTransactions(ctx context.Context, username, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, status string, amountTransferred, amountReceived float64) error {
	return nil
}

func (m *MockRepo) SQLTransactionMoneyTransfer(ctx context.Context, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, confirmedStatus, receivedStatus string, amountTransferred, amountReceived float64) error {
	return nil
}
