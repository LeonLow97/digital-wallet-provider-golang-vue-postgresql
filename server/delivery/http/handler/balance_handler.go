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
	"github.com/LeonLow97/go-clean-architecture/utils/contextstore"
	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
)

type BalanceHandler struct {
	balanceUsecase domain.BalanceUsecase
}

func NewBalanceHandler(uc domain.BalanceUsecase) *BalanceHandler {
	handler := &BalanceHandler{
		balanceUsecase: uc,
	}

	return handler
}

func (h *BalanceHandler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// retrieve balance id from url params
	balanceID, err := jsonutil.ReadURLParamsInt(w, r, "id")
	if err != nil {
		return
	}

	resp, err := h.balanceUsecase.GetBalanceHistory(ctx, userID, balanceID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceHistoryNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceHistoryNotFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// retrieve balance id from url params
	balanceID, err := jsonutil.ReadURLParamsInt(w, r, "id")
	if err != nil {
		return
	}

	resp, err := h.balanceUsecase.GetBalance(ctx, userID, balanceID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *BalanceHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	balances, err := h.balanceUsecase.GetBalances(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrBalancesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, balances)
}

func (h *BalanceHandler) GetUserBalanceCurrencies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	resp, err := h.balanceUsecase.GetUserBalanceCurrencies(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserCurrenciesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUserCurrenciesNotFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *BalanceHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.DepositRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in deposit handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.DepositSanitize()

	if err := h.balanceUsecase.Deposit(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrDepositCurrencyNotAllowed):
			jsonutil.ErrorJSON(w, apiErr.ErrDepositCurrencyNotAllowed, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.WithdrawRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in withdraw handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.WithdrawSanitize()

	if err := h.balanceUsecase.Withdraw(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrWithdrawCurrencyNotAllowed):
			jsonutil.ErrorJSON(w, apiErr.ErrWithdrawCurrencyNotAllowed, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFunds):
			jsonutil.ErrorJSON(w, apiErr.ErrInsufficientFundsForWithdrawal, http.StatusBadRequest)
		case errors.Is(err, exception.ErrBalanceNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *BalanceHandler) CurrencyExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CurrencyExchangeRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in withdraw handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.CurrencyExchangeSanitize()

	if err := h.balanceUsecase.CurrencyExchange(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrBalancesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrInsufficientFundsForCurrencyExchange):
			jsonutil.ErrorJSON(w, apiErr.ErrInsufficientFundsForCurrencyExchange, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

// PreviewExchange allows users to preview the exchange rate, no data manipulation
func (h *BalanceHandler) PreviewExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.PreviewExchangeRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in withdraw handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.PreviewExchangeSanitize()

	resp := h.balanceUsecase.PreviewExchange(ctx, req)
	jsonutil.WriteJSON(w, http.StatusOK, resp)
}
