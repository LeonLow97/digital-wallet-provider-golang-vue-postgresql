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
	"github.com/LeonLow97/go-clean-architecture/utils/constants"
	"github.com/LeonLow97/go-clean-architecture/utils/contextstore"
	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(uc domain.UserUsecase) *UserHandler {
	handler := &UserHandler{
		userUsecase: uc,
	}

	return handler
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.LoginRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in login handler", err, errMessage)
		jsonutil.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	req.LoginSanitize()

	resp, err := h.userUsecase.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound) || errors.Is(err, exception.ErrInvalidCredentials):
			jsonutil.ErrorJSON(w, apiErr.ErrInvalidCredentials, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrInactiveUser):
			jsonutil.ErrorJSON(w, apiErr.ErrInactiveUser, http.StatusUnauthorized)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SignUpRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	req.SignUpSanitize()

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in sign up handler", errMessage)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	if err = h.userUsecase.SignUp(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUserFound, http.StatusBadRequest)
		case errors.Is(err, exception.ErrInvalidPassword):
			jsonutil.ErrorJSON(w, apiErr.ErrInvalidPassword, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID, err := contextstore.SessionIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
	}

	// remove sessionID from Redis
	h.userUsecase.RemoveSessionFromRedis(r.Context(), sessionID)

	// override cookie in browser
	cookie := &http.Cookie{
		Name:     constants.JWT_COOKIE,
		Value:    "",
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		Secure:   false, // For HTTPS, `Secure: true`. Using HTTP, so `Secure: false`
		HttpOnly: true,  // prevent client-side scripts from accessing cookie, like `document.cookie`
	}
	http.SetCookie(w, cookie)

	jsonutil.WriteNoContent(w, http.StatusOK)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.ChangePasswordRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.ChangePasswordSanitize()

	if err := h.userUsecase.ChangePassword(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUserNotFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrInvalidCredentials):
			jsonutil.ErrorJSON(w, apiErr.ErrCurrentPasswordIncorrect, http.StatusBadRequest)
		case errors.Is(err, exception.ErrSamePassword):
			jsonutil.ErrorJSON(w, apiErr.ErrSamePassword, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) ConfigureMFA(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.ConfigureMFARequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.ConfigureMFASanitize()

	token, err := h.userUsecase.ConfigureMFA(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidMFACode):
			jsonutil.ErrorJSON(w, apiErr.ErrInvalidMFACode, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrUserNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrTOTPSecretExists):
			jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("X-CSRF-Token", token.CSRFToken)
	utils.IssueCookie(w, token.AccessToken)

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.VerifyMFARequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.VerifyMFASanitize()

	token, err := h.userUsecase.VerifyMFA(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, exception.ErrInvalidMFACode):
			jsonutil.ErrorJSON(w, apiErr.ErrInvalidMFACode, http.StatusUnauthorized)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("X-CSRF-Token", token.CSRFToken)
	utils.IssueCookie(w, token.AccessToken)

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.PasswordResetRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.PasswordResetSanitize()

	if err := h.userUsecase.PasswordReset(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrSamePassword):
			jsonutil.ErrorJSON(w, apiErr.ErrSamePassword, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SendPasswordResetEmailRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.SendPasswordResetEmailSanitize()

	if err := h.userUsecase.SendPasswordResetEmail(ctx, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			// return the same status when email is not found, prevent
			// cyber attacks from brute forcing and retrieving valid emails
			jsonutil.WriteNoContent(w, http.StatusNoContent)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.UpdateUserRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UpdateUserSanitize()

	if err := h.userUsecase.UpdateUser(ctx, userID, req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *UserHandler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID, err := contextstore.SessionIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	csrfToken, err := h.userUsecase.ExtendUserSessionInRedis(ctx, sessionID, constants.SESSION_EXPIRY)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// set csrf token in response headers because when user refreshes the page,
	// in memory csrf token in Pinia store (frontend) loses its state
	w.Header().Set("X-CSRF-Token", csrfToken)

	// TODO: UPDATE TO get user details
	jsonutil.WriteNoContent(w, http.StatusOK)
}
