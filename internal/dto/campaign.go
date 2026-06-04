package dto

import "github.com/shopspring/decimal"

type CreateCampaignInput struct {
	Title        string          `json:"title"`
	CreatorID    int64           `json:"creatorID"`
	Description  string          `json:"description"`
	TargetAmount decimal.Decimal `json:"targetAmount"`
}

type CampaignBalanceOperation struct {
	ID            int64           `json:"id"`
	Amount decimal.Decimal `json:"currentAmount"`
}

type PaginationParams struct {
	Limit  int32
	Offset int32
}
