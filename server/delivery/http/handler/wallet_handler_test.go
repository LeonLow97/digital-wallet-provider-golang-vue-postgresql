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
	wallet := testdata.Wallet()

	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
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
		},
		{
			Title:                  "ReturnsError_InternalServerError",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallet: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
		{
			Title:                  "ReturnsError_MissingUserIDFromContext",
			GivenUserIDWithContext: 0,
			ExpectedStatus:         http.StatusUnauthorized,
		},
		{
			Title:                  "ReturnsError_MissingWalletIDFromParams",
			GivenUserIDWithContext: 1,
			ExpectedStatus:         http.StatusBadRequest,
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

		})
	}
}

func TestWalletHandler_GetWallets(t *testing.T) {
	wallets := []domain.Wallet{*testdata.Wallet()}

	type testCase struct {
		Title                     string
		GivenUserIDWithContext    int
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
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
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Title:                  "ReturnsError_InternalServerError",
			GivenUserIDWithContext: 1,
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWallets: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
		{
			Title:                  "ReturnsError_MissingUserIDFromContext",
			GivenUserIDWithContext: 0,
			ExpectedStatus:         http.StatusUnauthorized,
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
		})
	}
}

func TestWalletHandler_GetWalletTypes(t *testing.T) {
	type testCase struct {
		Title                     string
		WalletUsecaseReturnValues mocks.WalletUsecaseReturnValues
		ExpectedStatus            int
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
			Title: "ReturnsError_InternalServerError",
			WalletUsecaseReturnValues: mocks.WalletUsecaseReturnValues{
				GetWalletTypes: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedStatus: http.StatusInternalServerError,
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
		})
	}
}
