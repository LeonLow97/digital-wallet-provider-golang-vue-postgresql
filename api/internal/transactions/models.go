package transactions

import "time"

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	SenderFirstName           string    `json:"sender_first_name"`
	SenderLastName            string    `json:"sender_last_name"`
	SenderUsername            string    `json:"sender_username"`
	BeneficiaryFirstName      string    `json:"beneficiary_first_name"`
	BeneficiaryLastName       string    `json:"beneficiary_last_name"`
	BeneficiaryUsername       string    `json:"beneficiary_username"`
	TransferredAmount         float64   `json:"transferred_amount"`
	TransferredAmountCurrency string    `json:"transferred_amount_currency"`
	ReceivedAmount            float64   `json:"received_amount"`
	ReceivedAmountCurrency    string    `json:"received_amount_currency"`
	Status                    string    `json:"status"`
	TransferredDate           time.Time `json:"transferred_date"`
	ReceivedDate              time.Time `json:"received_date"`
}
