package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeonLow97/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

// Stub the Transaction Service
type mockService struct {
	mock.Mock
}

func (s *mockService) GetTransactions(ctx context.Context, userId, page, pageSize int) (*Transactions, int, bool, error) {
	args := s.Called(ctx, userId, page, pageSize)
	return args.Get(0).(*Transactions), args.Int(1), args.Bool(2), args.Error(3)
}

func (s *mockService) CreateTransaction(ctx context.Context, username, BeneficiaryName, BeneficiaryNumber, AmountTransferredCurrency, AmountTransferred string) error {
	return nil
}

func createTestContextWithUserID(userId int) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.ContextUserIdKey, userId)
	return ctx
}

func Test_GetTransactions_Handler(t *testing.T) {
	testCases := []struct {
		Test                 string
		Page                 interface{}
		PageSize             interface{}
		Endpoint             string
		ExpectedTransactions *Transactions
		ExpectErr            bool
		MockError            error
		ExpectedJSONResponse string
		ExpectedStatusCode   int
	}{
		{
			Test:     "Successful Get Transactions",
			Page:     5,
			PageSize: 20,
			Endpoint: "/transactions?page=5&pageSize=20",
			ExpectedTransactions: &Transactions{
				Transactions: []Transaction{
					{
						SenderFirstName: "testFirstName",
					},
				},
			},
			ExpectErr:            false,
			MockError:            nil,
			ExpectedJSONResponse: `{"isLastPage":false,"total_pages":10,"transactions":{"transactions":[{"beneficiary_first_name":"","beneficiary_last_name":"","beneficiary_username":"","received_amount":0,"received_amount_currency":"","received_date":"0001-01-01T00:00:00Z","sender_first_name":"testFirstName","sender_last_name":"","sender_username":"","status":"","transferred_amount":0,"transferred_amount_currency":"","transferred_date":"0001-01-01T00:00:00Z"}]}}`,
			ExpectedStatusCode:   http.StatusOK,
		},
		{
			Test:                 "Error in GetTransactions service",
			Page:                 6,
			PageSize:             21,
			Endpoint:             "/transactions?page=6&pageSize=21",
			ExpectedTransactions: nil,
			ExpectErr:            true,
			MockError:            utils.InternalServerError{Message: "error in GetTransactions service"},
			ExpectedJSONResponse: `{"error":true,"message":"Internal Server Error"}`,
			ExpectedStatusCode:   http.StatusInternalServerError,
		},
		{
			Test:                 "Invalid URL Query Params for page",
			Page:                 "invalid",
			PageSize:             20,
			Endpoint:             "/transactions?page=invalid&pageSize=20",
			ExpectedTransactions: nil,
			ExpectErr:            true,
			MockError:            errors.New("invalid URL params supplied"),
			ExpectedJSONResponse: `{"error":true,"message":"Please provide numerical page in URL Param."}`,
			ExpectedStatusCode:   http.StatusBadRequest,
		},
		{
			Test:                 "Invalid URL Query Params for page size",
			Page:                 5,
			PageSize:             "invalid",
			Endpoint:             "/transactions?page=5&pageSize=invalid",
			ExpectedTransactions: nil,
			ExpectErr:            true,
			MockError:            errors.New("invalid URL params supplied"),
			ExpectedJSONResponse: `{"error":true,"message":"Please provide numerical page size in URL Param."}`,
			ExpectedStatusCode:   http.StatusBadRequest,
		},
	}

	// creating the mock service
	mockService := mockService{}
	transactionHandler, err := NewTransactionHandler(&mockService)
	require.NoError(t, err, "getting transactionHandler with mockService in Transaction handler")

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// create a mock GET request
			req, err := http.NewRequest(http.MethodGet, tc.Endpoint, nil)
			require.NoError(t, err)

			if !tc.ExpectErr {
				mockService.On("GetTransactions", mock.Anything, 1, tc.Page, tc.PageSize).Return(tc.ExpectedTransactions, 10, false, tc.MockError)
			} else {
				mockService.On("GetTransactions", mock.Anything, 1, tc.Page, tc.PageSize).Return(tc.ExpectedTransactions, 10, false, tc.MockError)
			}

			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "application/json")

			// set userId int request context value
			req = req.WithContext(createTestContextWithUserID(1))

			// calling the handler method
			handler := http.HandlerFunc(transactionHandler.GetTransactions)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.ExpectedStatusCode, rr.Code)

			// parse the jsonResponse
			var response envelope
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Error parsing JSON response: %v", err)
			}

			jsonData, _ := json.Marshal(response)

			if !tc.ExpectErr {
				assert.Equal(t, tc.ExpectedJSONResponse, string(jsonData))
			} else {
				assert.Equal(t, tc.ExpectedJSONResponse, string(jsonData))
			}
		})
	}
}
