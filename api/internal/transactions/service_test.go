package transactions

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetTransactions_Service(t *testing.T) {
	testCases := []struct {
		Test                                  string
		UserId                                int
		Page                                  int
		PageSize                              int
		Offset                                int
		ExpectedTransactions                  *Transactions
		ExpectedTotalPages                    int
		ExpectedIsLastPage                    bool
		ExpectErr                             bool
		MockErrorGetTransactionsCountByUserId error
		MockErrorGetTransactionsByUserId      error
	}{
		{
			Test:     "Successfully Get Transactions",
			UserId:   1,
			Page:     1,
			PageSize: 20,
			Offset:   0,
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName: "TestFirstName",
					},
				},
			},
			ExpectedTotalPages:                    50,
			ExpectedIsLastPage:                    false,
			ExpectErr:                             false,
			MockErrorGetTransactionsCountByUserId: nil,
			MockErrorGetTransactionsByUserId:      nil,
		},
		{
			Test:     "Page is negative, should reassign page = 1",
			UserId:   1,
			Page:     -10,
			PageSize: 20,
			Offset:   0,
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName: "TestFirstName",
					},
				},
			},
			ExpectedTotalPages:                    50,
			ExpectedIsLastPage:                    false,
			ExpectErr:                             false,
			MockErrorGetTransactionsCountByUserId: nil,
			MockErrorGetTransactionsByUserId:      nil,
		},
		{
			Test:     "Last Page",
			UserId:   2,
			Page:     100,
			PageSize: 30,
			Offset:   990,
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName: "TestFirstName",
					},
				},
			},
			ExpectedTotalPages:                    34,
			ExpectedIsLastPage:                    true,
			ExpectErr:                             false,
			MockErrorGetTransactionsCountByUserId: nil,
			MockErrorGetTransactionsByUserId:      nil,
		},
		{
			Test:     "Error in GetTransactionsCountByUserId",
			UserId:   513,
			Page:     1,
			PageSize: 20,
			Offset:   0,
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName: "TestFirstName",
					},
				},
			},
			ExpectedTotalPages:                    50,
			ExpectedIsLastPage:                    false,
			ExpectErr:                             true,
			MockErrorGetTransactionsCountByUserId: errors.New("Internal Server Error!"),
			MockErrorGetTransactionsByUserId:      nil,
		},
		{
			Test:     "Error in GetTransactionsByUserId",
			UserId:   514,
			Page:     1,
			PageSize: 20,
			Offset:   0,
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName: "TestFirstName",
					},
				},
			},
			ExpectedTotalPages:                    50,
			ExpectedIsLastPage:                    false,
			ExpectErr:                             true,
			MockErrorGetTransactionsCountByUserId: nil,
			MockErrorGetTransactionsByUserId:      errors.New("Internal Server Error!"),
		},
	}

	mockRepo := mockRepo{}
	s := NewService(&mockRepo)

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			mockRepo.On("GetTransactionsCountByUserId", mock.Anything, tc.UserId).Return(1000, tc.MockErrorGetTransactionsCountByUserId)
			mockRepo.On("GetTransactionsByUserId", mock.Anything, tc.UserId, tc.PageSize, tc.Offset).Return(tc.ExpectedTransactions, tc.MockErrorGetTransactionsByUserId)

			transactions, totalPages, isLastPage, err := s.GetTransactions(context.Background(), tc.UserId, tc.Page, tc.PageSize)

			if !tc.ExpectErr {
				require.NoError(t, err)
				assert.Equal(t, totalPages, tc.ExpectedTotalPages, "totalPages")
				assert.Equal(t, isLastPage, tc.ExpectedIsLastPage, "isLastPage")
				assert.Equal(t, fmt.Sprintf("%+v", transactions), fmt.Sprintf("%+v", tc.ExpectedTransactions))
			} else {
				require.Error(t, err)
				assert.Equal(t, totalPages, 0, "totalPages")
				assert.Equal(t, isLastPage, false, "isLastPage")
				assert.Nil(t, transactions, "transactions")
			}
		})
	}
}

