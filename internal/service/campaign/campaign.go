package campaign

import (
	"context"
	"fmt"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/shopspring/decimal"
)

type CampaignService struct {
	Store repository.Store
}
func NewCampaignService(store repository.Store) *CampaignService {
	return &CampaignService{
		Store: store,
	}
}

func (c *CampaignService) CreateCampaign(
	ctx context.Context,
	input model.CreateCampaignInput,
) (*model.Campaign, error) {
	campaignParams := c.CampaignParams(input)
	postCampaign, err := c.Store.CreateCampaign(ctx, campaignParams)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("campaign with title '%s' already exists", input.Title)
		}
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	campaign := c.CampaignResponse(postCampaign)

	return campaign, nil
}

func (c *CampaignService) GetCampaignByID(ctx context.Context, id int32) (*model.Campaign, error) {
	postCampaign, err := c.Store.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, err
	}
	campaign := c.CampaignResponse(postCampaign)
	return campaign, nil
}

func (c *CampaignService) GetCurrentAmount(ctx context.Context, id int32) (*decimal.Decimal, error) {
	сurrentAmount, err := c.Store.GetCurrentAmount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &сurrentAmount, nil
}

func (c *CampaignService) DeleteCampaign(ctx context.Context, id int32) error {
	err := c.Store.DeleteCampaign(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CampaignService) ListCampaigns(ctx context.Context, pagination model.PaginationParams) ([]*model.Campaign, error) {
	params := postgres.ListCampaignsParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	}
	postList, err := c.Store.ListCampaigns(ctx, params)
	if err != nil {
		return nil, err
	}
	result := c.MapCampaignList(postList)
	return result, nil
}
