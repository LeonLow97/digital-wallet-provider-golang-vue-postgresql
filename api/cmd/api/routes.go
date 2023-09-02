package main

import (
	"log"
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

	authRepo, err := auth.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	authService, err := auth.NewService(authRepo)
	if err != nil {
		log.Fatal(err)
	}
	authHandler, err := auth.NewAuthHandler(authService)
	if err != nil {
		log.Fatal(err)
	}

	userRepo, err := users.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	userService, err := users.NewService(userRepo)
	if err != nil {
		log.Fatal(err)
	}
	userHandler, err := users.NewUserHandler(userService)
	if err != nil {
		log.Fatal(err)
	}

	beneficiaryRepo, err := beneficiaries.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	beneficiaryService, err := beneficiaries.NewService(beneficiaryRepo)
	if err != nil {
		log.Fatal(err)
	}
	beneficiaryHandler, err := beneficiaries.NewBeneficiaryHandler(beneficiaryService)
	if err != nil {
		log.Fatal(err)
	}

	transactionRepo, err := transactions.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	transactionService, err := transactions.NewService(transactionRepo)
	if err != nil {
		log.Fatal(err)
	}
	transactionHandler, err := transactions.NewTransactionHandler(transactionService)
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/user", userHandler.GetUser).Methods(http.MethodGet)

	// sub-router for /user routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(authTokenMiddleware)
	userRouter.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods(http.MethodPost)

	router.HandleFunc("/beneficiaries", beneficiaryHandler.GetBeneficiaries).Methods(http.MethodGet)
	router.HandleFunc("/transaction", transactionHandler.CreateTransaction).Methods(http.MethodPost)

	router.Methods(http.MethodOptions).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	return router
}
