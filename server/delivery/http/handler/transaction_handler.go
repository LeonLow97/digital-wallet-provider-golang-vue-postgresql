package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
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
	transactionRouter.HandleFunc("/all", handler.GetTransactions).Methods(http.MethodGet)
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in create transaction handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	req.Sanitize()

	err := h.transactionUsecase.CreateTransaction(ctx, req, userID)
	log.Println("err herer", err)

	switch {
	case errors.Is(err, exception.ErrAmountMustBePositive):
		utils.ErrorJSON(w, apiErr.ErrAmountMustBePositive, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserIsInactive):
		utils.ErrorJSON(w, apiErr.ErrInactiveUser, http.StatusBadRequest)
	case errors.Is(err, exception.ErrBeneficiaryIsInactive):
		utils.ErrorJSON(w, apiErr.ErrBeneficiaryIsInactive, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
		utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
	case errors.Is(err, exception.ErrInsufficientFundsInWallet):
		utils.ErrorJSON(w, apiErr.ErrInsufficientFundsInWallet, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
		utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserAndWalletAssociationNotFound):
		utils.ErrorJSON(w, apiErr.ErrUserAndWalletAssociationNotFound, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteNoContent(w, http.StatusCreated)
	}
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	transactions, err := h.transactionUsecase.GetTransactions(ctx, userID)
	switch {
	case errors.Is(err, exception.ErrNoTransactionsFound):
		utils.ErrorJSON(w, apiErr.ErrNoTransactionsFound, http.StatusNotFound)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, transactions)
	}
}
