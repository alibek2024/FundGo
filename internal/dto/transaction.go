package dto

import "github.com/shopspring/decimal"

type TransactionInput struct {
	UserID        int64           `json:"userID" schema:"user_id" validate:"required,gt=0"`
	DonationID    *int64          `json:"donationID" schema:"donation_id" validate:"required,gt=0"`
	Type          string          `json:"type" schema:"type" validate:"required,gt=0"`
	Amount        decimal.Decimal `json:"amount" schema:"amount" validate:"required,gt=0"`
	BalanceBefore decimal.Decimal `json:"-"`
	BalanceAfter  decimal.Decimal `json:"-"`
}

type BalanceOperationInput struct {
	ID     int64           `json:"-" validate:"required,gt=0"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}
