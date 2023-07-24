package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LeonLow97/internal/utils"
)

func (h userHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	username, err := utils.RetrieveJWTClaimsUsername(request)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
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

	user, err := h.service.GetUser(ctx, username)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unauthorized"))
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling user to JSON: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonData)
}
