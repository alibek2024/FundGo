package campaign

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

type RefundManager struct {
	store.DonationStore
	Tx store.TransactionManager
}

func NewRefundManager(tx store.TransactionManager, store store.DonationStore) *RefundManager {
	return &RefundManager{
		DonationStore: store,
		Tx:            tx,
	}
}

func (c *RefundManager) RefundDonations(ctx context.Context, campaignID int64) error {
	return c.Tx.ExecTx(ctx, func(q store.Store) error {
		donors, err := q.GetListDonations(ctx, campaignID)
		if err != nil {
			return fmt.Errorf("get list donations: %w", err)
		}

		for _, v := range donors {
			user, err := q.GetByID(ctx, v.UserID)
			if err != nil {
				if err == store.ErrNotFound {
					return contracts.ErrUserNotFound
				}
				return fmt.Errorf("get user by id: %w", err)
			}

			if v.Status == model.DonationRefund {
				continue
			}

			refundBalanceParams := dto.CampaignBalanceOperation{
				ID:     v.CampaignID,
				Amount: v.Amount,
			}

			_, err = q.DecreaseCampaignBalance(ctx, refundBalanceParams)
			if err != nil {
				if errors.Is(err, store.ErrDataConflict) {
					return contracts.ErrDataConflict
				}
				return fmt.Errorf("decrease campaign: %w", err)
			}

			increaseBalance := dto.BalanceOperationInput{
				ID:     v.UserID,
				Amount: v.Amount,
			}

			err = q.AddBalance(ctx, increaseBalance)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return contracts.ErrUserNotFound
				}
				return fmt.Errorf("add balance: %w", err)
			}

			_, err = q.RefundDonationStatus(ctx, v.ID)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return contracts.ErrDonationNotFound
				}
				return fmt.Errorf("refund donation status: %w", err)
			}

			BalanceAfter := user.Balance.Add(increaseBalance.Amount)

			txParams := transaction.ToTransactionModel(
				v.UserID,
				nil,
				string(model.TransactionRefund),
				v.Amount,
				user.Balance,
				BalanceAfter,
			)

			_, err = q.CreateTransaction(ctx, txParams)
			if err != nil {
				return contracts.ErrDataConflict
			}
		}
		return nil
	})
}

func (c *RefundManager) GetCampaignDonors(ctx context.Context, campaignID int64) ([]model.Donation, error) {
	listUsers, err := c.DonationStore.GetListDonations(ctx, campaignID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, contracts.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("get list Donation: %w", err)
	}
	return listUsers, nil
}
