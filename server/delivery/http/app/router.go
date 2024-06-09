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
	baseRouter := mux.NewRouter().StrictSlash(true)           // trim trailing slashes in endpoint
	baseRouter = baseRouter.PathPrefix("/api/v1").Subrouter() // api versioning v1

	// standard health check endpoint
	baseRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	}).Methods(http.MethodGet).Name("HEALTH")

	router := baseRouter.PathPrefix("/").Subrouter()

	// Handlers
	userHandler := handlers.NewUserHandler(app.UserUsecase)
	balanceHandler := handlers.NewBalanceHandler(app.BalanceUsecase)
	beneficiaryHandler := handlers.NewBeneficiaryHandler(app.BeneficiaryUsecase)
	walletHandler := handlers.NewWalletHandler(app.WalletUsecase)
	transactionHandler := handlers.NewTransactionHandler(app.TransactionUsecase)

	// authentication routes
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/logout", userHandler.Logout).Methods(http.MethodPost)
	router.HandleFunc("/change-password", userHandler.ChangePassword).Methods(http.MethodPatch)
	router.HandleFunc("/configure-mfa", userHandler.ConfigureMFA).Methods(http.MethodPost)
	router.HandleFunc("/verify-mfa", userHandler.VerifyMFA).Methods(http.MethodPost)

	// password reset routes
	router.HandleFunc("/password-reset/send", userHandler.SendPasswordResetEmail).Methods(http.MethodPost)
	router.HandleFunc("/password-reset/reset", userHandler.PasswordReset).Methods(http.MethodPatch)

	// user routes
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/profile", userHandler.UpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("/me", userHandler.GetUserDetail).Methods(http.MethodGet)

	// balance routes
	balanceRouter := router.PathPrefix("/balances").Subrouter()
	balanceRouter.HandleFunc("", balanceHandler.GetBalances).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/{id:[0-9]+}", balanceHandler.GetBalance).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/history/{id:[0-9]+}", balanceHandler.GetBalanceHistory).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/currencies", balanceHandler.GetUserBalanceCurrencies).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/deposit", balanceHandler.Deposit).Methods(http.MethodPost)
	balanceRouter.HandleFunc("/withdraw", balanceHandler.Withdraw).Methods(http.MethodPost)
	balanceRouter.HandleFunc("/currency-exchange", balanceHandler.CurrencyExchange).Methods(http.MethodPatch)
	balanceRouter.HandleFunc("/preview-exchange", balanceHandler.PreviewExchange).Methods(http.MethodPost)

	// beneficiary routes
	beneficiaryRouter := router.PathPrefix("/beneficiary").Subrouter()
	beneficiaryRouter.HandleFunc("", beneficiaryHandler.CreateBeneficiary).Methods(http.MethodPost)
	beneficiaryRouter.HandleFunc("", beneficiaryHandler.UpdateBeneficiary).Methods(http.MethodPut)
	beneficiaryRouter.HandleFunc("/{id:[0-9]+}", beneficiaryHandler.GetBeneficiary).Methods(http.MethodGet)
	beneficiaryRouter.HandleFunc("", beneficiaryHandler.GetBeneficiaries).Methods(http.MethodGet)

	// wallet routes
	walletRouter := router.PathPrefix("/wallet").Subrouter()
	walletRouter.HandleFunc("/{id:[0-9]+}", walletHandler.GetWallet).Methods(http.MethodGet)
	walletRouter.HandleFunc("/all", walletHandler.GetWallets).Methods(http.MethodGet)
	walletRouter.HandleFunc("/types", walletHandler.GetWalletTypes).Methods(http.MethodGet)
	walletRouter.HandleFunc("", walletHandler.CreateWallet).Methods(http.MethodPost)
	walletRouter.HandleFunc("/topup/{id:[0-9]+}", walletHandler.TopUpWallet).Methods(http.MethodPut)
	walletRouter.HandleFunc("/cashout/{id:[0-9]+}", walletHandler.CashOutWallet).Methods(http.MethodPut)

	// transaction routes
	transactionRouter := router.PathPrefix("/transaction").Subrouter()
	transactionRouter.HandleFunc("", transactionHandler.CreateTransaction).Methods(http.MethodPost)
	transactionRouter.HandleFunc("/all", transactionHandler.GetTransactions).Methods(http.MethodGet)

	return baseRouter, nil
}
