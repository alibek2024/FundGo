package campaign

import (
	"errors"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (c *CampaignService) defaultParams(input model.CreateCampaignInput) (
	int64, string, pgtype.Text, decimal.Decimal) {

	return input.CreatorID,
		input.Title,
		helperfunc.Text(input.Description),
		input.TargetAmount
}

func (c *CampaignService) CampaignParams(
	input model.CreateCampaignInput,
) postgres.CreateCampaignParams {
	id, title, description, TargetAmount := c.defaultParams(input)
	return postgres.CreateCampaignParams{
		CreatorID:    id,
		Title:        title,
		Description:  description,
		TargetAmount: TargetAmount,
	}
}

func (u *CampaignService) CampaignResponse(input postgres.Campaign) *model.Campaign {
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

func (u *CampaignService) MapCampaignList(input []postgres.Campaign) []*model.Campaign {
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

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// "23505" — это код уникального нарушения в PostgreSQL
		return pgErr.Code == "23505"
	}
	return false
}
