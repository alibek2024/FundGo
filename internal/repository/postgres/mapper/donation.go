package mapper

import (
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func ToDonationsModels(input []generated.Donation) []model.Donation {
	var result = make([]model.Donation, len(input))
	for i, v := range input {
		result[i] = model.Donation{
			ID:         v.ID,
			UserID:     v.UserID.Int64,
			CampaignID: v.CampaignID,
			Amount:     v.Amount,
			CreatedAt:  v.CreatedAt.Time,
		}
	}
	return result
}

func ToSqlcModel(input dto.DonateInput) generated.CreateDonationParams {
	return generated.CreateDonationParams{
		UserID: Int(input.UserID),
		CampaignID: input.CampaignID,
		Amount: input.Amount,
	}
}

func ToDonationModel(input generated.Donation) model.Donation {
	return model.Donation{
		ID:         input.ID,
		UserID:     input.UserID.Int64,
		CampaignID: input.CampaignID,
		Amount:     input.Amount,
		CreatedAt:  input.CreatedAt.Time,
	}
}
