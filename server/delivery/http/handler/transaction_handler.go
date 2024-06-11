package handlers

import (
	"errors"
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/utils/context"
	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
)

type TransactionHandler struct {
	transactionUsecase domain.TransactionUsecase
}

func NewTransactionHandler(uc domain.TransactionUsecase) *TransactionHandler {
	handler := &TransactionHandler{
		transactionUsecase: uc,
	}

	return handler
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := context.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateTransactionRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.Sanitize()

	if err = h.transactionUsecase.CreateTransaction(ctx, req, userID); err != nil {
		switch {
		case errors.Is(err, exception.ErrBeneficiaryIsInactive) ||
			errors.Is(err, exception.ErrBeneficiaryMFANotConfigured):
			jsonutil.ErrorJSON(w, apiErr.ErrBeneficiaryAccountNotRegistered, http.StatusBadRequest)
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			jsonutil.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrSenderWalletInvalid):
			jsonutil.ErrorJSON(w, apiErr.ErrSenderWalletInvalid, http.StatusForbidden)
		case errors.Is(err, exception.ErrInsufficientFundsInWallet):
			jsonutil.ErrorJSON(w, apiErr.ErrInsufficientFundsInWallet, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusCreated)
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := context.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	transactions, err := h.transactionUsecase.GetTransactions(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoTransactionsFound):
			jsonutil.ErrorJSON(w, apiErr.ErrNoTransactionsFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, transactions)
}
