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

type BalanceHandler struct {
	balanceUsecase domain.BalanceUsecase
}

func NewBalanceHandler(router *mux.Router, uc domain.BalanceUsecase) {
	handler := &BalanceHandler{
		balanceUsecase: uc,
	}

	balanceRouter := router.PathPrefix("/balances").Subrouter()

	balanceRouter.HandleFunc("", handler.GetBalances).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/{id:[0-9]+}", handler.GetBalance).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/history/{id:[0-9]+}", handler.GetBalanceHistory).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/currencies", handler.GetUserBalanceCurrencies).Methods(http.MethodGet)
	balanceRouter.HandleFunc("/deposit", handler.Deposit).Methods(http.MethodPost)
	balanceRouter.HandleFunc("/withdraw", handler.Withdraw).Methods(http.MethodPost)
	balanceRouter.HandleFunc("/currency-exchange", handler.CurrencyExchange).Methods(http.MethodPatch)
}

func (h *BalanceHandler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// retrieve balance id from url params
	balanceID, err := utils.ReadParamsInt(w, r, "id")
	if err != nil {
		return
	}

	resp, err := h.balanceUsecase.GetBalanceHistory(ctx, userID, balanceID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceHistoryNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceHistoryNotFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// retrieve balance id from url params
	balanceID, err := utils.ReadParamsInt(w, r, "id")
	if err != nil {
		return
	}

	resp, err := h.balanceUsecase.GetBalance(ctx, userID, balanceID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *BalanceHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	balances, err := h.balanceUsecase.GetBalances(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, balances)
}

func (h *BalanceHandler) GetUserBalanceCurrencies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	resp, err := h.balanceUsecase.GetUserBalanceCurrencies(ctx, userID)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *BalanceHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.DepositRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in deposit handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.DepositSanitize()

	if err := h.balanceUsecase.Deposit(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrDepositCurrencyNotAllowed):
			utils.ErrorJSON(w, apiErr.ErrDepositCurrencyNotAllowed, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.WithdrawRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in withdraw handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.WithdrawSanitize()

	if err := h.balanceUsecase.Withdraw(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrWithdrawCurrencyNotAllowed):
			utils.ErrorJSON(w, apiErr.ErrWithdrawCurrencyNotAllowed, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFunds):
			utils.ErrorJSON(w, apiErr.ErrInsufficientFundsForWithdrawal, http.StatusBadRequest)
		case errors.Is(err, exception.ErrBalanceNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *BalanceHandler) CurrencyExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CurrencyExchangeRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in withdraw handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.CurrencyExchangeSanitize()

	if err := h.balanceUsecase.CurrencyExchange(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrInsufficientFundsForCurrencyExchange):
			utils.ErrorJSON(w, apiErr.ErrInsufficientFundsForCurrencyExchange, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}
