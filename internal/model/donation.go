package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Donation struct {
	ID         int64           `json:"id" db:"id"`
	UserID     int64          `json:"userID" db:"userID"`
	CampaignID int64           `json:"campaignID" db:"campaignID"`
	Amount     decimal.Decimal `json:"Amount" db:"Amount"`
	CreatedAt  time.Time       `json:"CreatedAt" db:"CreatedAt"`
}

type DonateInput struct {
	UserID     int64           `json:"userID" db:"userID"`
	CampaignID int64           `json:"campaignID" db:"campaignID"`
	Amount     decimal.Decimal `json:"Amount" db:"Amount"`
}
