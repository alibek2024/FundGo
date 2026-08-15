package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/hibiken/asynq"
)

type CampaignTaskHandler struct {
	campaignService campaign.Service
}

func NewCampaignTaskHandler(campaignService campaign.Service) CampaignTaskHandler {
	return CampaignTaskHandler{campaignService: campaignService}
}

func (h *CampaignTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p campaign.CloseCampaignPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	if err := h.campaignService.WrapUpCampaign(ctx, p.CampaignID); err != nil {
		return fmt.Errorf("failed to close campaign %d: %w", p.CampaignID, err)
	}
	return nil
}