package app_test

import (
	"strings"
	"testing"

	"github.com/LeonLow97/go-clean-architecture/delivery/http/app"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	mocks "github.com/LeonLow97/go-clean-architecture/mocks/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateRouter(t *testing.T) {
	mockConfig := new(infrastructure.Config)
	mockRedisClient := *new(infrastructure.RedisClient)
	mockUserUsecase := new(mocks.UserUsecase)
	mockBalanceUsecase := new(mocks.BalanceUsecase)
	mockBeneficiaryUsecase := new(mocks.BeneficiaryUsecase)
	mockWalletUsecase := new(mocks.WalletUsecase)
	mockTransactionUsecase := new(mocks.TransactionUsecase)

	app := app.Application{
		Cfg:                mockConfig,
		RedisClient:        mockRedisClient,
		UserUsecase:        mockUserUsecase,
		BalanceUsecase:     mockBalanceUsecase,
		BeneficiaryUsecase: mockBeneficiaryUsecase,
		WalletUsecase:      mockWalletUsecase,
		TransactionUsecase: mockTransactionUsecase,
	}

	router, err := app.CreateRouter()
	require.NoError(t, err)

	expectedRoutes := []struct {
		Path   string
		Method string
	}{
		{"/api/v1/health", "GET"},
		{"/api/v1/login", "POST"},
		{"/api/v1/signup", "POST"},
		{"/api/v1/logout", "POST"},
		{"/api/v1/change-password", "PATCH"},
		{"/api/v1/configure-mfa", "POST"},
		{"/api/v1/verify-mfa", "POST"},
		{"/api/v1/password-reset/send", "POST"},
		{"/api/v1/password-reset/reset", "PATCH"},
		{"/api/v1/users/profile", "PUT"},
		{"/api/v1/users/me", "GET"},
		{"/api/v1/balances", "GET"},
		{"/api/v1/balances/{id:[0-9]+}", "GET"},
		{"/api/v1/balances/history/{id:[0-9]+}", "GET"},
		{"/api/v1/balances/currencies", "GET"},
		{"/api/v1/balances/deposit", "POST"},
		{"/api/v1/balances/withdraw", "POST"},
		{"/api/v1/balances/currency-exchange", "PATCH"},
		{"/api/v1/balances/preview-exchange", "POST"},
		{"/api/v1/beneficiary", "POST"},
		{"/api/v1/beneficiary", "PUT"},
		{"/api/v1/beneficiary/{id:[0-9]+}", "GET"},
		{"/api/v1/beneficiary", "GET"},
		{"/api/v1/wallet/{id:[0-9]+}", "GET"},
		{"/api/v1/wallet/all", "GET"},
		{"/api/v1/wallet/types", "GET"},
		{"/api/v1/wallet", "POST"},
		{"/api/v1/wallet/update/{id:[0-9]+}/{operation}", "PUT"},
		{"/api/v1/transaction", "POST"},
		{"/api/v1/transaction/all", "GET"},
	}

	// Check if each expected route exists in the router with the correct method
	for _, route := range expectedRoutes {
		require.True(t, routeExists(router, route.Path, route.Method), "Route not found: %s %s", route.Method, route.Path)
	}
}

func routeExists(router *mux.Router, testPath, testMethod string) bool {
	found := false

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		// Build the full path from ancestors and the current route
		fullPath := ""
		for _, ancestor := range ancestors {
			ancestorPath, err := ancestor.GetPathTemplate()
			if err == nil {
				fullPath += ancestorPath
			}
		}
		fullPath += path

		methods, err := route.GetMethods()
		if err != nil && err != mux.ErrMethodMismatch {
			return err
		}

		// Check if the full path and method match
		if strings.EqualFold(fullPath, testPath) {
			for _, method := range methods {
				if strings.EqualFold(method, testMethod) {
					found = true
					break
				}
			}
		}

		return nil
	})

	return found
}
