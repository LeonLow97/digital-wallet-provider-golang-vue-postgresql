package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	mocks "github.com/LeonLow97/go-clean-architecture/mocks/repository"
	"github.com/LeonLow97/go-clean-architecture/testdata"
	"github.com/LeonLow97/go-clean-architecture/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWalletUsecase_NewWalletUsecase(t *testing.T) {
	dbConn := &sqlx.DB{}
	walletRepo := new(mocks.WalletRepository)
	balanceRepo := new(mocks.BalanceRepository)

	walletUsecase := usecase.NewWalletUsecase(dbConn, walletRepo, balanceRepo)
	require.NotNil(t, walletUsecase)

	// Use require.Implements to check if walletUsecase implements WalletUsecase interface
	require.Implements(t, (*domain.WalletUsecase)(nil), walletUsecase)
}

func TestWalletUsecase_GetWallet(t *testing.T) {
	testCases := []struct {
		Title                        string
		GivenUserID                  int
		GivenWalletID                int
		WalletRepositoryReturnValues mocks.WalletRepositoryReturnValues
		ExpectedWallet               *domain.Wallet
		ExpectedError                error
	}{
		{
			Title:         "ReturnsSuccessfully",
			GivenUserID:   1,
			GivenWalletID: 1,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{testdata.MockWallet(), nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.MockWalletCurrencyAmounts(), nil},
			},
			ExpectedWallet: testdata.MockWalletWithCurrencyAmounts(),
		},
		{
			Title:         "ReturnsError_GetWalletByWalletID",
			GivenUserID:   1,
			GivenWalletID: 1,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:         "ReturnsError_GetWalletBalancesByUserIDAndWalletID",
			GivenUserID:   1,
			GivenWalletID: 1,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{testdata.MockWallet(), nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			dbConn := &sqlx.DB{}
			walletRepo := new(mocks.WalletRepository)
			balanceRepo := new(mocks.BalanceRepository)

			walletUsecase := usecase.NewWalletUsecase(dbConn, walletRepo, balanceRepo)
			require.NotNil(t, walletUsecase)

			walletRepo.On("GetWalletByWalletID", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletByWalletID...)
			walletRepo.On("GetWalletBalancesByUserIDAndWalletID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletBalancesByUserIDAndWalletID...)

			wallet, err := walletUsecase.GetWallet(context.Background(), tc.GivenUserID, tc.GivenWalletID)

			if !strings.HasPrefix(tc.Title, "ReturnsError") {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedWallet, wallet)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError.Error(), err.Error())
				require.Nil(t, wallet)
			}
		})
	}
}

