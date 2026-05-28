package repo

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/alibek2024/FundGo/internal/repository/store"
)

var ErrCampaignInactive = errors.New("campaign inactive")
var ErrInsufficientFunds = errors.New("insufficient funds")

type DonationStore struct {
	Store store.SQLStore
}

func CreateDonateService(store store.SQLStore) *DonationStore {
	return &DonationStore{
		Store: store,
	}
}

func (d *DonationStore) DonateToCampaign(ctx context.Context, input model.DonateInput) error {
	return d.Store.ExecTx(ctx, func(q store.SQLStore) error {
		status, err := q.DB.GetCampaignStatus(ctx, input.CampaignID)
		if err != nil {
			return err
		}
		if !mapper.CheckStatusCampaign(status) {
			return ErrCampaignInactive
		}

		rows, err := q.DB.SubtractBalance(ctx, generated.SubtractBalanceParams{
			ID:      input.UserID,
			Balance: input.Amount,
		})
		if err != nil {
			return err
		}
		if rows == 0 {
			return ErrInsufficientFunds
		}

		_, err = q.DB.IncreaseCampaignAmount(ctx, generated.IncreaseCampaignAmountParams{
			ID:            input.CampaignID,
			CurrentAmount: input.Amount,
		})
		if err != nil {
			return err
		}

		err = d.CreateDonation(ctx, q.DB, generated.CreateDonationParams{
			UserID:     mapper.Int(input.UserID),
			CampaignID: input.CampaignID,
			Amount:     input.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *DonationStore) CreateDonation(ctx context.Context, q generated.Querier, input generated.CreateDonationParams) error {
	donate, err := q.CreateDonation(ctx, input)
	if err != nil {
		return err
	}
	balance, err := q.GetBalance(ctx, input.UserID.Int64)
	if err != nil {
		return err
	}

	q.CreateTransaction(ctx, generated.CreateTransactionParams{
		UserID:          input.UserID,
		DonationID:      mapper.Int(donate.ID),
		TransactionType: generated.TransactionType(model.TransactionDonation),
		Amount:          input.Amount,
		BalanceBefore:   balance,
		BalanceAfter:    balance.Sub(input.Amount),
	})

	return nil
}
