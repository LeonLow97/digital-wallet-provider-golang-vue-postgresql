package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlers "github.com/LeonLow97/go-clean-architecture/delivery/http/handler"
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	mocks "github.com/LeonLow97/go-clean-architecture/mocks/usecase"
	"github.com/LeonLow97/go-clean-architecture/testdata"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWalletHandler_NewWalletHandler(t *testing.T) {
	t.Run("ReturnsAnInstanceOfWalletHandler", func(t *testing.T) {
		walletUsecase := &mocks.WalletUsecase{}
		instance := handlers.NewWalletHandler(walletUsecase)

		require.IsType(t, handlers.WalletHandler{}, *instance)
	})
}

func TestWalletHandler_GetWallet(t *testing.T) {
	wallet := testdata.MockWallet()

	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
		ExpectedErrorResponse     string
	}

	testCases := []testCase{
		{
			Title:                  "ReturnsSuccessfully",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallet: []interface{}{wallet, nil},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Title:                  "ReturnsError_NoWalletFound",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallet: []interface{}{nil, exception.ErrNoWalletFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"Wallet not found."}`,
		},
		{
			Title:                  "ReturnsError_NoWalletBalancesFound",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallet: []interface{}{nil, exception.ErrWalletBalancesNotFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"No wallet balances found. Please top up."}`,
		},
		{
			Title:                  "ReturnsError_InternalServerError",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallet: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus:        http.StatusInternalServerError,
			ExpectedErrorResponse: `{"status":500,"message":"Internal Server Error"}`,
		},
		{
			Title:                  "ReturnsError_MissingUserIDFromContext",
			GivenUserIDWithContext: 0,
			ExpectedStatus:         http.StatusUnauthorized,
			ExpectedErrorResponse:  `{"status":401,"message":"Unauthorized"}`,
		},
		{
			Title:                  "ReturnsError_MissingWalletIDFromParams",
			GivenUserIDWithContext: 1,
			ExpectedStatus:         http.StatusBadRequest,
			ExpectedErrorResponse:  `{"status":400,"message":"Bad Request"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			walletUsecase := &mocks.WalletUsecase{}

			walletUsecase.On("GetWallet", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletUsecaseReturnValues.GetWallet...)

			walletHandler := handlers.NewWalletHandler(walletUsecase)

			req, err := http.NewRequest(http.MethodGet, "/api/v1/wallet/1", nil)
			require.NoError(t, err)

			ctx := testdata.InjectUserIDIntoContext(req.Context(), tc.GivenUserIDWithContext)
			req = req.WithContext(ctx)

			// test for walletID in params
			if !strings.HasSuffix(tc.Title, "MissingWalletIDFromParams") {
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
			}

			rr := httptest.NewRecorder()

			walletHandler.GetWallet(rr, req)

			require.Equal(t, tc.ExpectedStatus, rr.Code)

			if strings.HasPrefix(tc.Title, "ReturnsError") {
				require.Equal(t, tc.ExpectedErrorResponse, rr.Body.String())
			}
		})
	}
}

func TestWalletHandler_GetWallets(t *testing.T) {
	wallets := []domain.Wallet{*testdata.MockWallet()}

	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
		ExpectedErrorResponse     string
	}

	testCases := []testCase{
		{
			Title:                  "ReturnsSuccessfully",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallets: []interface{}{&wallets, nil},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Title:                  "ReturnsError_NoWalletsFound",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallets: []interface{}{nil, exception.ErrNoWalletsFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"No wallets found."}`,
		},
		{
			Title:                  "ReturnsError_NoWalletBalancesFound",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallets: []interface{}{nil, exception.ErrWalletBalancesNotFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"No wallet balances found. Please top up."}`,
		},
		{
			Title:                  "ReturnsError_InternalServerError",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallets: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus:        http.StatusInternalServerError,
			ExpectedErrorResponse: `{"status":500,"message":"Internal Server Error"}`,
		},
		{
			Title:                  "ReturnsError_MissingUserIDFromContext",
			GivenUserIDWithContext: 0,
			ExpectedStatus:         http.StatusUnauthorized,
			ExpectedErrorResponse:  `{"status":401,"message":"Unauthorized"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			walletUsecase := &mocks.WalletUsecase{}

			walletUsecase.On("GetWallets", mock.Anything, mock.Anything).
				Return(tc.WalletUsecaseReturnValues.GetWallets...)

			walletHandler := handlers.NewWalletHandler(walletUsecase)

			req, err := http.NewRequest(http.MethodGet, "/api/v1/wallet/all", nil)
			require.NoError(t, err)

			ctx := testdata.InjectUserIDIntoContext(req.Context(), tc.GivenUserIDWithContext)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			walletHandler.GetWallets(rr, req)

			require.Equal(t, tc.ExpectedStatus, rr.Code)

			if strings.HasPrefix(tc.Title, "ReturnsError") {
				require.Equal(t, tc.ExpectedErrorResponse, rr.Body.String())
			}
		})
	}
}

