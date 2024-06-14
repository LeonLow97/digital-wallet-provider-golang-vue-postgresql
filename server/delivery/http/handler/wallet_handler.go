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

type WalletHandler struct {
	walletUseCase domain.WalletUsecase
}

func NewWalletHandler(uc domain.WalletUsecase) *WalletHandler {
	handler := &WalletHandler{
		walletUseCase: uc,
	}

	return handler
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// retrieve wallet id from url params
	walletID, err := jsonutil.ReadURLParamsInt(w, r, "id")
	if err != nil {
		return
	}

	resp, err := h.walletUseCase.GetWallet(ctx, userID, walletID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletFound):
			jsonutil.ErrorJSON(w, apiErr.ErrNoWalletFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrWalletBalancesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrNoWalletBalancesFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	resp, err := h.walletUseCase.GetWallets(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletsFound):
			jsonutil.ErrorJSON(w, apiErr.ErrNoWalletsFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrWalletBalancesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrNoWalletBalancesFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) GetWalletTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.walletUseCase.GetWalletTypes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrWalletTypesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrWalletTypesNotFound, http.StatusNotFound)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateWalletRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create wallet handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.CreateWalletSanitize()

	if err = h.walletUseCase.CreateWallet(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrWalletTypeInvalid):
			jsonutil.ErrorJSON(w, apiErr.ErrWalletTypeInvalid, http.StatusBadRequest)
		case errors.Is(err, exception.ErrWalletAlreadyExists):
			jsonutil.ErrorJSON(w, apiErr.ErrWalletAlreadyExists, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFunds):
			jsonutil.ErrorJSON(w, apiErr.ErrInsufficientFundsInAccount, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}
	jsonutil.WriteNoContent(w, http.StatusCreated)
}

func (h *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	walletID, err := jsonutil.ReadURLParamsInt(w, r, "id")
	if err != nil {
		return
	}

	operation, err := jsonutil.ReadURLParamsString(w, r, "operation")
	if err != nil {
		return
	}

	var req dto.UpdateWalletRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in update wallet handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	// to determine topup or withdrawal
	switch operation {
	case "topup":
		err = h.walletUseCase.TopUpWallet(ctx, userID, walletID, req)
	case "withdraw":
		err = h.walletUseCase.CashOutWallet(ctx, userID, walletID, req)
	default:
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrWalletBalancesNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrNoWalletBalancesFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrBalanceNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFunds):
			jsonutil.ErrorJSON(w, apiErr.ErrInsufficientFundsInAccount, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFundsForWithdrawal):
			jsonutil.ErrorJSON(w, apiErr.ErrInsufficientFundsForWithdrawal, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}
