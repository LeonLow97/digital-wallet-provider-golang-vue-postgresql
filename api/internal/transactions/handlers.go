package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/utils"
)

type transactionHandler struct {
	service Service
}

func NewTransactionHandler(s Service) (*transactionHandler, error) {
	return &transactionHandler{
		service: s,
	}, nil
}

type transactionRequest struct {
	BeneficiaryName           string `json:"beneficiary_name"`
	BeneficiaryNumber         string `json:"mobile_number"`
	AmountTransferred         string `json:"amount_transferred"`
	AmountTransferredCurrency string `json:"amount_transferred_currency"`
}

func (t transactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// retrieve userId from request context
	userId := r.Context().Value(utils.ContextUserIdKey).(int)

	// Default page and page size if not provided
	page := 1
	pageSize := 20

	// retrieve page and pageSize from url query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, errors.New("Please provide numerical page in URL Param."))
		return
	}

	pageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, errors.New("Please provide numerical page size in URL Param."))
		return
	}

	transactions, totalPages, isLastPage, err := t.service.GetTransactions(r.Context(), userId, page, pageSize)
	if err != nil {
		log.Println(err)
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, envelope{"total_pages": totalPages, "isLastPage": isLastPage, "transactions": transactions})
}

func (t transactionHandler) CreateTransaction(writer http.ResponseWriter, request *http.Request) {
	username, err := utils.RetrieveJWTClaimsUsername(request)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
	}

	decoder := json.NewDecoder(request.Body)
	var req transactionRequest
	err = decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setting timeout for requests
	ctx := request.Context()
	ctx, cancel := context.WithTimeout(ctx, utils.CONSTANTS.TIMEOUT)
	defer cancel()

	err = t.service.CreateTransaction(ctx, username, req.BeneficiaryName, req.BeneficiaryNumber, req.AmountTransferredCurrency, req.AmountTransferred)
	if err != nil {
		if s, ok := err.(utils.ServiceError); ok {
			log.Printf("Service Error: %s", err)
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(s.Error()))
			return
		}
		log.Printf("Internal Server Error: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if request.Method != http.MethodPost {
		log.Printf("Error marshaling transactions to JSON: %s", err)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}
