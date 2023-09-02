package transactions

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/LeonLow97/internal/utils"
)

var TRANSACTION_STATUS = struct {
	CONFIRMED string
	PENDING   string
	RECEIVED  string
}{
	CONFIRMED: "CONFIRMED",
	PENDING:   "PENDING",
	RECEIVED:  "RECEIVED",
}

type Service interface {
	CreateTransaction(ctx context.Context, senderName, beneficiaryName, beneficiaryNumber, amountTransferredCurrency, amountTransferredString string) error
	GetTransactions(ctx context.Context, userId int) (*Transactions, error)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

func (s *service) GetTransactions(ctx context.Context, userId int) (*Transactions, error) {
	transactions, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Sorting date_transferred by descending order
	sort.Slice(transactions.Transactions, func(i, j int) bool {
		return transactions.Transactions[i].DateTransferred.After(transactions.Transactions[j].DateTransferred)
	})

	return transactions, nil
}

func (s *service) CreateTransaction(ctx context.Context, senderName, beneficiaryName, beneficiaryNumber, amountTransferredCurrency, amountTransferredString string) error {
	amountTransferred, err := strconv.ParseFloat(amountTransferredString, 64)
	if err != nil {
		return utils.ServiceError{Message: "Please enter a valid numeric amount."}
	}
	if amountTransferred == 0 {
		return utils.ServiceError{Message: "Please specify an amount to be transferred."}
	}
	if amountTransferred < 10 {
		return utils.ServiceError{Message: "Minimum amount allowed per transfer is 10. Please try again."}
	}
	if amountTransferred > 10000 {
		return utils.ServiceError{Message: "Maximum amount allowed per transfer is 10000. Please try again."}
	}
	err = utils.ValidateFloatPrecision(amountTransferred)
	if err != nil {
		return utils.ServiceError{Message: "Amount Transferred must be up to 2 decimal places."}
	}
	senderName = strings.TrimSpace(senderName)
	beneficiaryName = strings.TrimSpace(beneficiaryName)
	beneficiaryNumber = strings.TrimSpace(beneficiaryNumber)
	amountTransferredCurrency = strings.TrimSpace(amountTransferredCurrency)

	// -------------------- MONEY TRANSFER --------------------
	// 1. Determine if sender is a registered user by verifying username
	count, err := s.repo.GetCountByUsername(ctx, senderName)
	if err != nil {
		return err
	}
	if count == 0 {
		return utils.ServiceError{Message: "This sender does not exist."}
	}
	if count > 1 {
		return utils.ServiceError{Message: "Duplicate senders."}
	}

	// Check if sender is sending to himself (not allowed)

	// 2. Determine if sender is linked to beneficiary by sender name and beneficiary mobile number
	count, err = s.repo.GetCountByUsernameAndBeneficiaryNumber(ctx, senderName, beneficiaryNumber)
	if err != nil {
		return err
	}
	if count == 0 {
		return utils.ServiceError{Message: "This user is not linked to the specified beneficiary"}
	}

	// 2. Determine if the beneficiary is a registered user of the mobile application via mobile number
	count, err = s.repo.GetCountByBeneficiaryNameAndBeneficiaryNumber(ctx, beneficiaryName, beneficiaryNumber)
	if err != nil {
		return err
	}
	if count == 0 {
		// Check if mobile number exist in the `users` table, if exist, it is a registered user and the sender may
		// be referring to another name
		username, err := s.repo.GetUsernameByBeneficiaryNumber(ctx, beneficiaryNumber)
		if err != nil {
			return err
		}
		if len(username) > 0 {
			return utils.ServiceError{Message: "Do you mean: " + username + " ?"}
		}

		log.Println("External Beneficiary...")
		// External Beneficiary! For future development!
		return nil
	}
	if count > 1 {
		return utils.ServiceError{Message: "Duplicate beneficiary"}
	}

	// 3. Get Beneficiary Currency
	amountReceivedCurrency, err := s.repo.GetCurrencyByBeneficiaryMobileNumber(ctx, beneficiaryNumber)
	if err != nil {
		return err
	}

	// 4. Calculate Beneficiary received amount (Perform currency conversion in backend)
	amountReceived := utils.CurrencyConversion(amountTransferred, amountTransferredCurrency, amountReceivedCurrency)

	// Check if amount to be sent if less than or equal to the user's balance (SELECT)
	// userBalance, err := s.repo.GetUserBalanceByUsername(ctx, senderName)
	// if err != nil {
	// 	return err
	// }

	// if userBalance == 0.0 {
	// 	return utils.ServiceError{Message: "Account has 0 balance. Please top up."}
	// }
	// if amountTransferred > userBalance {
	// 	return utils.ServiceError{Message: "User does not have sufficient funds to make the transfer. Please top up."}
	// }

	// 6. Deduct the transferred amount from the current user balance to obtain the final balance
	// finalUserBalance := userBalance - amountTransferred

	// 7. Update the beneficiary balance by adding the amount transferred to the current balance
	// beneficiaryBalance, err := s.repo.GetUserBalanceByUsername(ctx, beneficiaryName)
	// if err != nil {
	// 	return err
	// }
	// finalBeneficiaryBalance := beneficiaryBalance + amountReceived

	// 8. Perform SQL Transaction to ensure data integrity and
	// that the amounts were deducted and transferred to the correct recipients.
	err = s.repo.SQLTransactionMoneyTransfer(ctx, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, TRANSACTION_STATUS.CONFIRMED, TRANSACTION_STATUS.RECEIVED, amountTransferred, amountReceived)
	if err != nil {
		return err
	}

	return nil
}