func TestWalletUsecase_GetWallets(t *testing.T) {
	walletBalances := []domain.WalletCurrencyAmount{}
	walletBalances = append(walletBalances, testdata.MockWalletCurrencyAmountsByWalletID(1)...)
	walletBalances = append(walletBalances, testdata.MockWalletCurrencyAmountsByWalletID(2)...)
	walletBalances = append(walletBalances, testdata.MockWalletCurrencyAmountsByWalletID(3)...)

	testCases := []struct {
		Title                        string
		GivenUserID                  int
		WalletRepositoryReturnValues mocks.WalletRepositoryReturnValues
		ExpectedWallets              *[]domain.Wallet
		ExpectedError                error
	}{
		{
			Title:       "ReturnsSuccessfully",
			GivenUserID: 1,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWallets:                []interface{}{testdata.MockWallets(3), nil},
				GetWalletBalancesByUserID: []interface{}{walletBalances, nil},
			},
			ExpectedWallets: &[]domain.Wallet{
				{ID: 1, WalletType: "Personal", WalletTypeID: 1, UserID: 1, CreatedAt: "2024-06-15T13:45:00Z", CurrencyAmount: testdata.MockWalletCurrencyAmountsByWalletID(1)},
				{ID: 2, WalletType: "Personal", WalletTypeID: 1, UserID: 1, CreatedAt: "2024-06-15T13:45:00Z", CurrencyAmount: testdata.MockWalletCurrencyAmountsByWalletID(2)},
				{ID: 3, WalletType: "Personal", WalletTypeID: 1, UserID: 1, CreatedAt: "2024-06-15T13:45:00Z", CurrencyAmount: testdata.MockWalletCurrencyAmountsByWalletID(3)},
			},
		},
		{
			Title:       "ReturnsError_GetWalletsInternalServerError",
			GivenUserID: 1,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWallets: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:       "ReturnsError_GetWalletBalancesByUserIDInternalServerError",
			GivenUserID: 1,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWallets:                []interface{}{testdata.MockWallets(3), nil},
				GetWalletBalancesByUserID: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			dbConn := &sqlx.DB{}
			walletRepo := new(mocks.WalletRepository)
			balanceRepo := new(mocks.BalanceRepository)

			walletUsecase := usecase.NewWalletUsecase(dbConn, walletRepo, balanceRepo)
			require.NotNil(t, walletUsecase)

			walletRepo.On("GetWallets", mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWallets...)
			walletRepo.On("GetWalletBalancesByUserID", mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletBalancesByUserID...)

			wallets, err := walletUsecase.GetWallets(context.Background(), tc.GivenUserID)

			if !strings.HasPrefix(tc.Title, "ReturnsError") {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedWallets, wallets)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError.Error(), err.Error())
				require.Nil(t, wallets)
			}
		})
	}
}

func TestWalletUsecase_GetWalletTypes(t *testing.T) {
	testCases := []struct {
		Title                        string
		WalletRepositoryReturnValues mocks.WalletRepositoryReturnValues
		ExpectedWalletTypes          *[]dto.GetWalletTypesResponse
		ExpectedError                error
	}{
		{
			Title: "ReturnsSuccessfully",
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletTypes: []interface{}{testdata.MockGetWalletTypesResponses(), nil},
			},
			ExpectedWalletTypes: testdata.MockGetWalletTypesResponses(),
		},
		{
			Title: "ReturnsError_GetWalletTypesInternalServerError",
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletTypes: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			dbConn := &sqlx.DB{}
			walletRepo := new(mocks.WalletRepository)
			balanceRepo := new(mocks.BalanceRepository)

			walletUsecase := usecase.NewWalletUsecase(dbConn, walletRepo, balanceRepo)
			require.NotNil(t, walletUsecase)

			walletRepo.On("GetWalletTypes", mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletTypes...)

			resp, err := walletUsecase.GetWalletTypes(context.Background())

			if !strings.HasPrefix(tc.Title, "ReturnsError") {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedWalletTypes, resp)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError.Error(), err.Error())
				require.Nil(t, resp)
			}
		})
	}
}

