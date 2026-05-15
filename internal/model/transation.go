package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Donate     TransactionType = "donation"
	Refund     TransactionType = "refund"
	Withdrawal TransactionType = "withdrawal"
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
	Amount     decimal.Decimal
}
