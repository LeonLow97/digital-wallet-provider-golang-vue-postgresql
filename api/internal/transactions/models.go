package transactions

import "time"

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	SenderName                string    `json:"sender_name"`
	BeneficiaryName           string    `json:"beneficiary_name"`
	AmountTransferred         float64   `json:"amount_transferred"`
	AmountTransferredCurrency string    `json:"amount_transferred_currency"`
	AmountReceived            float64   `json:"amount_received"`
	AmountReceivedCurrency    string    `json:"amount_received_currency"`
	Status                    string    `json:"status"`
	DateTransferred           time.Time `json:"date_transferred"`
	DateReceived              time.Time `json:"date_received"`
}
