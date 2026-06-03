package donate

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/alibek2024/FundGo/internal/service/transaction"
)

type Service struct {
	TX store.TransactionManager
}

func CreateDonateService(tx store.TransactionManager) *Service {
	return &Service{
		TX: tx,
	}
}

func (d *Service) DonateToCampaign(
	ctx context.Context,
	input dto.DonateInput,
) error {
	return d.TX.ExecTx(ctx, func(q store.Store) error {
		status, err := q.GetCampaignStatus(ctx, input.CampaignID)
		if err != nil {
			return err
		}

		if status == nil || *status != model.Active {
			return contracts.ErrCampaignNotActive
		}

		balance, err := q.GetBalance(ctx, input.UserID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrUserNotFound
			}
			return fmt.Errorf("get balance: %w", err)
		}

		subtractBalance := dto.BalanceOperationInput{
			ID:     input.UserID,
			Amount: input.Amount,
		}

		err = q.SubtractBalance(ctx, subtractBalance)
		if err != nil {
			if errors.Is(err, store.ErrDataConflict) {
				return contracts.ErrInsufficientBalance
			}
			return fmt.Errorf("subtract balance: %w", err)
		}

		increaseBalance := dto.CampaignBalanceOperation{
			ID:     input.CampaignID,
			Amount: input.Amount,
		}

		_, err = q.IncreaseCampaignAmount(ctx, increaseBalance)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrCampaignNotFound
			}
			return fmt.Errorf("increase balance: %w", err)
		}

		donation, err := q.CreateDonation(ctx, input)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrCampaignNotFound
			}
			return fmt.Errorf("create donation: %w", err)
		}

		balanceAfter := balance.Sub(input.Amount)

		txParams := transaction.ToTransactionModel(donation.UserID,
			&donation.ID, string(model.TransactionDonation), donation.Amount, *balance, balanceAfter)

		_, err = q.CreateTransaction(ctx, txParams)
		if err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		return nil
	})
}

func (c *Service) RefundDonation(ctx context.Context, donationID int64) error {
	return c.TX.ExecTx(ctx, func(q store.Store) error {
		donation, err := q.GetDonationByID(ctx, donationID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrDonationNotFound
			}
			return fmt.Errorf("get by donation: %w", err)
		}
		if donation.Status == model.DonationRefund {
			return contracts.ErrDonateRefunded
		}
		if donation.Status == model.DonationArchived {
			return contracts.ErrCampaignClosed
		}

		campaign, err := q.GetCampaignByID(ctx, donation.CampaignID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrCampaignNotFound
			}
			return fmt.Errorf("get by donation: %w", err)
		}

		if campaign.Status != model.Active {
			return contracts.ErrCampaignNotActive
		}

		user, err := q.GetByID(ctx, donation.UserID)
		if err != nil {
			if err == store.ErrNotFound {
				return contracts.ErrUserNotFound
			}
			return fmt.Errorf("get user by id: %w", err)
		}

		refundBalanceParams := dto.CampaignBalanceOperation{
			ID:     donation.CampaignID,
			Amount: donation.Amount,
		}

		_, err = q.DecreaseCampaignBalance(ctx, refundBalanceParams)
		if err != nil {
			if errors.Is(err, store.ErrDataConflict) {
				return contracts.ErrDataConflict
			}
			return fmt.Errorf("decrease campaign: %w", err)
		}

		increaseBalance := dto.BalanceOperationInput{
			ID:     donation.UserID,
			Amount: donation.Amount,
		}

		err = q.AddBalance(ctx, increaseBalance)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrUserNotFound
			}
			return err
		}

		_, err = q.RefundDonationStatus(ctx, donationID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrDonationNotFound
			}
			return fmt.Errorf("refund donation status: %w", err)
		}

		BalanceAfter := user.Balance.Add(increaseBalance.Amount)

		txParams := transaction.ToTransactionModel(
			user.ID,
			nil,
			string(model.TransactionRefund),
			donation.Amount,
			user.Balance,
			BalanceAfter,
		)

		_, err = q.CreateTransaction(ctx, txParams)
		if err != nil {
			return contracts.ErrDataConflict
		}

		return nil
	})
}
