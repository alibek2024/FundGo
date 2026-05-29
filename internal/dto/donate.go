package dto

import "github.com/shopspring/decimal"

type DonateInput struct {
	UserID     int64           `json:"userID" db:"userID"`
	CampaignID int64           `json:"campaignID" db:"campaignID"`
	Amount     decimal.Decimal `json:"Amount" db:"Amount"`
}
