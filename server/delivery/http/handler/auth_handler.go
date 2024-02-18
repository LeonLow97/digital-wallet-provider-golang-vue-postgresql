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
	redisClient infrastructure.RedisClient
}

func NewAuthHandler(router *mux.Router, uc domain.UserUsecase, redisClient infrastructure.RedisClient) {
	handler := &AuthHandler{
		authUseCase: uc,
		redisClient: redisClient,
	}

	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/signup", handler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/logout", handler.Logout).Methods(http.MethodPost)
	// TODO: reset password
	// TODO: update user details
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
		// TODO: utilise refresh token or remove it from auth use case
		utils.IssueCookie(w, token.AccessToken)

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
	// retrieve sessionID from context
	sessionID, ok := r.Context().Value(utils.SessionIDKey).(string)
	if !ok {
		log.Println("Unable to retrieve session id from context")
	}

	// remove sessionID from Redis
	if err := h.authUseCase.RemoveSessionFromRedis(r.Context(), sessionID); err != nil {
		// failed to remove sessionID but don't block the logout
		log.Println("failed to remove sessionID from Redis", err)
	}

	// override cookie in browser
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
