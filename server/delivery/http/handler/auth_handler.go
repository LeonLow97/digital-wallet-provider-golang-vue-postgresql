package handlers

import (
	"context"
	"encoding/json"
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

type AuthHandler struct {
	authUseCase domain.UserUsecase
}

func NewAuthHandler(router *mux.Router, uc domain.UserUsecase) {
	handler := &AuthHandler{
		authUseCase: uc,
	}

	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/signup", handler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/logout", handler.Logout).Methods(http.MethodPost)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in login handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in login handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.LoginSanitize()

	resp, token, err := h.authUseCase.Login(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserNotFound) || errors.Is(err, exception.ErrInvalidCredentials):
		utils.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
	case errors.Is(err, exception.ErrInactiveUser):
		utils.ErrorJSON(w, apiErr.ErrInactiveUser, http.StatusUnauthorized)
	case err != nil:
		log.Println("error in login handler", err)
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		cookie := &http.Cookie{
			Name:     "mw-token",
			Value:    token.AccessToken,
			MaxAge:   3600,
			Path:     "/",
			Domain:   "localhost", // TODO: replace with config domain name
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req dto.SignUpRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		log.Println("error decoding req body in sign up handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	req.SignUpSanitize()

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in sign up handler", errMessage)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	err = h.authUseCase.SignUp(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserFound):
		utils.ErrorJSON(w, apiErr.ErrUserFound, http.StatusBadRequest)
	case errors.Is(err, exception.ErrInvalidPassword):
		utils.ErrorJSON(w, apiErr.ErrInvalidPassword, http.StatusBadRequest)
	case err != nil:
		log.Println("error in sign up handler", err)
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteNoContent(w, http.StatusOK)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "mw-token",
		Value:    "",
		MaxAge:   -1,
		Path:     "",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	utils.WriteNoContent(w, http.StatusOK)
}
