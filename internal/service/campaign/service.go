package campaign

import (
	"context"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
)

type Service struct {
	store store.CampaignStore
	Tx    RefundManager
}

func NewCampaignService(store store.CampaignStore, tx RefundManager) *Service {
	return &Service{
		store: store,
		Tx:    tx,
	}
}

func (c *Service) CreateCampaign(
	ctx context.Context,
	input dto.CreateCampaignInput,
) (*model.Campaign, error) {
	campaign, err := c.store.CreateCampaign(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create camapaign :%w", err)
	}

	return campaign, nil
}

func (c *Service) GetCampaignByID(
	ctx context.Context,
	campaignID int64,
) (*model.Campaign, error) {
	camapaign, err := c.store.GetCampaignByID(ctx, campaignID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, contracts.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("get campaign by id: %w", err)
	}
	return camapaign, nil
}

func (c *Service) SearchCampaign(
	ctx context.Context,
	name string,
) ([]*model.Campaign, error) {
	campaigns, err := c.store.GetCampaignByTitle(ctx, name)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, contracts.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("get campaign by id: %w", err)
	}
	return campaigns, nil
}

func (c *Service) WrapUpCampaign(ctx context.Context, campaignID int64) error {
	camapaign, err := c.store.GetCampaignByID(ctx, campaignID)
	if err != nil {
		if err == store.ErrNotFound {
			return contracts.ErrCampaignNotFound
		}
		return fmt.Errorf("get campaign by id: %w", err)
	}

	if camapaign.Status != model.Active {
		return contracts.ErrCampaignNotActive
	}

	if 0 > camapaign.CurrentAmount.Cmp(camapaign.TargetAmount) {
		err := c.Tx.RefundDonations(ctx, campaignID)
		if err != nil {
			return fmt.Errorf("refund donations: %w", err)
		}

		failed := model.Failed
		_, err = c.store.UpdateCampaignStatus(ctx, campaignID, failed)
		if err != nil {
			return fmt.Errorf("update status: %w", err)
		}
		return nil
	}

	successfull := model.Successful
	_, err = c.store.UpdateCampaignStatus(ctx, campaignID, successfull)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}

func (c *Service) ForceDeleteCampaign(ctx context.Context, campaignID int64) error {
	camapaign, err := c.store.GetCampaignByID(ctx, campaignID)
	if err != nil {
		if err == store.ErrNotFound {
			return contracts.ErrCampaignNotFound
		}
		return fmt.Errorf("get campaign by id: %w", err)
	}

	if camapaign.Status != model.Active {
		return contracts.ErrCampaignNotActive
	}

	err = c.Tx.RefundDonations(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("refund donations: %w", err)
	}

	failed := model.Failed
	_, err = c.store.UpdateCampaignStatus(ctx, campaignID, failed)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}
