package donate

import (
	"context"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
)

func CheckStatusCampaign(input postgres.CampaignStatus) bool {
	status := model.CampaignStatus(input)
	if status == model.Active {
		return false
	}
	return true
}

func CreateDonateService(store repository.Store) *DonateService {
	return &DonateService{
		Store: store,
	}
}

func (d *DonateService) createDonation(ctx context.Context, q postgres.Querier, input postgres.CreateDonationParams) error {
	donate, err := q.CreateDonation(ctx, input)
	if err != nil {
		return err
	}

	q.CreateTransaction(ctx, postgres.CreateTransactionParams{
		UserID:          input.UserID,
		DonationID:      helperfunc.Int(donate.ID),
		TransactionType: postgres.TransactionType(model.TransactionDonation),
		Amount:          input.Amount,
	})

	return nil
}
