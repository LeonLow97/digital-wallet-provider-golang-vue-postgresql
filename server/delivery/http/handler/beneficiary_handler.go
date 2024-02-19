package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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
	beneficiaryRouter.HandleFunc("", handler.UpdateBeneficiary).Methods(http.MethodPatch)
	beneficiaryRouter.HandleFunc("/{id:[0-9]+}", handler.GetBeneficiary).Methods(http.MethodGet)
	beneficiaryRouter.HandleFunc("", handler.GetBeneficiaries).Methods(http.MethodGet)
}

func (h *BeneficiaryHandler) CreateBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateBeneficiaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in create beneficiary handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
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
	req.UserID = userID

	err = h.beneficiaryUsecase.CreateBeneficiary(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserNotFound):
		utils.ErrorJSON(w, apiErr.ErrUserNotFound, http.StatusNotFound)
	case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
		utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
	case errors.Is(err, exception.ErrBeneficiaryAlreadyExists):
		utils.ErrorJSON(w, apiErr.ErrBeneficiaryAlreadyExists, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteNoContent(w, http.StatusCreated)
	}
}

func (h *BeneficiaryHandler) UpdateBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.UpdateBeneficiaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding req body in update beneficiary handler", err)
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	}

	errMessage, err := infrastructure.ValidateStruct(req)
	if err != nil {
		log.Println("error validating req struct in create beneficiary handler", err)
		utils.ErrorJSON(w, errMessage, http.StatusBadRequest)
		return
	}

	req.UserID = userID

	err = h.beneficiaryUsecase.UpdateBeneficiary(ctx, req)
	switch {
	case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
		utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
		utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteNoContent(w, http.StatusNoContent)
	}
}

func (h *BeneficiaryHandler) GetBeneficiary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve beneficiary id from url params
	var beneficiaryID int
	vars := mux.Vars(r)
	if beneficiaryIDString, ok := vars["id"]; !ok {
		log.Println("unable to get beneficiary id from url params")
		utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return
	} else {
		id, err := strconv.Atoi(beneficiaryIDString)
		if err != nil {
			log.Println("Unable to convert beneficiary ID to string")
			utils.ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
			return
		}
		beneficiaryID = id
	}

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// get one beneficiary
	resp, err := h.beneficiaryUsecase.GetBeneficiary(ctx, beneficiaryID, userID)
	switch {
	case errors.Is(err, exception.ErrUserIDEqualBeneficiaryID):
		utils.ErrorJSON(w, apiErr.ErrUserIDEqualBeneficiaryID, http.StatusBadRequest)
	case errors.Is(err, exception.ErrUserNotLinkedToBeneficiary):
		utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToBeneficiary, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (h *BeneficiaryHandler) GetBeneficiaries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user id from context
	userID, ok := ctx.Value(utils.UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// get beneficiaries
	resp, err := h.beneficiaryUsecase.GetBeneficiaries(ctx, userID)
	switch {
	case errors.Is(err, exception.ErrUserHasNoBeneficiary):
		utils.ErrorJSON(w, apiErr.ErrUserNotLinkedToAnyBeneficiary, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}
