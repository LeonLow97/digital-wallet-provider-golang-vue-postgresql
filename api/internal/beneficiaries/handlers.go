package beneficiaries

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LeonLow97/internal/utils"
)

type beneficiaryHandler struct {
	service Service
}

func NewBeneficiaryHandler(s Service) *beneficiaryHandler {
	return &beneficiaryHandler{
		service: s,
	}
}

func (b beneficiaryHandler) GetBeneficiaries(writer http.ResponseWriter, request *http.Request) {
	username, err := utils.RetrieveJWTClaimsUsername(request)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if request.Method != http.MethodGet {
		log.Println("Invalid Method")
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Setting timeout for requests
	ctx := request.Context()
	ctx, cancel := context.WithTimeout(ctx, utils.CONSTANTS.TIMEOUT)
	defer cancel()

	beneficiaries, err := b.service.GetBeneficiaries(ctx, username)
	if err != nil {
		if s, ok := err.(utils.ServiceError); ok {
			log.Printf("Service Error: %s", err)
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(s.Error()))
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(beneficiaries)
	if err != nil {
		log.Printf("Error marshaling beneficiaries to JSON: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonData)

}

// TODO: Add Beneficiary
// TODO: Update Beneficiary
// TODO: Delete Beneficiary