func TestWalletUsecase_CreateWallet(t *testing.T) {
	basicRequest := &dto.CreateWalletRequest{
		WalletTypeID: 1,
		CurrencyAmount: []dto.CurrencyAmount{
			{Amount: 100, Currency: "USD"},
			{Amount: 200, Currency: "SGD"},
		},
	}

	testCases := []struct {
		Title                         string
		GivenCreateWalletRequest      *dto.CreateWalletRequest
		WalletRepositoryReturnValues  mocks.WalletRepositoryReturnValues
		BalanceRepositoryReturnValues mocks.BalanceRepositoryReturnValues
		ExpectedError                 error
	}{
		{
			Title:                    "ReturnsSuccessfully",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
				CreateWallet:                    []interface{}{1, nil},
				InsertWalletCurrencyAmount:      []interface{}{nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
		},
		{
			Title:                    "ReturnsError_CheckWalletExistsByWalletTypeID_InternalServerError",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_CheckWalletExistsByWalletTypeID_WalletAlreadyExists",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{true, nil},
			},
			ExpectedError: exception.ErrWalletAlreadyExists,
		},
		{
			Title:                    "ReturnsError_CheckWalletTypeExists_InternalServerError",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{false, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_CheckWalletTypeExists_WalletTypeDoesNotExist",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{false, nil},
			},
			ExpectedError: exception.ErrWalletTypeInvalid,
		},
		{
			Title:                    "ReturnsError_GetBalances_InternalServerError",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_ErrBalanceNotFound_ForOneCurrency",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{[]domain.Balance{
					{ID: 1, Balance: 500.0, Currency: "EUR", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
				}, nil},
			},
			ExpectedError: exception.ErrBalanceNotFound,
		},
		{
			Title:                    "ReturnsError_ErrInsufficientFunds_ForOneCurrency",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{[]domain.Balance{
					{ID: 1, Balance: 1.0, Currency: "USD", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
					{ID: 3, Balance: 200.0, Currency: "SGD", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
				}, nil},
			},
			ExpectedError: exception.ErrInsufficientFunds,
		},
		{
			Title:                    "ReturnsError_ErrInsufficientFunds_ForAllCurrencies",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{[]domain.Balance{
					{ID: 1, Balance: 1.0, Currency: "USD", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
					{ID: 3, Balance: 0, Currency: "SGD", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
				}, nil},
			},
			ExpectedError: exception.ErrInsufficientFunds,
		},
		{
			Title:                    "ReturnsError_UpdateBalances_InternalServerError",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_CreateWallet_InternalServerError",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
				CreateWallet:                    []interface{}{0, errors.New("internal server error")},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_InsertWalletCurrencyAmount_InternalServerError",
			GivenCreateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				CheckWalletExistsByWalletTypeID: []interface{}{false, nil},
				CheckWalletTypeExists:           []interface{}{true, nil},
				CreateWallet:                    []interface{}{1, nil},
				InsertWalletCurrencyAmount:      []interface{}{errors.New("internal server error")},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
			ExpectedError: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			// Create sqlmock database connection and mock object
			db, sqlmock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Wrap the sqlmock connection with sqlx
			sqlxDB := sqlx.NewDb(db, "sqlmock")

			// Initialize mocks for repository
			walletRepo := new(mocks.WalletRepository)
			balanceRepo := new(mocks.BalanceRepository)

			// Create instance of usecase
			walletUsecase := usecase.NewWalletUsecase(sqlxDB, walletRepo, balanceRepo)
			require.NotNil(t, walletUsecase)

			walletRepo.On("CheckWalletExistsByWalletTypeID", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.CheckWalletExistsByWalletTypeID...)
			walletRepo.On("CheckWalletTypeExists", mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.CheckWalletTypeExists...)
			balanceRepo.On("GetBalances", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.BalanceRepositoryReturnValues.GetBalances...)
			balanceRepo.On("UpdateBalances", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.BalanceRepositoryReturnValues.UpdateBalances...)
			walletRepo.On("CreateWallet", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.CreateWallet...)
			walletRepo.On("InsertWalletCurrencyAmount", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.InsertWalletCurrencyAmount...)

			// Setting up expectations
			sqlmock.ExpectBegin()

			// Handle different scenarios for transaction outcome
			if tc.ExpectedError != nil {
				sqlmock.ExpectRollback()
			} else {
				sqlmock.ExpectCommit()
			}

			// Call the function under test
			givenUserID := 1
			err = walletUsecase.CreateWallet(context.Background(), givenUserID, *tc.GivenCreateWalletRequest)

			// Validate the outcome based on expected errors
			if tc.ExpectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError.Error(), err.Error())
			}

			// Ensure all expectations were met
			require.NoError(t, sqlmock.ExpectationsWereMet())
		})
	}
}

