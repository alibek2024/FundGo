package dto

import "github.com/shopspring/decimal"

type DonateInput struct {
	UserID     int64           `json:"userID" schema:"user_id" validate:"required,gt=0"`
	CampaignID int64           `json:"campaignID" schema:"campaign_id" validate:"required,gt=0"`
	Amount     decimal.Decimal `json:"Amount" schema:"amount" validate:"required"`
}
type UpdateDonationStatus struct {
	DonationID int64  `json:"donation_id" schema:"donation_id" validate:"required,gt=0"`
	Status     string `json:"status" schema:"status" validate:"required"`
}
