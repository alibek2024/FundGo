package wallet

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/shopspring/decimal"
)

type WalletService struct {
	Store repository.Store
}

func NewWalletService(store repository.Store) *WalletService {
	return &WalletService{
		Store: store,
	}
}

func (w *WalletService) GetBalance(ctx context.Context, account int64) (*decimal.Decimal, error) {
	balance, err := w.Store.GetBalance(ctx, account)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (w *WalletService) TopUpBalance(ctx context.Context, input model.Balance) error {
	return w.Store.ExecTx(ctx, func(q postgres.Querier) error {
		params := postgres.AddBalanceParams{
			ID:      input.ID,
			Balance: input.Amount,
		}
		err := q.AddBalance(ctx, params)
		if err != nil {
			return err
		}

		q.CreateTransaction(ctx, postgres.CreateTransactionParams{
			UserID:          helperfunc.Int(input.ID),
			TransactionType: postgres.TransactionType(model.TransactionTopUp),
			Amount:          input.Amount,
		})

		return nil
	})
}

func (w *WalletService) WithDraw(ctx context.Context, input model.Balance) error {
	return w.Store.ExecTx(ctx, func(q postgres.Querier) error {
		rows, err := q.SubtractBalance(ctx, postgres.SubtractBalanceParams{
			ID:      input.ID,
			Balance: input.Amount,
		})
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("insufficient funds")
		}

		_, err = q.CreateTransaction(ctx, postgres.CreateTransactionParams{
			UserID: helperfunc.Int(input.ID),
			Amount: input.Amount,
			TransactionType: postgres.TransactionType(
				model.TransactionWithdraw),
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (w *WalletService) DebitUserBalance(ctx context.Context, input model.Balance) error {
	rows, err := w.Store.SubtractBalance(ctx, postgres.SubtractBalanceParams{
		ID:      input.ID,
		Balance: input.Amount,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("insufficient funds")
	}
	return nil
}
