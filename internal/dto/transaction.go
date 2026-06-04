package dto

import "github.com/shopspring/decimal"

type TransactionInput struct {
	UserID        int64
	DonationID    *int64
	Type          string
	Amount        decimal.Decimal
	BalanceBefore decimal.Decimal
	BalanceAfter  decimal.Decimal
}

type BalanceOperationInput struct {
	ID     int64
	Amount decimal.Decimal
}