func TestWalletUsecase_TopUpWallet(t *testing.T) {
	basicRequest := dto.UpdateWalletRequest{
		CurrencyAmount: []dto.CurrencyAmount{
			{Amount: 100, Currency: "USD"},
			{Amount: 200, Currency: "SGD"},
		},
	}

	testCases := []struct {
		Title                         string
		GivenUpdateWalletRequest      dto.UpdateWalletRequest
		WalletRepositoryReturnValues  mocks.WalletRepositoryReturnValues
		BalanceRepositoryReturnValues mocks.BalanceRepositoryReturnValues
		ExpectedError                 error
	}{
		{
			Title:                    "ReturnSuccessfully",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
				TopUpWalletBalances:                  []interface{}{nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
		},
		{
			Title:                    "ReturnsError_GetWalletByWalletID_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_GetBalances_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID: []interface{}{nil, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_GetWalletBalancesByUserIDAndWalletID_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{nil, errors.New("internal server error")},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{testdata.NewBalances(), nil},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_ErrBalanceNotFound_ForOneCurrency",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{[]domain.Balance{
					{ID: 1, Balance: 500.0, Currency: "EUR", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
				}, nil},
			},
			ExpectedError: exception.ErrBalanceNotFound,
		},
		{
			Title:                    "ReturnsError_InsufficientFundsForTopUpWallet",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{[]domain.Balance{
					{ID: 1, Balance: 1.0, Currency: "USD", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
					{ID: 3, Balance: 0, Currency: "SGD", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
				}, nil},
			},
			ExpectedError: exception.ErrInsufficientFunds,
		},
		{
			Title:                    "ReturnsError_UpdateBalances_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_TopUpWalletBalances_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
				TopUpWalletBalances:                  []interface{}{errors.New("internal server error")},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
			ExpectedError: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			// Create sqlmock database connection and mock object
			db, sqlmock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Wrap the sqlmock connection with sqlx
			sqlxDB := sqlx.NewDb(db, "sqlmock")

			// Initialize mocks for repository
			walletRepo := new(mocks.WalletRepository)
			balanceRepo := new(mocks.BalanceRepository)

			// Create instance of usecase
			walletUsecase := usecase.NewWalletUsecase(sqlxDB, walletRepo, balanceRepo)
			require.NotNil(t, walletUsecase)

			walletRepo.On("GetWalletByWalletID", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletByWalletID...)
			balanceRepo.On("GetBalances", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.BalanceRepositoryReturnValues.GetBalances...)
			balanceRepo.On("UpdateBalances", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.BalanceRepositoryReturnValues.UpdateBalances...)
			walletRepo.On("GetWalletBalancesByUserIDAndWalletID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletBalancesByUserIDAndWalletID...)
			walletRepo.On("TopUpWalletBalances", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.TopUpWalletBalances...)

			// Setting up expectations
			sqlmock.ExpectBegin()

			// Handle different scenarios for transaction outcome
			if tc.ExpectedError != nil {
				sqlmock.ExpectRollback()
			} else {
				sqlmock.ExpectCommit()
			}

			// Call the function under test
			givenUserID, givenWalletID := 1, 1
			err = walletUsecase.TopUpWallet(context.Background(), givenUserID, givenWalletID, tc.GivenUpdateWalletRequest)

			// Validate the outcome based on expected errors
			if tc.ExpectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError.Error(), err.Error())
			}

			// Ensure all expectations were met
			require.NoError(t, sqlmock.ExpectationsWereMet())
		})
	}
}

