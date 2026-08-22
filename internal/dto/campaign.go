package dto

import "github.com/shopspring/decimal"

type CreateCampaignInput struct {
	Title        string          `json:"title" schema:"title" validate:"required,min=3,max=75"`
	CreatorID    int64           `json:"creatorID" schema:"creator_id" validate:"required,validate:gt=0"`
	Description  string          `json:"description" schema:"description" validate:"required,min=5"`
	TargetAmount decimal.Decimal `json:"targetAmount" schema:"target_amount" validate:"required,validate:"`
}

type CampaignBalanceOperation struct {
	ID     int64           `json:"id" schema:"id" validate:"required,validate:gt=0"`
	Amount decimal.Decimal `json:"currentAmount" schema:"amount" validate:"required,validate:"`
}

type PaginationParams struct {
	Limit  int32
	Offset int32
}