func Test_CreateTransaction_Service(t *testing.T) {
	testCases := []struct {
		Test                                               string
		UserId                                             int
		BeneficiaryId                                      int
		Transaction                                        *CreateTransaction
		SenderTransaction                                  TransactionEntity
		BeneficiaryTransaction                             TransactionEntity
		ExpectedUserCount                                  int
		ExpectedIsLinked                                   int
		ExpectedBeneficiaryHasTransferredCurrency          int
		ExpectedBeneficiaryBalanceId                       int
		ExpectedBeneficiaryCurrency                        string
		ExpectedSenderHasTransferredCurrency               int
		ExpectedSenderBalanceId                            int
		ExpectedErrorMessage                               string
		MockErrorGetUserCountByUserId                      error
		MockErrorGetUserIdByMobileNumber                   error
		MockErrorGetCountByUserIdAndBeneficiaryId          error
		MockErrorGetCountByUserIdAndCurrencyForBeneficiary error
		MockErrorGetBalanceIdByUserIdAndPrimary            error
		MockErrorGetCountByUserIdAndCurrencyForSender      error
		MockErrorCreateTransactionSQLTransaction           error
		ExpectErr                                          bool
	}{
		{
			Test:   "TransferredAmount is 0",
			UserId: 1,
			Transaction: &CreateTransaction{
				TransferredAmount: 0,
			},
			ExpectedErrorMessage: "Transfer amount must be between $10 and $10,000.",
			ExpectErr:            true,
		},
		{
			Test:   "TransferredAmount is negative",
			UserId: 1,
			Transaction: &CreateTransaction{
				TransferredAmount: -1234,
			},
			ExpectedErrorMessage: "Transfer amount must be between $10 and $10,000.",
			ExpectErr:            true,
		},
		{
			Test:   "TransferredAmount is less than 10",
			UserId: 1,
			Transaction: &CreateTransaction{
				TransferredAmount: 9,
			},
			ExpectedErrorMessage: "Transfer amount must be between $10 and $10,000.",
			ExpectErr:            true,
		},
		{
			Test:   "TransferredAmount is more than 10,000",
			UserId: 1,
			Transaction: &CreateTransaction{
				TransferredAmount: 60000,
			},
			ExpectedErrorMessage: "Transfer amount must be between $10 and $10,000.",
			ExpectErr:            true,
		},
		{
			Test:   "TransferredAmount is more than 10,000",
			UserId: 1,
			Transaction: &CreateTransaction{
				TransferredAmount: 60000,
			},
			ExpectedErrorMessage: "Transfer amount must be between $10 and $10,000.",
			ExpectErr:            true,
		},
		{
			Test:   "TransferredAmount has more than 2 decimal places",
			UserId: 1,
			Transaction: &CreateTransaction{
				TransferredAmount: 123.831590843059894158121342314,
			},
			ExpectedErrorMessage: "Transfer amount must be up to 2 decimal places.",
			ExpectErr:            true,
		},
		{
			Test:          "InternalServerError in GetUserCountByUserId",
			UserId:        2,
			BeneficiaryId: 3,
			Transaction: &CreateTransaction{
				TransferredAmount: 9000,
			},
			ExpectedUserCount:             1,
			ExpectErr:                     true,
			MockErrorGetUserCountByUserId: errors.New("Internal Server Error"),
			ExpectedErrorMessage:          "Internal Server Error",
		},
		{
			Test:          "BadRequest in GetUserCountByUserId",
			UserId:        3,
			BeneficiaryId: 4,
			Transaction: &CreateTransaction{
				TransferredAmount: 9000,
			},
			ExpectedUserCount:    0,
			ExpectErr:            true,
			ExpectedErrorMessage: "The specified sender does not exist.",
		},
		{
			Test:          "InternalServerError in GetUserIdByMobileNumber",
			UserId:        4,
			BeneficiaryId: 5,
			Transaction: &CreateTransaction{
				BeneficiaryNumber: "+65 98765432",
				TransferredAmount: 9000,
			},
			ExpectedUserCount:                1,
			ExpectErr:                        true,
			MockErrorGetUserIdByMobileNumber: errors.New("Internal Server Error"),
			ExpectedErrorMessage:             "Internal Server Error",
		},
		{
			Test:          "BadRequest (beneficiaryId == 0) in GetUserIdByMobileNumber",
			UserId:        5,
			BeneficiaryId: 0,
			Transaction: &CreateTransaction{
				TransferredAmount: 9000,
				BeneficiaryNumber: "0",
			},
			ExpectedUserCount:    1,
			ExpectErr:            true,
			ExpectedErrorMessage: "The specified beneficiary does not exist.",
		},
		{
			Test:          "BadRequest (userId == beneficiaryId) in GetUserIdByMobileNumber",
			UserId:        5,
			BeneficiaryId: 5,
			Transaction: &CreateTransaction{
				TransferredAmount: 9000,
				BeneficiaryNumber: "01",
			},
			ExpectedUserCount: 1,
			ExpectedIsLinked:  1,
			ExpectErr:         true,
			MockErrorGetCountByUserIdAndBeneficiaryId: errors.New("Internal Server Error"),
			ExpectedErrorMessage:                      "Unable to send money to yourself.",
		},
		{
			Test:          "InternalServerError in GetCountByUserIdAndBeneficiaryId",
			UserId:        5,
			BeneficiaryId: 6,
			Transaction: &CreateTransaction{
				TransferredAmount: 9000,
				BeneficiaryNumber: "02",
			},
			ExpectedUserCount: 1,
			ExpectedIsLinked:  1,
			ExpectErr:         true,
			MockErrorGetCountByUserIdAndBeneficiaryId: errors.New("Internal Server Error"),
			ExpectedErrorMessage:                      "Internal Server Error",
		},
		{
			Test:          "InternalServerError in GetCountByUserIdAndBeneficiaryId",
			UserId:        6,
			BeneficiaryId: 7,
			Transaction: &CreateTransaction{
				TransferredAmount: 9000,
				BeneficiaryNumber: "03",
			},
			ExpectedUserCount:    1,
			ExpectedIsLinked:     0,
			ExpectErr:            true,
			ExpectedErrorMessage: "Unable to transfer funds. Sender is not linked to the specified beneficiary.",
		},
		{
			Test:          "InternalServerError in GetCountByUserIdAndCurrency",
			UserId:        7,
			BeneficiaryId: 8,
			Transaction: &CreateTransaction{
				TransferredAmount:         9000,
				BeneficiaryNumber:         "04",
				TransferredAmountCurrency: "SGD",
			},
			ExpectedUserCount:                         1,
			ExpectedIsLinked:                          1,
			ExpectedBeneficiaryHasTransferredCurrency: 0,
			ExpectedBeneficiaryBalanceId:              0,
			ExpectErr:                                 true,
			MockErrorGetCountByUserIdAndCurrencyForBeneficiary: errors.New("Internal Server Error"),
			ExpectedErrorMessage: "Internal Server Error",
		},
		{
			Test:          "InternalServerError in GetBalanceIdByUserIdAndPrimary",
			UserId:        8,
			BeneficiaryId: 9,
			Transaction: &CreateTransaction{
				TransferredAmount:         9000,
				BeneficiaryNumber:         "05",
				TransferredAmountCurrency: "SGD",
			},
			ExpectedUserCount:                         1,
			ExpectedIsLinked:                          1,
			ExpectedBeneficiaryHasTransferredCurrency: 0,
			ExpectedBeneficiaryBalanceId:              1,
			ExpectedBeneficiaryCurrency:               "USD",
			ExpectErr:                                 true,
			MockErrorGetBalanceIdByUserIdAndPrimary:   errors.New("Internal Server Error"),
			ExpectedErrorMessage:                      "Internal Server Error",
		},
		{
			Test:          "InternalServerError in GetCountByUserIdAndCurrency",
			UserId:        10,
			BeneficiaryId: 11,
			Transaction: &CreateTransaction{
				TransferredAmount:         9000,
				BeneficiaryNumber:         "06",
				TransferredAmountCurrency: "SGD",
			},
			ExpectedUserCount:                             1,
			ExpectedIsLinked:                              1,
			ExpectedBeneficiaryHasTransferredCurrency:     1,
			ExpectedBeneficiaryBalanceId:                  1,
			ExpectedBeneficiaryCurrency:                   "USD",
			ExpectedSenderHasTransferredCurrency:          1,
			ExpectedSenderBalanceId:                       2,
			ExpectErr:                                     true,
			MockErrorGetCountByUserIdAndCurrencyForSender: errors.New("Internal Server Error"),
			ExpectedErrorMessage:                          "Internal Server Error",
		},
		{
			Test:          "BadRequest in GetCountByUserIdAndCurrency",
			UserId:        11,
			BeneficiaryId: 12,
			Transaction: &CreateTransaction{
				TransferredAmount:         9000,
				BeneficiaryNumber:         "07",
				TransferredAmountCurrency: "YEN",
			},
			ExpectedUserCount:                         1,
			ExpectedIsLinked:                          1,
			ExpectedBeneficiaryHasTransferredCurrency: 1,
			ExpectedBeneficiaryBalanceId:              1,
			ExpectedBeneficiaryCurrency:               "USD",
			ExpectedSenderHasTransferredCurrency:      0,
			ExpectedSenderBalanceId:                   2,
			ExpectErr:                                 true,
			ExpectedErrorMessage:                      "You do not have balance in the specified currency. Please use another currency.",
		},
		{
			Test:          "InternalServerError in CreateTransactionSQLTransaction",
			UserId:        12,
			BeneficiaryId: 13,
			Transaction: &CreateTransaction{
				TransferredAmount:         9000,
				BeneficiaryNumber:         "08",
				TransferredAmountCurrency: "YEN",
			},
			ExpectedUserCount:                         1,
			ExpectedIsLinked:                          1,
			ExpectedBeneficiaryHasTransferredCurrency: 1,
			ExpectedBeneficiaryBalanceId:              1,
			ExpectedBeneficiaryCurrency:               "USD",
			ExpectedSenderHasTransferredCurrency:      0,
			ExpectedSenderBalanceId:                   2,
			ExpectErr:                                 false,
		},
	}

	mockRepo := mockRepo{}
	s := NewService(&mockRepo)

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			mockRepo.On("GetUserCountByUserId", mock.Anything, tc.UserId).Return(tc.ExpectedUserCount, tc.MockErrorGetUserCountByUserId)
			mockRepo.On("GetUserIdByMobileNumber", mock.Anything, tc.Transaction.BeneficiaryNumber).Return(tc.BeneficiaryId, tc.MockErrorGetUserIdByMobileNumber)
			mockRepo.On("GetCountByUserIdAndBeneficiaryId", mock.Anything, tc.UserId, tc.BeneficiaryId).Return(tc.ExpectedIsLinked, tc.MockErrorGetCountByUserIdAndBeneficiaryId)
			mockRepo.On("GetCountByUserIdAndCurrency", mock.Anything, tc.BeneficiaryId, tc.Transaction.TransferredAmountCurrency).Return(tc.ExpectedBeneficiaryHasTransferredCurrency, tc.ExpectedBeneficiaryBalanceId, tc.MockErrorGetCountByUserIdAndCurrencyForBeneficiary)
			mockRepo.On("GetBalanceIdByUserIdAndPrimary", mock.Anything, tc.BeneficiaryId).Return(tc.ExpectedBeneficiaryBalanceId, tc.ExpectedBeneficiaryCurrency, tc.MockErrorGetBalanceIdByUserIdAndPrimary)
			mockRepo.On("GetCountByUserIdAndCurrency", mock.Anything, tc.UserId, tc.Transaction.TransferredAmountCurrency).Return(tc.ExpectedSenderHasTransferredCurrency, tc.ExpectedSenderBalanceId, tc.MockErrorGetCountByUserIdAndCurrencyForSender)
			mockRepo.On("CreateTransactionSQLTransaction", mock.Anything, mock.Anything, mock.Anything).Return(tc.MockErrorCreateTransactionSQLTransaction)

			err := s.CreateTransaction(context.Background(), tc.UserId, tc.Transaction)

			if !tc.ExpectErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, err.Error(), tc.ExpectedErrorMessage)

			}
		})
	}
}
