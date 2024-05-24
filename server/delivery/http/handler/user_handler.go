package handlers

import (
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

type UserHandler struct {
	userUsecase domain.UserUsecase
	redisClient infrastructure.RedisClient
}

func NewUserHandler(router *mux.Router, uc domain.UserUsecase, redisClient infrastructure.RedisClient) {
	handler := &UserHandler{
		userUsecase: uc,
		redisClient: redisClient,
	}

	// authentication routes
	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/signup", handler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/logout", handler.Logout).Methods(http.MethodPost)
	router.HandleFunc("/change-password", handler.ChangePassword).Methods(http.MethodPatch)
	router.HandleFunc("/configure-mfa", handler.ConfigureMFA).Methods(http.MethodPost)
	router.HandleFunc("/verify-mfa", handler.VerifyMFA).Methods(http.MethodPost)

	// password reset
	router.HandleFunc("/password-reset/send", handler.SendPasswordResetEmail).Methods(http.MethodPost)
	router.HandleFunc("/password-reset/reset", handler.PasswordReset).Methods(http.MethodPatch)

	// user routes
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/profile", handler.UpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("/me", handler.GetUserDetail).Methods(http.MethodGet)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	resp, token, err := h.userUsecase.Login(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserNotFound) || errors.Is(err, exception.ErrInvalidCredentials):
		utils.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
	case errors.Is(err, exception.ErrInactiveUser):
		utils.ErrorJSON(w, apiErr.ErrInactiveUser, http.StatusUnauthorized)
	case err != nil:
		log.Println("error in login handler", err)
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		// TODO: utilise refresh token or remove it from user use case
		w.Header().Set("X-CSRF-Token", token.CSRFToken)
		utils.IssueCookie(w, token.AccessToken)

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	err = h.userUsecase.SignUp(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserFound):
		utils.ErrorJSON(w, apiErr.ErrUserFound, http.StatusBadRequest)
	case errors.Is(err, exception.ErrInvalidPassword):
		utils.ErrorJSON(w, apiErr.ErrInvalidPassword, http.StatusBadRequest)
	case err != nil:
		log.Println("error in sign up handler", err)
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteNoContent(w, http.StatusNoContent)
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID, err := utils.SessionIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
	}

	// remove sessionID from Redis
	h.userUsecase.RemoveSessionFromRedis(r.Context(), sessionID)

	// override cookie in browser
	cookie := &http.Cookie{
		Name:     utils.JWT_COOKIE,
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

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in change password handler", err)
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.ChangePasswordSanitize()

	if err := h.userUsecase.ChangePassword(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			utils.ErrorJSON(w, apiErr.ErrUserNotFound, http.StatusNotFound)
			return
		case errors.Is(err, exception.ErrInvalidCredentials):
			utils.ErrorJSON(w, apiErr.ErrCurrentPasswordIncorrect, http.StatusBadRequest)
			return
		case errors.Is(err, exception.ErrSamePassword):
			utils.ErrorJSON(w, apiErr.ErrSamePassword, http.StatusBadRequest)
			return
		default:
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) ConfigureMFA(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.ConfigureMFARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body with error:", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.ConfigureMFASanitize()

	if err := h.userUsecase.ConfigureMFA(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidMFACode):
			utils.ErrorJSON(w, apiErr.ErrInvalidMFACode, http.StatusUnauthorized)
			return
		case errors.Is(err, exception.ErrUserNotFound):
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		case errors.Is(err, exception.ErrTOTPSecretExists):
			utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
			return
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.VerifyMFARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body with error:", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.VerifyMFASanitize()

	if err := h.userUsecase.VerifyMFA(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		case errors.Is(err, exception.ErrInvalidMFACode):
			utils.ErrorJSON(w, apiErr.ErrInvalidMFACode, http.StatusUnauthorized)
			return
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in password reset handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.PasswordResetSanitize()

	if err := h.userUsecase.PasswordReset(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrSamePassword):
			utils.ErrorJSON(w, apiErr.ErrSamePassword, http.StatusBadRequest)
			return
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SendPasswordResetEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in send password reset email handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.SendPasswordResetEmailSanitize()

	if err := h.userUsecase.SendPasswordResetEmail(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			// return the same status when email is not found, prevent
			// cyber attacks from brute forcing and retrieving valid emails
			utils.WriteNoContent(w, http.StatusNoContent)
			return
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in update user handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UpdateUserSanitize()

	if err := h.userUsecase.UpdateUser(ctx, userID, req); err != nil {
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID, err := utils.SessionIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	csrfToken, err := h.userUsecase.ExtendUserSessionInRedis(ctx, sessionID, utils.SESSION_EXPIRY)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// set csrf token in response headers because when user refreshes the page,
	// in memory csrf token in Pinia store (frontend) loses its state
	w.Header().Set("X-CSRF-Token", csrfToken)

	// TODO: UPDATE TO get user details
	utils.WriteNoContent(w, http.StatusOK)
}
