package dto

import "github.com/shopspring/decimal"

type TransactionInput struct {
	UserID        int64           `json:"userID" validate:"required,gt=0"`
	DonationID    *int64          `json:"donationID" validate:"omitempty,gt=0"`
	Type          string          `json:"type" validate:"required"`
	Amount        decimal.Decimal `json:"amount" validate:"decimal_gt_zero"`
	BalanceBefore decimal.Decimal `json:"-"`
	BalanceAfter  decimal.Decimal `json:"-"`
}

type BalanceOperationInput struct {
	ID     int64           `json:"-" validate:"required,gt=0"`
	Amount decimal.Decimal `json:"amount" validate:"decimal_gt_zero"`
}
