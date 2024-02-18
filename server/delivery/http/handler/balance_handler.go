package handlers

import (
	"context"
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

type BalanceHandler struct {
	balanceUsecase domain.BalanceUsecase
}

func NewBalanceHandler(router *mux.Router, uc domain.BalanceUsecase) {
	handler := &BalanceHandler{
		balanceUsecase: uc,
	}

	balanceRouter := router.PathPrefix("/balance").Subrouter()

	balanceRouter.HandleFunc("/deposit", handler.Deposit).Methods(http.MethodPost)
	balanceRouter.HandleFunc("/withdraw", handler.Withdraw).Methods(http.MethodPatch)
}

func (h *BalanceHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// retrieve user id from context
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.DepositRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		log.Println("error decoding req body in deposit handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.DepositSanitize()

	resp, err := h.balanceUsecase.Deposit(ctx, req)
	switch {
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// retrieve user id from context
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.WithdrawRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		log.Println("error decoding req body in withdraw handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.WithdrawSanitize()

	resp, err := h.balanceUsecase.Withdraw(ctx, req)
	switch {
	case errors.Is(err, exception.ErrBalanceNotFound):
		utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
	case errors.Is(err, exception.ErrInsufficientFunds):
		utils.ErrorJSON(w, apiErr.ErrInsufficientFundsForWithdrawal, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}
