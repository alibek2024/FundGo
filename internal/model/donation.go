package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DonationStatus string

const (
	DonationActive   DonationStatus = "active"
	DonationArchived DonationStatus = "archived"
	DonationRefund   DonationStatus = "refund"
)

type Donation struct {
	ID         int64           `json:"id" db:"id"`
	UserID     int64           `json:"userID" db:"userID"`
	CampaignID int64           `json:"campaignID" db:"campaignID"`
	Status     DonationStatus  `json:"status" db:"status"`
	Amount     decimal.Decimal `json:"Amount" db:"Amount"`
	CreatedAt  time.Time       `json:"CreatedAt" db:"CreatedAt"`
}
