package transactions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/utils"
)

type transactionHandler struct {
	service Service
}

func NewTransactionHandler(s Service) (*transactionHandler) {
	return &transactionHandler{
		service: s,
	}
}

// Retrieves a list of transactions for the user to view the transaction history
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
		utils.ErrorJSON(w, utils.BadRequestError{Message: "Please provide numerical page in URL Param."})
		return
	}

	pageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, utils.BadRequestError{Message: "Please provide numerical page size in URL Param."})
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

// CreateTransaction handler for transferring funds to a beneficiary
func (t transactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {

	var transaction *CreateTransaction

	err := utils.ReadJSON(w, r, &transaction)
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, utils.BadRequestError{Message: "Bad Request!"})
		return
	}

	userId := r.Context().Value(utils.ContextUserIdKey).(int)

	// calling the CreateTransaction service
	err = t.service.CreateTransaction(r.Context(), userId, transaction)
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, envelope{"result": "Successfully created a transaction!"})
}
