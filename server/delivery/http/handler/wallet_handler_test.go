package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/LeonLow97/go-clean-architecture/delivery/http/handler"
	"github.com/LeonLow97/go-clean-architecture/exception"
	mocks "github.com/LeonLow97/go-clean-architecture/mocks/usecase"
	"github.com/LeonLow97/go-clean-architecture/testdata"
	"github.com/LeonLow97/go-clean-architecture/utils"
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
	wallet := testdata.Wallet()

	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
		ExpectedError             error
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
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  exception.ErrNoWalletFound,
		},
		{
			Title:                  "ReturnsError_InternalServerError",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallet: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.New("internal server error"),
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

			ctx := context.WithValue(req.Context(), utils.UserIDKey, tc.GivenUserIDWithContext)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			walletHandler.GetWallet(rr, req)

			require.Equal(t, tc.ExpectedStatus, rr.Code)
		})
	}
}
