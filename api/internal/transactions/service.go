package transactions

import (
	"context"
	"strings"
	"time"

	"github.com/LeonLow97/internal/utils"
)

type Service interface {
	CreateTransaction(ctx context.Context, userId int, transaction *CreateTransaction) error
	GetTransactions(ctx context.Context, userId, page, pageSize int) (*Transactions, int, bool, error)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

// Pagination on GetTransactions which returns a list of transactions
func (s *service) GetTransactions(ctx context.Context, userId, page, pageSize int) (*Transactions, int, bool, error) {
	// handle edge cases
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// calculate total number of pages and find the last page
	totalRecords, err := s.repo.GetTransactionsCountByUserId(ctx, userId)
	if err != nil {
		return nil, 0, false, err
	}

	totalPages := (totalRecords + pageSize - 1) / pageSize
	isLastPage := page >= totalPages

	if isLastPage {
		// if user specifies a very high page number, we default it to the last page for better user experience
		offset = (totalPages - 1) * pageSize
	}

	transactions, err := s.repo.GetTransactionsByUserId(ctx, userId, pageSize, offset)
	if err != nil {
		return nil, 0, false, err
	}

	return transactions, totalPages, isLastPage, nil
}

func (s *service) CreateTransaction(ctx context.Context, userId int, transaction *CreateTransaction) error {
	if transaction.TransferredAmount == 0 || transaction.TransferredAmount < 10 || transaction.TransferredAmount > 10000 {
		return utils.BadRequestError{Message: "Transfer amount must be between $10 and $10,000."}
	}
	if err := ValidateFloatPrecision(transaction.TransferredAmount); err != nil {
		return utils.BadRequestError{Message: "Transfer amount must be up to 2 decimal places."}
	}

	// trim necessary strings
	transaction.BeneficiaryNumber = strings.TrimSpace(transaction.BeneficiaryNumber)
	transaction.TransferredAmountCurrency = strings.TrimSpace(transaction.TransferredAmountCurrency)

	// check if sender is a registered user
	userCount, err := s.repo.GetUserCountByUserId(ctx, userId)
	if err != nil {
		return err
	}
	if userCount != 1 {
		return utils.BadRequestError{Message: "The specified sender does not exist."}
	}

	// check if sender is sending to himself
	beneficiaryId, err := s.repo.GetUserIdByMobileNumber(ctx, transaction.BeneficiaryNumber)
	if err != nil {
		return err
	}
	if beneficiaryId == 0 {
		return utils.BadRequestError{Message: "The specified beneficiary does not exist."}
	}
	if userId == beneficiaryId {
		return utils.BadRequestError{Message: "Unable to send money to yourself."}
	}

	// check if sender is linked to the beneficiary
	isLinked, err := s.repo.GetCountByUserIdAndBeneficiaryId(ctx, userId, beneficiaryId)
	if err != nil {
		return err
	}
	if isLinked != 1 {
		return utils.BadRequestError{Message: "Unable to transfer funds. Sender is not linked to the specified beneficiary."}
	}

	// valid beneficiary up to this point
	// check if transferred currency exists in beneficiary's list of currencies
	var beneficiaryBalanceId int
	beneficiaryCurrency := transaction.TransferredAmountCurrency
	beneficiaryHasTransferredCurrency, beneficiaryBalanceId, err := s.repo.GetCountByUserIdAndCurrency(ctx, beneficiaryId, transaction.TransferredAmountCurrency)
	if err != nil {
		return err
	}

	// if transferred currency is not in beneficiary's list of currencies, retrieve primary currency
	if beneficiaryHasTransferredCurrency == 0 {
		beneficiaryBalanceId, beneficiaryCurrency, err = s.repo.GetBalanceIdByUserIdAndPrimary(ctx, beneficiaryId)
		if err != nil {
			return err
		}
	}

	// retrieve user balance. check for user currency availability and
	// if user has sufficient funds for the transfer.
	count, userBalanceId, err := s.repo.GetCountByUserIdAndCurrency(ctx, userId, transaction.TransferredAmountCurrency)
	if err != nil {
		return err
	}
	if count != 1 {
		return utils.BadRequestError{Message: "You do not have balance in the specified currency. Please use another currency."}
	}

	// create transaction entries for sender and beneficiary
	senderTransaction := TransactionEntity{
		UserId:                    userId,
		SenderId:                  userId,
		BalanceId:                 userBalanceId,
		BeneficiaryId:             beneficiaryId,
		TransferredAmount:         transaction.TransferredAmount,
		TransferredAmountCurrency: transaction.TransferredAmountCurrency,
		ReceivedAmountCurrency:    beneficiaryCurrency,
		Status:                    utils.TRANSACTION_STATUS.COMPLETED,
		TransferredDate:           time.Now(),
		ReceivedDate:              time.Now(),
	}

	beneficiaryTransaction := TransactionEntity{
		UserId:                            beneficiaryId,
		SenderId:                          userId,
		BalanceId:                         beneficiaryBalanceId,
		BeneficiaryId:                     beneficiaryId,
		TransferredAmount:                 transaction.TransferredAmount,
		TransferredAmountCurrency:         transaction.TransferredAmountCurrency,
		ReceivedAmountCurrency:            beneficiaryCurrency,
		BeneficiaryHasTransferredCurrency: beneficiaryHasTransferredCurrency,
		Status:                            utils.TRANSACTION_STATUS.RECEIVED,
		TransferredDate:                   time.Now(),
		ReceivedDate:                      time.Now(),
	}

	// CreateTransaction SQL Transaction
	err = s.repo.CreateTransactionSQLTransaction(ctx, senderTransaction, beneficiaryTransaction)
	if err != nil {
		return err
	}

	return nil
}
