package repo

import (
	"context"
	"fmt"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/shopspring/decimal"
)

type Service struct {
	Store store.SQLStore
}

func NewCampaignService(store store.SQLStore) *Service {
	return &Service{
		Store: store,
	}
}

func (c *Service) CreateCampaign(
	ctx context.Context,
	input model.CreateCampaignInput,
) (*model.Campaign, error) {
	campaignParams := mapper.CampaignParams(input)
	postCampaign, err := c.Store.DB.CreateCampaign(ctx, campaignParams)
	if err != nil {
		if mapper.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("campaign with title '%s' already exists", input.Title)
		}
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	campaign := mapper.CampaignResponse(postCampaign)

	return campaign, nil
}

func (c *Service) GetCampaignByID(ctx context.Context, id int64) (*model.Campaign, error) {
	postCampaign, err := c.Store.DB.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, err
	}
	campaign := mapper.CampaignResponse(postCampaign)
	return campaign, nil
}

func (c *Service) GetCurrentAmount(ctx context.Context, id int64) (*decimal.Decimal, error) {
	currentAmount, err := c.Store.DB.GetCurrentAmount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &currentAmount, nil
}

func (c *Service) DeleteCampaign(ctx context.Context, id int64) error {
	err := c.Store.DB.DeleteCampaign(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Service) ListCampaigns(ctx context.Context, pagination model.PaginationParams) ([]*model.Campaign, error) {
	params := mapper.PaginationParams(pagination)
	postList, err := c.Store.DB.ListCampaigns(ctx, params)
	if err != nil {
		return nil, err
	}
	result := mapper.MapCampaignList(postList)
	return result, nil
}
