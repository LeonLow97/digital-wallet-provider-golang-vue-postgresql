package handlers

import (
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/delivery/http/middleware"
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	transactionUsecase domain.TransactionUsecase
}

func NewTransactionHandler(router *mux.Router, uc domain.TransactionUsecase) {
	handler := &TransactionHandler{
		transactionUsecase: uc,
	}

	transactionRouter := router.PathPrefix("/transaction").Subrouter()
	transactionRouter.Use(middleware.AuthenticationMiddleware)

	transactionRouter.HandleFunc("/deposit", handler.Deposit).Methods(http.MethodPost)
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, "test")
}
