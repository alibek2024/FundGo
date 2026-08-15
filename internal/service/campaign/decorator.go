package campaign

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/hibiken/asynq"
)

type CampaignScheduler struct {
	service Service
	asyncClient *asynq.Client
}

const TypeCloseCampaign = "campaign:close"

type CloseCampaignPayload struct {
	CampaignID int64 `json:"campaign_id"`
}

func NewCampaignScheduler(
	service Service, 
	asynqClient *asynq.Client,
) *CampaignScheduler {
	return &CampaignScheduler{
		service:        service,
		asyncClient: asynqClient,
	}
}

func (c *CampaignScheduler) CreateCampaign(
	ctx context.Context, 
	input dto.CreateCampaignInput,
) (*model.Campaign, error) {
	campaign, err := c.service.CreateCampaign(ctx, input)
	if err != nil {
	  return nil, fmt.Errorf("create camapaign :%w", err)
	}

	payload, err := json.Marshal(CloseCampaignPayload{
		CampaignID: campaign.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeCloseCampaign, payload)
	_, err = c.asyncClient.EnqueueContext(
		ctx, 
		task, 
		asynq.ProcessAt(campaign.EndDate),
		asynq.TaskID(fmt.Sprintf("close-campaign-%d", campaign.ID)),
	)
	if err != nil {
		fmt.Printf("failed to enqueue close task for campaign %d: %v\n", campaign.ID, err)
	}
	return campaign, nil
}

func (c *CampaignScheduler) ForceCloseCampaign(ctx context.Context, id int64) error {
	return c.service.ForceDeleteCampaign(ctx, id)
}
func (c *CampaignScheduler) CloseCampaign(ctx context.Context, id int64) error {
	return c.service.WrapUpCampaign(ctx, id)
}