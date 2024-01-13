package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	AuthUseCase domain.UserUsecase
}

func NewAuthHandler(router *mux.Router, uc domain.UserUsecase) {
	handler := &AuthHandler{
		AuthUseCase: uc,
	}

	router.HandleFunc("/login", handler.LoginHandler)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	user, token, err := h.AuthUseCase.Login(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserNotFound):
		utils.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
	case errors.Is(err, exception.ErrInactiveUser):
		utils.ErrorJSON(w, apiErr.ErrInactiveUser, http.StatusUnauthorized)
	case errors.Is(err, exception.ErrInvalidCredentials):
		utils.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		resp := dto.LoginResponse{
			Email:    user.Email,
			Username: user.Username,
			Active:   user.Active,
			Admin:    user.Admin,
		}
		cookie := &http.Cookie{
			Name:     "mw-token",
			Value:    token.AccessToken,
			MaxAge:   3600,
			Path:     "/",
			Domain:   "mobilewallet.com",
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}
