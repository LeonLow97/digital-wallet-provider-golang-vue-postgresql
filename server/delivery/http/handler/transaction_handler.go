package handlers

import (
	"fmt"
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/domain"
	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
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

	transactionRouter.HandleFunc("", handler.CreateTransaction).Methods(http.MethodPost)
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	fmt.Println("UserID", userID)
}
