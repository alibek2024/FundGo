package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CampaignStatus string

const (
	Active     CampaignStatus = "deposit"
	Successful CampaignStatus = "successful"
	Failed     CampaignStatus = "failed"
	Archived   CampaignStatus = "archived"
)

type Campaign struct {
	Title         string          `json:"title" db:"title"`
	CreatorID     int64           `json:"creatorID" db:"creatorID"`
	ID            int64           `json:"id" db:"id"`
	Description   string          `json:"description" db:"description"`
	TargetAmount  decimal.Decimal `json:"targetAmount" db:"targetAmount"`
	CurrentAmount decimal.Decimal `json:"currentAmount" db:"currentAmount"`
	Status        CampaignStatus  `json:"status" db:"status"`
	EndDate       time.Time       `json:"endDate" db:"endDate"`
	CreatedAt     time.Time       `json:"createdAt" db:"createdAt"`
}
