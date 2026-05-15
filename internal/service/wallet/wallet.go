package wallet

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
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

func (w *WalletService) GetBalance(ctx context.Context, account int32) (*decimal.Decimal, error) {
	balance, err := w.Store.GetBalance(ctx, account)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (w *WalletService) TopUpBalance(ctx context.Context, input model.TopUpParams) error {
	params := postgres.TopUpParams{
		ID:      input.ID,
		Balance: input.Balance,
	}
	err := w.Store.TopUp(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (w *WalletService) WithDraw(ctx context.Context, input model.TopUpParams) error {
	params := postgres.WithdrawParams{
		ID:      input.ID,
		Balance: input.Balance,
	}
	rows, err := w.Store.Withdraw(ctx, params)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("insufficient funds")
	}
	return nil
}

