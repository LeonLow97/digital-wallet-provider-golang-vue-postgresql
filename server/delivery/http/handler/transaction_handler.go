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

	var err error
	if req.SenderWalletID > 0 {
		// if sender wallet id is more than 0, user is trying to transfer funds with wallet
		err = h.transactionUsecase.CreateTransactionByWallet(ctx, req, userID)
	}

	switch {
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
		utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserAndWalletAssociationNotFound):
		utils.ErrorJSON(w, apiErr.ErrUserAndWalletAssociationNotFound, http.StatusBadRequest)
	default:
		utils.WriteNoContent(w, http.StatusCreated)
	}
}
