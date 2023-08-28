package users

import (
	"log"
	"net/http"

	"github.com/LeonLow97/internal/utils"
)

type userHandler struct {
	service Service
}

func NewUserHandler(s Service) (*userHandler, error) {
	return &userHandler{
		service: s,
	}, nil
}

type envelope map[string]interface{}

func (h userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := utils.ReadJSON(w, r, &creds)
	if err != nil {
		log.Printf("invalid json supplied, or json missing entirely: %s", err)
		utils.ErrorJSON(w, utils.BadRequestError{Message: "Bad Request!"})
		return
	}

	// Invoking the Login service
	user, token, err := h.service.Login(r.Context(), &creds)
	if err != nil {
		log.Printf("Login issue: %s", err)
		utils.ErrorJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusUnauthorized, envelope{"user": user, "token": token})
}
