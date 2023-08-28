package main

import (
	"log"
	"net/http"

	"github.com/LeonLow97/internal/beneficiaries"
	"github.com/LeonLow97/internal/handlers"
	"github.com/LeonLow97/internal/transactions"
	"github.com/LeonLow97/internal/users"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func routes(db *sqlx.DB) *mux.Router {
	router := mux.NewRouter()

	// middleware
	router.Use(setAccessControlHeader)

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
	beneficiaryHandler, err := handlers.NewBeneficiaryHandler(beneficiaryService)
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
	transactionHandler, err := handlers.NewTransactionHandler(transactionService)
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/user", userHandler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/beneficiaries", beneficiaryHandler.GetBeneficiaries).Methods(http.MethodGet)
	router.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods(http.MethodGet)
	router.HandleFunc("/transaction", transactionHandler.CreateTransaction).Methods(http.MethodPost)

	router.Methods(http.MethodOptions).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	return router
}

func setAccessControlHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}