func TestWalletUsecase_CashOutWallet(t *testing.T) {
	basicRequest := dto.UpdateWalletRequest{
		CurrencyAmount: []dto.CurrencyAmount{
			{Amount: 50, Currency: "USD"},
			{Amount: 50, Currency: "SGD"},
		},
	}
	basicRequestCurrencyNotFoundInWallet := dto.UpdateWalletRequest{
		CurrencyAmount: []dto.CurrencyAmount{
			{Amount: 50, Currency: "USD"},
			{Amount: 50, Currency: "IDR"},
		},
	}
	basicRequestInsufficientFundsInWalletForWithdrawal := dto.UpdateWalletRequest{
		CurrencyAmount: []dto.CurrencyAmount{
			{Amount: 1000000, Currency: "USD"},
			{Amount: 2000000, Currency: "SGD"},
		},
	}

	testCases := []struct {
		Title                         string
		GivenUpdateWalletRequest      dto.UpdateWalletRequest
		WalletRepositoryReturnValues  mocks.WalletRepositoryReturnValues
		BalanceRepositoryReturnValues mocks.BalanceRepositoryReturnValues
		ExpectedError                 error
	}{
		{
			Title:                    "ReturnsSuccessfully",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
				CashOutWalletBalances:                []interface{}{nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
		},
		{
			Title:                    "ReturnsError_GetWalletByWalletID_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_GetBalances_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID: []interface{}{nil, nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{nil, errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_GetWalletBalancesByUserIDAndWalletID_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{nil, errors.New("internal server error")},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{testdata.NewBalances(), nil},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_ErrWalletBalanceNotFound_ForOneCurrency",
			GivenUpdateWalletRequest: basicRequestCurrencyNotFoundInWallet,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{[]domain.Balance{
					{ID: 1, Balance: 500.0, Currency: "EUR", UserID: 1, CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
				}, nil},
			},
			ExpectedError: exception.ErrWalletBalanceNotFound,
		},
		{
			Title:                    "ReturnsError_ErrInsufficientFundsForWithdrawal",
			GivenUpdateWalletRequest: basicRequestInsufficientFundsInWalletForWithdrawal,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances: []interface{}{testdata.NewBalances(), nil},
			},
			ExpectedError: exception.ErrInsufficientFundsForWithdrawal,
		},
		{
			Title:                    "ReturnsError_UpdateBalances_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{errors.New("internal server error")},
			},
			ExpectedError: errors.New("internal server error"),
		},
		{
			Title:                    "ReturnsError_CashOutWalletBalances_InternalServerError",
			GivenUpdateWalletRequest: basicRequest,
			WalletRepositoryReturnValues: mocks.WalletRepositoryReturnValues{
				GetWalletByWalletID:                  []interface{}{nil, nil},
				GetWalletBalancesByUserIDAndWalletID: []interface{}{testdata.NewWalletCurrencyAmount(), nil},
				CashOutWalletBalances:                []interface{}{errors.New("internal server error")},
			},
			BalanceRepositoryReturnValues: mocks.BalanceRepositoryReturnValues{
				GetBalances:    []interface{}{testdata.NewBalances(), nil},
				UpdateBalances: []interface{}{nil},
			},
			ExpectedError: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			// Create sqlmock database connection and mock object
			db, sqlmock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Wrap the sqlmock connection with sqlx
			sqlxDB := sqlx.NewDb(db, "sqlmock")

			// Initialize mocks for repository
			walletRepo := new(mocks.WalletRepository)
			balanceRepo := new(mocks.BalanceRepository)

			// Create instance of usecase
			walletUsecase := usecase.NewWalletUsecase(sqlxDB, walletRepo, balanceRepo)
			require.NotNil(t, walletUsecase)

			walletRepo.On("GetWalletByWalletID", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletByWalletID...)
			balanceRepo.On("GetBalances", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.BalanceRepositoryReturnValues.GetBalances...)
			balanceRepo.On("UpdateBalances", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.BalanceRepositoryReturnValues.UpdateBalances...)
			walletRepo.On("GetWalletBalancesByUserIDAndWalletID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.GetWalletBalancesByUserIDAndWalletID...)
			walletRepo.On("CashOutWalletBalances", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.WalletRepositoryReturnValues.CashOutWalletBalances...)

			// Setting up expectations
			sqlmock.ExpectBegin()

			// Handle different scenarios for transaction outcome
			if tc.ExpectedError != nil {
				sqlmock.ExpectRollback()
			} else {
				sqlmock.ExpectCommit()
			}

			// Call the function under test
			givenUserID, givenWalletID := 1, 1
			err = walletUsecase.CashOutWallet(context.Background(), givenUserID, givenWalletID, tc.GivenUpdateWalletRequest)

			// Validate the outcome based on expected errors
			if tc.ExpectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError.Error(), err.Error())
			}

			// Ensure all expectations were met
			require.NoError(t, sqlmock.ExpectationsWereMet())
		})
	}
}