func TestWalletHandler_GetWalletTypes(t *testing.T) {
	type testCase struct {
		Title                     string
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
		ExpectedErrorResponse     string
	}

	testCases := []testCase{
		{
			Title: "ReturnsSuccessfully",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWalletTypes: []interface{}{
					&[]dto.GetWalletTypesResponse{{ID: 1, WalletType: "personal"}}, nil,
				},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Title: "ReturnsError_NoWalletTypesFound",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWalletTypes: []interface{}{nil, exception.ErrWalletTypesNotFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"No wallet types found."}`,
		},
		{
			Title: "ReturnsError_InternalServerError",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWalletTypes: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus:        http.StatusInternalServerError,
			ExpectedErrorResponse: `{"status":500,"message":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			walletUsecase := &mocks.WalletUsecase{}

			walletUsecase.On("GetWalletTypes", mock.Anything).
				Return(tc.WalletUsecaseReturnValues.GetWalletTypes...)

			walletHandler := handlers.NewWalletHandler(walletUsecase)

			req, err := http.NewRequest(http.MethodGet, "/api/v1/wallet/types", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			walletHandler.GetWalletTypes(rr, req)

			require.Equal(t, tc.ExpectedStatus, rr.Code)

			if strings.HasPrefix(tc.Title, "ReturnsError") {
				require.Equal(t, tc.ExpectedErrorResponse, rr.Body.String())
			}
		})
	}
}

func TestWalletHandler_CreateWallet(t *testing.T) {
	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		GivenCreateWalletRequest  string
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
		ExpectedErrorResponse     string
	}

	testCases := []testCase{
		{
			Title:                    "ReturnsSuccessfully",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":1,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CreateWallet: []interface{}{nil},
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Title:                  "ReturnsError_MissingUserIDFromContext",
			GivenUserIDWithContext: 0,
			ExpectedStatus:         http.StatusUnauthorized,
			ExpectedErrorResponse:  `{"status":401,"message":"Unauthorized"}`,
		},
		{
			Title:                    "ReturnsError_ReadJSONBody",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: ``,
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Bad Request"}`,
		},
		{
			Title:                    "ReturnsError_WalletTypeID=0",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":0,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Key: 'CreateWalletRequest.WalletTypeID' Error:Field validation for 'WalletTypeID' failed on the 'required' tag"}`,
		},
		{
			Title:                    "ReturnsError_MissingCurrencyAmount",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":0}`,
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Key: 'CreateWalletRequest.WalletTypeID' Error:Field validation for 'WalletTypeID' failed on the 'required' tag Key: 'CreateWalletRequest.CurrencyAmount' Error:Field validation for 'CurrencyAmount' failed on the 'required' tag"}`,
		},
		{
			Title:                    "ReturnsError_BalanceNotFound",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":1,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CreateWallet: []interface{}{exception.ErrBalanceNotFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"Balance not found. Please deposit to create a new balance."}`,
		},
		{
			Title:                    "ReturnsError_WalletTypeInvalid",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":1,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CreateWallet: []interface{}{exception.ErrWalletTypeInvalid},
			},
			ExpectedStatus:        http.StatusBadRequest,
			ExpectedErrorResponse: `{"status":400,"message":"Wallet type is invalid. Please try another wallet type."}`,
		},
		{
			Title:                    "ReturnsError_WalletAlreadyExists",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":1,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CreateWallet: []interface{}{exception.ErrWalletAlreadyExists},
			},
			ExpectedStatus:        http.StatusBadRequest,
			ExpectedErrorResponse: `{"status":400,"message":"The wallet you are trying to create already exist. Please try again."}`,
		},
		{
			Title:                    "ReturnsError_InsufficientFunds",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":1,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CreateWallet: []interface{}{exception.ErrInsufficientFunds},
			},
			ExpectedStatus:        http.StatusBadRequest,
			ExpectedErrorResponse: `{"status":400,"message":"Insufficient funds in account. Please top up."}`,
		},
		{
			Title:                    "ReturnsError_InternalServerError",
			GivenUserIDWithContext:   1,
			GivenCreateWalletRequest: `{"wallet_type_id":1,"currency_amount":[{"amount":100.0,"currency":"USD"},{"amount":200.0,"currency":"EUR"}]}`,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CreateWallet: []interface{}{errors.New("internal server error")},
			},
			ExpectedStatus:        http.StatusInternalServerError,
			ExpectedErrorResponse: `{"status":500,"message":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			walletUsecase := &mocks.WalletUsecase{}

			walletUsecase.On("CreateWallet", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletUsecaseReturnValues.CreateWallet...)

			walletHandler := handlers.NewWalletHandler(walletUsecase)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(tc.GivenCreateWalletRequest))
			require.NoError(t, err)

			ctx := testdata.InjectUserIDIntoContext(req.Context(), tc.GivenUserIDWithContext)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			walletHandler.CreateWallet(rr, req)

			require.Equal(t, tc.ExpectedStatus, rr.Code)
			require.Equal(t, tc.ExpectedErrorResponse, rr.Body.String())
		})
	}
}

func TestWalletHandler_UpdateWallet(t *testing.T) {
	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		GivenUpdateWalletRequest  string
		GivenOperation            string
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
		ExpectedErrorResponse     string
	}

	testCases := []testCase{
		{
			Title:                    "ReturnsSuccessfully_TopUpWallet",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "topup",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				TopUpWallet: []interface{}{nil},
			},
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Title:                    "ReturnsSuccessfully_CashOutWallet",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "withdraw",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CashOutWallet: []interface{}{nil},
			},
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Title:                  "ReturnsError_MissingUserIDFromContext",
			GivenUserIDWithContext: 0,
			ExpectedStatus:         http.StatusUnauthorized,
			ExpectedErrorResponse:  `{"status":401,"message":"Unauthorized"}`,
		},
		{
			Title:                  "ReturnsError_MissingWalletIDFromParams",
			GivenUserIDWithContext: 1,
			ExpectedStatus:         http.StatusBadRequest,
			ExpectedErrorResponse:  `{"status":400,"message":"Bad Request"}`,
		},
		{
			Title:                  "ReturnsError_MissingOperationFromParams",
			GivenUserIDWithContext: 1,
			ExpectedStatus:         http.StatusBadRequest,
			ExpectedErrorResponse:  `{"status":400,"message":"Bad Request"}`,
		},
		{
			Title:                    "ReturnsError_EmptyJSONBody",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: ``,
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Bad Request"}`,
		},
		{
			Title:                    "ReturnsError_NegativeAmount",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":-100,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Key: 'UpdateWalletRequest.CurrencyAmount[0].Amount' Error:Field validation for 'Amount' failed on the 'gt' tag"}`,
		},
		{
			Title:                    "ReturnsError_InvalidCurrency",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":-100,"currency":"ABCDEF"},{"amount":20.00,"currency":"EUR"}]}`,
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Key: 'UpdateWalletRequest.CurrencyAmount[0].Amount' Error:Field validation for 'Amount' failed on the 'gt' tag Key: 'UpdateWalletRequest.CurrencyAmount[0].Currency' Error:Field validation for 'Currency' failed on the 'max' tag"}`,
		},
		{
			Title:                    "ReturnsError_InvalidOperationURLParams",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "exchange",
			ExpectedStatus:           http.StatusBadRequest,
			ExpectedErrorResponse:    `{"status":400,"message":"Bad Request"}`,
		},
		{
			Title:                    "ReturnsError_TopUpWallet_ErrNoWalletFound",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "topup",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				TopUpWallet: []interface{}{exception.ErrNoWalletFound},
			},
			ExpectedStatus:        http.StatusUnauthorized,
			ExpectedErrorResponse: `{"status":401,"message":"Unauthorized"}`,
		},
		{
			Title:                    "ReturnsError_TopUpWallet_ErrWalletBalancesNotFound",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "topup",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				TopUpWallet: []interface{}{exception.ErrWalletBalancesNotFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"No wallet balances found. Please top up."}`,
		},
		{
			Title:                    "ReturnsError_TopUpWallet_ErrBalanceNotFound",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "topup",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				TopUpWallet: []interface{}{exception.ErrBalanceNotFound},
			},
			ExpectedStatus:        http.StatusNotFound,
			ExpectedErrorResponse: `{"status":404,"message":"Balance not found. Please deposit to create a new balance."}`,
		},
		{
			Title:                    "ReturnsError_TopUpWallet_ErrInsufficientFunds",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "topup",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				TopUpWallet: []interface{}{exception.ErrInsufficientFunds},
			},
			ExpectedStatus:        http.StatusBadRequest,
			ExpectedErrorResponse: `{"status":400,"message":"Insufficient funds in account. Please top up."}`,
		},
		{
			Title:                    "ReturnsError_CashOutWallet_ErrInsufficientFundsForWithdrawal",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "withdraw",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CashOutWallet: []interface{}{exception.ErrInsufficientFundsForWithdrawal},
			},
			ExpectedStatus:        http.StatusBadRequest,
			ExpectedErrorResponse: `{"status":400,"message":"The specified amount for withdrawal is more than the current balance amount. Please lower the withdrawal amount."}`,
		},
		{
			Title:                    "ReturnsError_CashOutWallet_InternalServerError",
			GivenUserIDWithContext:   1,
			GivenUpdateWalletRequest: `{"currency_amount":[{"amount":10.50,"currency":"USD"},{"amount":20.00,"currency":"EUR"}]}`,
			GivenOperation:           "withdraw",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				CashOutWallet: []interface{}{errors.New("internal server error")},
			},
			ExpectedStatus:        http.StatusInternalServerError,
			ExpectedErrorResponse: `{"status":500,"message":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			walletUsecase := &mocks.WalletUsecase{}

			walletUsecase.On("TopUpWallet", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletUsecaseReturnValues.TopUpWallet...)
			walletUsecase.On("CashOutWallet", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletUsecaseReturnValues.CashOutWallet...)

			walletHandler := handlers.NewWalletHandler(walletUsecase)

			req, err := http.NewRequest(http.MethodPut, "/api/v1/wallet/update/1/topup", strings.NewReader(tc.GivenUpdateWalletRequest))
			require.NoError(t, err)

			// Set URL variables for testing as the mux router would
			vars := make(map[string]string)

			// test for walletID in params
			if !strings.HasSuffix(tc.Title, "MissingWalletIDFromParams") {
				vars["id"] = "1"
			}

			// test for operation in params
			if !strings.HasSuffix(tc.Title, "MissingOperationFromParams") {
				vars["operation"] = tc.GivenOperation
			}

			if len(vars) > 0 {
				req = mux.SetURLVars(req, vars)
			}

			ctx := testdata.InjectUserIDIntoContext(req.Context(), tc.GivenUserIDWithContext)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			walletHandler.UpdateWallet(rr, req)

			require.Equal(t, tc.ExpectedStatus, rr.Code)
			require.Equal(t, tc.ExpectedErrorResponse, rr.Body.String())
		})
	}
}
