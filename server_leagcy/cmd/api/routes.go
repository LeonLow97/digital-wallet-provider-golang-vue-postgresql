package main

import (
	"net/http"

	"github.com/LeonLow97/internal/auth"
	"github.com/LeonLow97/internal/beneficiaries"
	"github.com/LeonLow97/internal/transactions"
	"github.com/LeonLow97/internal/users"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func routes(db *sqlx.DB) *mux.Router {
	router := mux.NewRouter()

	// middleware
	router.Use(setAccessControlHeader)

	authRepo := auth.NewRepo(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	userRepo := users.NewRepo(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewUserHandler(userService)

	beneficiaryRepo := beneficiaries.NewRepo(db)
	beneficiaryService := beneficiaries.NewService(beneficiaryRepo)
	beneficiaryHandler := beneficiaries.NewBeneficiaryHandler(beneficiaryService)

	transactionRepo := transactions.NewRepo(db)
	transactionService := transactions.NewService(transactionRepo)
	transactionHandler := transactions.NewTransactionHandler(transactionService)

	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/user", userHandler.GetUser).Methods(http.MethodGet)

	// sub-router for /user routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(authTokenMiddleware)
	userRouter.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods(http.MethodGet)
	userRouter.HandleFunc("/transaction", transactionHandler.CreateTransaction).Methods(http.MethodPost)

	router.HandleFunc("/beneficiaries", beneficiaryHandler.GetBeneficiaries).Methods(http.MethodGet)

	router.Methods(http.MethodOptions).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	return router
}
