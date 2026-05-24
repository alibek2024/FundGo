package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionDeposit  TransactionType = "deposit"
	TransactionWithdraw TransactionType = "withdraw"
	TransactionDonation TransactionType = "donation"
	TransactionTransfer TransactionType = "transfer"
	TransactionRefund   TransactionType = "refund"
	TransactionTopUp TransactionType = "purchase"
)

type Transaction struct {
	ID         int32
	UserID     int32
	DonationID *int32
	Type       TransactionType
	Amount     decimal.Decimal
	CreatedAt  time.Time
}

type TransactionInput struct {
	UserID     int32
	DonationID *int32
	Type       TransactionType
	Amount     decimal.Decimal
}
