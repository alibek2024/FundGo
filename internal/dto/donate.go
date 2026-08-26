package dto

import "github.com/shopspring/decimal"

type DonateInput struct {
	UserID     int64           `json:"-" validate:"required,gt=0"`
	CampaignID int64           `json:"campaignID" validate:"required,gt=0"`
	Amount     decimal.Decimal `json:"amount" validate:"decimal_gt_zero"`
}

type UpdateDonationStatus struct {
	DonationID int64  `json:"donationID" validate:"required,gt=0"`
	Status     string `json:"status" validate:"required"`
}
