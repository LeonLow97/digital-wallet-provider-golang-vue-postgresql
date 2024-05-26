package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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
	walletRouter.HandleFunc("", handler.UpdateWallet).Methods(http.MethodPut)
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
	walletID, err := utils.ReadParamsInt(r, "id")
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	resp, err := h.walletUseCase.GetWallet(ctx, userID, walletID)
	switch {
	case errors.Is(err, exception.ErrNoWalletFound):
		utils.ErrorJSON(w, apiErr.ErrNoWalletFound, http.StatusNotFound)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
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
	switch {
	case errors.Is(err, exception.ErrNoWalletsFound):
		utils.ErrorJSON(w, apiErr.ErrNoWalletsFound, http.StatusNotFound)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in create wallet handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
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

func (h *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.UpdateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in create wallet handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create wallet handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.Type = strings.ToLower(req.Type)
	req.UserID = userID
	req.UpdateWalletSanitize()

	resp, err := h.walletUseCase.UpdateWallet(ctx, req)
	switch {
	case errors.Is(err, exception.ErrNoWalletFound):
		utils.ErrorJSON(w, apiErr.ErrNoWalletFound, http.StatusNotFound)
	case errors.Is(err, exception.ErrBalanceNotFound):
		utils.ErrorJSON(w, apiErr.ErrBalanceNotFound, http.StatusNotFound)
	case errors.Is(err, exception.ErrInsufficientFundsForWithdrawal):
		utils.ErrorJSON(w, apiErr.ErrInsufficientFundsForWithdrawalFromWallet, http.StatusBadRequest)
	case errors.Is(err, exception.ErrInsufficientFundsForDeposit):
		utils.ErrorJSON(w, apiErr.ErrInsufficientFundsForDepositToWallet, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}
