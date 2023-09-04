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
	s, err := NewService(&mockRepo)
	require.NoError(t, err, "getting service with mock repo in GetTransactions service")

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
