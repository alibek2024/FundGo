package mapper

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func CheckStatusCampaign(input generated.CampaignStatus) bool {
	status := model.CampaignStatus(input)
	if status == model.Active {
		return false
	}
	return true
}

func IncreaseCampaignAmount(input model.DonateInput) generated.IncreaseCampaignAmountParams {
	return generated.IncreaseCampaignAmountParams{
		ID: input.CampaignID,
		CurrentAmount: input.Amount,
	}
}