package mapper

import (
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func UpdateStatusParams(id int64, status model.CampaignStatus) generated.UpdateCampainStatusParams {
	return generated.UpdateCampainStatusParams{
		ID:     id,
		Status: generated.CampaignStatus(status),
	}
}

func DefaultCampaignParams(input dto.CreateCampaignInput) (
	int64, string, pgtype.Text, decimal.Decimal) {

	return input.CreatorID,
		input.Title,
		Text(input.Description),
		input.TargetAmount
}

func CampaignParams(input dto.CreateCampaignInput) generated.CreateCampaignParams {
	id, title, description, targetAmount := DefaultCampaignParams(input)

	return generated.CreateCampaignParams{
		CreatorID:    id,
		Title:        title,
		Description:  description,
		TargetAmount: targetAmount,
		EndDate: pgtype.Timestamptz{
			Time:  input.EndDate,
			Valid: true,
		},
	}
}

func CampaignResponse(input generated.Campaign) *model.Campaign {
	return &model.Campaign{
		ID:            input.ID,
		Title:         input.Title,
		CreatorID:     input.CreatorID,
		Description:   input.Description.String,
		TargetAmount:  input.TargetAmount,
		CurrentAmount: input.CurrentAmount,
		Status:        model.CampaignStatus(input.Status),
		EndDate:       input.EndDate.Time,
		CreatedAt:     input.CreatedAt.Time,
	}
}

func MapCampaignList(input []generated.Campaign) []*model.Campaign {
	result := make([]*model.Campaign, len(input))

	for i, c := range input {
		result[i] = &model.Campaign{
			ID:            c.ID,
			Title:         c.Title,
			CreatorID:     c.CreatorID,
			Description:   c.Description.String,
			TargetAmount:  c.TargetAmount,
			CurrentAmount: c.CurrentAmount,
			Status:        model.CampaignStatus(c.Status),
			EndDate:       c.EndDate.Time,
			CreatedAt:     c.CreatedAt.Time,
		}
	}
	return result
}

func PaginationParams(input dto.PaginationParams) generated.ListCampaignsParams {
	return generated.ListCampaignsParams{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func IncreaseCampaignAmount(input dto.CampaignBalanceOperation) generated.IncreaseCampaignAmountParams {
	return generated.IncreaseCampaignAmountParams{
		ID:            input.ID,
		CurrentAmount: input.Amount,
	}
}

func DecreaseCampaignAmount(input dto.CampaignBalanceOperation) generated.DecreaseCampaignAmountParams {
	return generated.DecreaseCampaignAmountParams{
		ID:            input.ID,
		CurrentAmount: input.Amount,
	}
}
