package postgres

import (
	"context"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/shopspring/decimal"
)

func (c *Repository) CreateCampaign(
	ctx context.Context,
	input dto.CreateCampaignInput,
) (*model.Campaign, error) {
	campaignParams := mapper.CampaignParams(input)
	postCampaign, err := c.DB.CreateCampaign(ctx, campaignParams)
	if err != nil {
		if mapper.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("campaign with title '%s' already exists", input.Title)
		}
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	campaign := mapper.CampaignResponse(postCampaign)

	return campaign, nil
}

func (c *Repository) GetCampaignByID(ctx context.Context, id int64) (*model.Campaign, error) {
	postCampaign, err := c.DB.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, err
	}
	campaign := mapper.CampaignResponse(postCampaign)
	return campaign, nil
}

func (c *Repository) GetCampaignByTitle(ctx context.Context, title string) (*model.Campaign, error) {
	postCampaign, err := c.DB.GetCampaignByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	campaign := mapper.CampaignResponse(postCampaign)
	return campaign, nil
}

func (c *Repository) GetCurrentAmount(ctx context.Context, id int64) (*decimal.Decimal, error) {
	currentAmount, err := c.DB.GetCurrentAmount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &currentAmount, nil
}
func (c *Repository) GetCampaignStatus(ctx context.Context, id int64) (*model.CampaignStatus, error) {
	storeStatus, err := c.DB.GetCampaignStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	status := model.CampaignStatus(storeStatus)
	return &status, nil
}

func (c *Repository) DeleteCampaign(ctx context.Context, id int64) error {
	err := c.DB.DeleteCampaign(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Repository) ListCampaigns(ctx context.Context, pagination dto.PaginationParams) ([]*model.Campaign, error) {
	params := mapper.PaginationParams(pagination)
	postList, err := c.DB.ListCampaigns(ctx, params)
	if err != nil {
		return nil, err
	}
	result := mapper.MapCampaignList(postList)
	return result, nil
}

func (c *Repository) IncreaseCampaignAmount(ctx context.Context, input dto.CampaignBalanceOperation) (*model.Campaign, error) {
	params := mapper.IncreaseCampaignAmount(input)
	storeAmount, err := c.DB.IncreaseCampaignAmount(ctx, params)
	if err != nil {
	  return nil, err
	}
	campaign := mapper.CampaignResponse(storeAmount)
	return campaign, nil
}
