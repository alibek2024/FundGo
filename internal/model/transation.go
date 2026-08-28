package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionDeposit  TransactionType = "deposit"
	TransactionDonation TransactionType = "donation"
	TransactionRefund   TransactionType = "refund"
	TransactionWithdraw TransactionType = "withdrawal"
)

type Transaction struct {
	ID            int64
	UserID        int64
	DonationID    *int64
	Type          TransactionType
	Amount        decimal.Decimal
	BalanceBefore decimal.Decimal
	BalanceAfter  decimal.Decimal
	CreatedAt     time.Time
}
