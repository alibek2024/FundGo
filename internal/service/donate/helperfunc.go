package donate

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
)

func (d *DonateService) toParams(input model.DonateInput) postgres.CreateDonationParams {
	return postgres.CreateDonationParams{
		UserID:     helperfunc.Int(input.UserID),
		CampaignID: input.CampaignID,
		Amount:     input.Amount,
	}
}

func (d *DonateService) toInput(p postgres.CreateDonationParams) model.DonateInput {
	return model.DonateInput{
		UserID:     p.UserID.Int32,
		CampaignID: p.CampaignID,
		Amount:     p.Amount,
	}
}

func (d *DonateService) toDonate(p postgres.Donation) model.Donation {
	return model.Donation{
		ID:         p.ID,
		UserID:     p.UserID.Int32,
		CampaignID: p.CampaignID,
		Amount:     p.Amount,
		CreatedAt:  p.CreatedAt.Time,
	}
}

func (d *DonateService) toDonationSlice(p []postgres.Donation) []*model.Donation {
	result := make([]*model.Donation, len(p))
	for i, v := range p {
		result[i] = &model.Donation{
			ID:         v.ID,
			UserID:     v.UserID.Int32,
			CampaignID: v.CampaignID,
			Amount:     v.Amount,
			CreatedAt:  v.CreatedAt.Time,
		}
	}
	return result
}
