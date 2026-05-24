package donate

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
)

var ErrCampaignInactive = errors.New("campaign inactive")
var ErrInsufficientFunds = errors.New("insufficient funds")

type DonateService struct {
	Store repository.Store
}

func (d *DonateService) DonateToCampaign(ctx context.Context, input model.DonateInput) error {
	return d.Store.ExecTx(ctx, func(q postgres.Querier) error {
		status, err := q.GetCampaignStatus(ctx, input.CampaignID)
		if err != nil {
			return err
		}
		if !CheckStatusCampaign(status) {
			return ErrCampaignInactive
		}

		rows, err := q.SubtractBalance(ctx, postgres.SubtractBalanceParams{
			ID:      input.UserID,
			Balance: input.Amount,
		})
		if err != nil {
			return err
		}
		if rows == 0 {
			return ErrInsufficientFunds
		}

		_, err = q.IncreaseCampaignAmount(ctx, postgres.IncreaseCampaignAmountParams{
			ID:            input.CampaignID,
			CurrentAmount: input.Amount,
		})
		if err != nil {
			return err
		}

		d.createDonation(ctx, q, postgres.CreateDonationParams{
			UserID:     helperfunc.Int(input.UserID),
			CampaignID: input.CampaignID,
			Amount:     input.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
