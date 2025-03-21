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

type BeneficiaryHandler struct {
	beneficiaryUsecase domain.BeneficiaryUsecase
}

func NewBeneficiaryHandler(uc domain.BeneficiaryUsecase) *BeneficiaryHandler {
	handler := &BeneficiaryHandler{
		beneficiaryUsecase: uc,
	}

	return handler
}

func (h *BeneficiaryHandler) CreateBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateBeneficiaryRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create beneficiary handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	// sanitize request body
	req.CreateBeneficiarySanitize()

	if err = h.beneficiaryUsecase.CreateBeneficiary(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			jsonutil.ErrorJSON(w, apiErr.ErrUserNotFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			jsonutil.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrBeneficiaryAlreadyExists):
			jsonutil.ErrorJSON(w, apiErr.ErrBeneficiaryAlreadyExists, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusCreated)
}

func (h *BeneficiaryHandler) UpdateBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.UpdateBeneficiaryRequest
	if err := jsonutil.ReadJSONBody(w, r, &req); err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create beneficiary handler", err)
		jsonutil.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	if err = h.beneficiaryUsecase.UpdateBeneficiary(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			jsonutil.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
			jsonutil.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteNoContent(w, http.StatusNoContent)
}

func (h *BeneficiaryHandler) GetBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve beneficiary id from url params
	beneficiaryID, err := jsonutil.ReadURLParamsInt(w, r, "id")
	if err != nil {
		return
	}

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// get one beneficiary
	resp, err := h.beneficiaryUsecase.GetBeneficiary(ctx, beneficiaryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			jsonutil.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
			jsonutil.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (h *BeneficiaryHandler) GetBeneficiaries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// get beneficiaries
	resp, err := h.beneficiaryUsecase.GetBeneficiaries(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserHasNoBeneficiary):
			jsonutil.ErrorJSON(w, apiErr.ErrUserNotLinkedToAnyBeneficiary, http.StatusBadRequest)
		default:
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, resp)
}
