package handlers

import (
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/delivery/http/middleware"
	"github.com/LeonLow97/go-clean-architecture/domain"
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
	balanceRouter.Use(middleware.AuthenticationMiddleware)

	balanceRouter.HandleFunc("/deposit", handler.Deposit).Methods(http.MethodPost)
}

func (h *BalanceHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, "deposit")
}
