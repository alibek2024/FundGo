package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateCampaignInput struct {
	Title        string          `json:"title" validate:"required,min=3,max=75"`
	CreatorID    int64           `json:"-" validate:"required,gt=0"`
	Description  string          `json:"description" validate:"required,min=5"`
	TargetAmount decimal.Decimal `json:"targetAmount" validate:"required"`
	EndDate      time.Time       `json:"endDate" validate:"required"`
}

type CampaignBalanceOperation struct {
	ID     int64           `json:"id" validate:"required,gt=0"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

type PaginationParams struct {
	Limit  int32
	Offset int32
}
