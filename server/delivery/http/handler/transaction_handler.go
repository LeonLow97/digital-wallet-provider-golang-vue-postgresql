package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
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
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateTransactionRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.Sanitize()

	if err = h.transactionUsecase.CreateTransaction(ctx, req, userID); err != nil {
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
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusCreated)
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	transactions, err := h.transactionUsecase.GetTransactions(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoTransactionsFound):
			utils.ErrorJSON(w, apiErr.ErrNoTransactionsFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, transactions)
}
