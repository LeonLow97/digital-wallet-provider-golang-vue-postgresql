package app

import (
	"net/http"

	handlers "github.com/LeonLow97/go-clean-architecture/delivery/http/handler"
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/gorilla/mux"
)

type Application struct {
	UserUsecase        domain.UserUsecase
	BalanceUsecase     domain.BalanceUsecase
	BeneficiaryUsecase domain.BeneficiaryUsecase
	WalletUsecase      domain.WalletUsecase
	TransactionUsecase domain.TransactionUsecase
}

func (app Application) CreateRouter() (*mux.Router, error) {
	// Create the base router with API version prefix
	apiRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api/v1").Subrouter()

	// standard health check endpoint
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	}).Methods(http.MethodGet).Name("HEALTH")

	// Handlers
	userHandler := handlers.NewUserHandler(app.UserUsecase)
	balanceHandler := handlers.NewBalanceHandler(app.BalanceUsecase)
	beneficiaryHandler := handlers.NewBeneficiaryHandler(app.BeneficiaryUsecase)
	walletHandler := handlers.NewWalletHandler(app.WalletUsecase)
	transactionHandler := handlers.NewTransactionHandler(app.TransactionUsecase)

	// authentication routes
	apiRouter.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	apiRouter.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	apiRouter.HandleFunc("/logout", userHandler.Logout).Methods(http.MethodPost)
	apiRouter.HandleFunc("/change-password", userHandler.ChangePassword).Methods(http.MethodPatch)
	apiRouter.HandleFunc("/configure-mfa", userHandler.ConfigureMFA).Methods(http.MethodPost)
	apiRouter.HandleFunc("/verify-mfa", userHandler.VerifyMFA).Methods(http.MethodPost)

	// password reset routes
	apiRouter.HandleFunc("/password-reset/send", userHandler.SendPasswordResetEmail).Methods(http.MethodPost)
	apiRouter.HandleFunc("/password-reset/reset", userHandler.PasswordReset).Methods(http.MethodPatch)

	// user routes
	apiRouter.HandleFunc("/users/profile", userHandler.UpdateUser).Methods(http.MethodPut)
	apiRouter.HandleFunc("/users/me", userHandler.GetUserDetail).Methods(http.MethodGet)

	// balance routes
	apiRouter.HandleFunc("/balances", balanceHandler.GetBalances).Methods(http.MethodGet)
	apiRouter.HandleFunc("/balances/{id:[0-9]+}", balanceHandler.GetBalance).Methods(http.MethodGet)
	apiRouter.HandleFunc("/balances/history/{id:[0-9]+}", balanceHandler.GetBalanceHistory).Methods(http.MethodGet)
	apiRouter.HandleFunc("/balances/currencies", balanceHandler.GetUserBalanceCurrencies).Methods(http.MethodGet)
	apiRouter.HandleFunc("/balances/deposit", balanceHandler.Deposit).Methods(http.MethodPost)
	apiRouter.HandleFunc("/balances/withdraw", balanceHandler.Withdraw).Methods(http.MethodPost)
	apiRouter.HandleFunc("/balances/currency-exchange", balanceHandler.CurrencyExchange).Methods(http.MethodPatch)
	apiRouter.HandleFunc("/balances/preview-exchange", balanceHandler.PreviewExchange).Methods(http.MethodPost)

	// beneficiary routes
	apiRouter.HandleFunc("/beneficiary", beneficiaryHandler.CreateBeneficiary).Methods(http.MethodPost)
	apiRouter.HandleFunc("/beneficiary", beneficiaryHandler.UpdateBeneficiary).Methods(http.MethodPut)
	apiRouter.HandleFunc("/beneficiary/{id:[0-9]+}", beneficiaryHandler.GetBeneficiary).Methods(http.MethodGet)
	apiRouter.HandleFunc("/beneficiary", beneficiaryHandler.GetBeneficiaries).Methods(http.MethodGet)

	// wallet routes
	apiRouter.HandleFunc("/wallet/{id:[0-9]+}", walletHandler.GetWallet).Methods(http.MethodGet)
	apiRouter.HandleFunc("/wallet/all", walletHandler.GetWallets).Methods(http.MethodGet)
	apiRouter.HandleFunc("/wallet/types", walletHandler.GetWalletTypes).Methods(http.MethodGet)
	apiRouter.HandleFunc("/wallet", walletHandler.CreateWallet).Methods(http.MethodPost)
	apiRouter.HandleFunc("/wallet/update/{id:[0-9]+}/{operation}", walletHandler.UpdateWallet).Methods(http.MethodPut)

	// transaction routes
	apiRouter.HandleFunc("/transaction", transactionHandler.CreateTransaction).Methods(http.MethodPost)
	apiRouter.HandleFunc("/transaction/all", transactionHandler.GetTransactions).Methods(http.MethodGet)

	return apiRouter, nil
}
