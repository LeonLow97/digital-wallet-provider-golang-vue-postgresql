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

type WalletHandler struct {
	walletUseCase domain.WalletUsecase
}

func NewWalletHandler(router *mux.Router, uc domain.WalletUsecase) {
	handler := &WalletHandler{
		walletUseCase: uc,
	}

	walletRouter := router.PathPrefix("/wallet").Subrouter()

	walletRouter.HandleFunc("/{id:[0-9]+}", handler.GetWallet).Methods(http.MethodGet)
	walletRouter.HandleFunc("/all", handler.GetWallets).Methods(http.MethodGet)
	walletRouter.HandleFunc("/types", handler.GetWalletTypes).Methods(http.MethodGet)
	walletRouter.HandleFunc("", handler.CreateWallet).Methods(http.MethodPost)
	walletRouter.HandleFunc("/topup/{id:[0-9]+}", handler.TopUpWallet).Methods(http.MethodPut)
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// retrieve wallet id from url params
	walletID, err := utils.ReadParamsInt(w, r, "id")
	if err != nil {
		return
	}

	resp, err := h.walletUseCase.GetWallet(ctx, userID, walletID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletFound):
			utils.ErrorJSON(w, apiErr.ErrNoWalletFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	resp, err := h.walletUseCase.GetWallets(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletsFound):
			utils.ErrorJSON(w, apiErr.ErrNoWalletsFound, http.StatusNotFound)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) GetWalletTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.walletUseCase.GetWalletTypes(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateWalletRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create wallet handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.CreateWalletSanitize()

	if err = h.walletUseCase.CreateWallet(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrBalanceNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrWalletTypeInvalid):
			utils.ErrorJSON(w, apiErr.ErrWalletTypeInvalid, http.StatusBadRequest)
		case errors.Is(err, exception.ErrWalletAlreadyExists):
			utils.ErrorJSON(w, apiErr.ErrWalletAlreadyExists, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFunds):
			utils.ErrorJSON(w, apiErr.ErrInsufficientFundsInAccount, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}
	utils.WriteNoContent(w, http.StatusCreated)
}

func (h *WalletHandler) TopUpWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	walletID, err := utils.ReadParamsInt(w, r, "id")
	if err != nil {
		return
	}

	var req dto.UpdateWalletRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		return
	}

	if err := h.walletUseCase.TopUpWallet(ctx, userID, walletID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletFound):
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrWalletTypeInvalid):
			utils.ErrorJSON(w, apiErr.ErrWalletTypeInvalid, http.StatusBadRequest)
		case errors.Is(err, exception.ErrBalanceNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFunds):
			utils.ErrorJSON(w, apiErr.ErrInsufficientFundsInAccount, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *WalletHandler) CashOutWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	walletID, err := utils.ReadParamsInt(w, r, "id")
	if err != nil {
		return
	}

	var req dto.UpdateWalletRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		return
	}

	if err := h.walletUseCase.CashOutWallet(ctx, userID, walletID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrNoWalletFound):
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrWalletTypeInvalid):
			utils.ErrorJSON(w, apiErr.ErrWalletTypeInvalid, http.StatusBadRequest)
		case errors.Is(err, exception.ErrWalletBalanceNotFound):
			utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInsufficientFundsForWithdrawal):
			utils.ErrorJSON(w, apiErr.ErrInsufficientFundsForWithdrawal, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}
