package wallet

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
	Tx store.TransactionManager
}

func NewWalletService(store store.TransactionManager) *Service {
	return &Service{
		Tx: store,
	}
}

func (w *Service) TopUpBalance(ctx context.Context, input dto.BalanceOperationInput) error {
	return w.Tx.ExecTx(ctx, func(q store.Store) error {
		user, err := q.GetByID(ctx, input.ID)
		if err != nil {
			if err == store.ErrNotFound {
				return contracts.ErrUserNotFound
			}
			return fmt.Errorf("get user by id: %w", err)
		}

		err = q.AddBalance(ctx, input)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrUserNotFound
			}
			return err
		}

		BalanceAfter := user.Balance.Add(input.Amount)

		params := transaction.ToTransactionModel(
			input.ID, nil,
			string(model.TransactionTopUp),
			input.Amount, user.Balance, BalanceAfter,
		)

		_, err = q.CreateTransaction(ctx, params)
		if err != nil {
			return contracts.ErrDataConflict
		}

		return nil
	})
}

func (w *Service) WithdrawBalance(ctx context.Context, input dto.BalanceOperationInput) error {
	return w.Tx.ExecTx(ctx, func(q store.Store) error {
		balance, err := q.GetBalance(ctx, input.ID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return contracts.ErrUserNotFound
			}
			return fmt.Errorf("get balance: %w", err)
		}

		subtractBalance := dto.BalanceOperationInput{
			ID:     input.ID,
			Amount: input.Amount,
		}

		err = q.SubtractBalance(ctx, subtractBalance)
		if err != nil {
			if errors.Is(err, store.ErrDataConflict) {
				return contracts.ErrInsufficientBalance
			}
			return fmt.Errorf("subtract balance: %w", err)
		}

		balanceAfter := balance.Sub(input.Amount)

		txParams := transaction.ToTransactionModel(input.ID,
			nil, string(model.TransactionWithdraw), input.Amount, *balance, balanceAfter)

		_, err = q.CreateTransaction(ctx, txParams)
		if err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		return nil
	})
}
