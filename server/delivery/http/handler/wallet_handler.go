package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/LeonLow97/go-clean-architecture/delivery/http/middleware"
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
	walletRouter.Use(middleware.AuthenticationMiddleware)

	walletRouter.HandleFunc("", handler.CreateWallet).Methods(http.MethodPost)
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// retrieve user id from context
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in create wallet handler")
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
	req.CreateWalletSanitize()

	err = h.walletUseCase.CreateWallet(ctx, req)
	switch {
	case errors.Is(err, exception.ErrWalletTypeInvalid):
		utils.ErrorJSON(w, apiErr.ErrWalletTypeInvalid, http.StatusBadRequest)
	case errors.Is(err, exception.ErrWalletAlreadyExists):
		utils.ErrorJSON(w, apiErr.ErrWalletAlreadyExists, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteNoContent(w, http.StatusCreated)
	}
}
