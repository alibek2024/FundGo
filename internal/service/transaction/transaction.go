package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/shopspring/decimal"
)

type Service struct {
	Store store.TrasactionStore
}

func CreateTX(db store.Store) *Service {
	return &Service{
		Store: db,
	}
}

func (t *Service) GetPaymentHistory(ctx context.Context,
	userID int64,
) ([]model.Transaction, error) {
	txHistory, err := t.Store.HistoryTX(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, contracts.ErrUserNotFound
		}
		return nil, fmt.Errorf("history tx: %w", err)
	}

	return txHistory, nil
}

func ToTransactionModel(
	userID int64,
	donationID *int64,
	txType string,
	amount, balanceBefore, balanceAfter decimal.Decimal,
) dto.TransactionInput {
	if donationID == nil {
		return dto.TransactionInput{
			UserID:        userID,
			DonationID:    nil,
			Type:          txType,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
		}
	}
	return dto.TransactionInput{
		UserID:        userID,
		DonationID:    donationID,
		Type:          txType,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
	}
}
