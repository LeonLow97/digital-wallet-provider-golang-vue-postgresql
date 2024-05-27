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

type BeneficiaryHandler struct {
	beneficiaryUsecase domain.BeneficiaryUsecase
}

func NewBeneficiaryHandler(router *mux.Router, uc domain.BeneficiaryUsecase) {
	handler := &BeneficiaryHandler{
		beneficiaryUsecase: uc,
	}

	beneficiaryRouter := router.PathPrefix("/beneficiary").Subrouter()

	beneficiaryRouter.HandleFunc("", handler.CreateBeneficiary).Methods(http.MethodPost)
	beneficiaryRouter.HandleFunc("", handler.UpdateBeneficiary).Methods(http.MethodPut)
	beneficiaryRouter.HandleFunc("/{id:[0-9]+}", handler.GetBeneficiary).Methods(http.MethodGet)
	beneficiaryRouter.HandleFunc("", handler.GetBeneficiaries).Methods(http.MethodGet)
}

func (h *BeneficiaryHandler) CreateBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateBeneficiaryRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create beneficiary handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	// sanitize request body
	req.CreateBeneficiarySanitize()

	if err = h.beneficiaryUsecase.CreateBeneficiary(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			utils.ErrorJSON(w, apiErr.ErrUserNotFound, http.StatusNotFound)
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrBeneficiaryAlreadyExists):
			utils.ErrorJSON(w, apiErr.ErrBeneficiaryAlreadyExists, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusCreated)
}

func (h *BeneficiaryHandler) UpdateBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.UpdateBeneficiaryRequest
	if err := utils.ReadJSONBody(w, r, &req); err != nil {
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create beneficiary handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	if err = h.beneficiaryUsecase.UpdateBeneficiary(ctx, userID, req); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
			utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteNoContent(w, http.StatusNoContent)
}

func (h *BeneficiaryHandler) GetBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve beneficiary id from url params
	beneficiaryID, err := utils.ReadParamsInt(w, r, "id")
	if err != nil {
		return
	}

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// get one beneficiary
	resp, err := h.beneficiaryUsecase.GetBeneficiary(ctx, beneficiaryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
			utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
		case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
			utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *BeneficiaryHandler) GetBeneficiaries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// get beneficiaries
	resp, err := h.beneficiaryUsecase.GetBeneficiaries(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserHasNoBeneficiary):
			utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToAnyBeneficiary, http.StatusBadRequest)
		default:
			utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
