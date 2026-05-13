package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CampaignStatus string

type Campaign struct {
	Title         string          `json:"title" db:"title"`
	CreatorID     int32           `json:"creatorID" db:"creatorID"`
	ID            int32           `json:"id" db:"id"`
	Description   string          `json:"description" db:"description"`
	TargetAmount  decimal.Decimal `json:"targetAmount" db:"targetAmount"`
	CurrentAmount decimal.Decimal `json:"currentAmount" db:"currentAmount"`
	Status        CampaignStatus  `json:"status" db:"status"`
	EndDate       time.Time       `json:"endDate" db:"endDate"`
	CreatedAt     time.Time       `json:"createdAt" db:"createdAt"`
}

type CreateCampaignInput struct {
	Title         string          `json:"title" db:"title"`
	CreatorID     int32           `json:"creatorID" db:"creatorID"`
	Description   string          `json:"description" db:"description"`
	TargetAmount  decimal.Decimal `json:"targetAmount" db:"targetAmount"`
}

type PaginationParams  struct {
    Limit  int32
    Offset int32
}