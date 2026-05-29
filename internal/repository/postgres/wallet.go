package postgres

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/shopspring/decimal"
)

func (w *Repository) GetBalance(ctx context.Context, account int64) (*decimal.Decimal, error) {
	balance, err := w.DB.GetBalance(ctx, account)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (w *Repository) AddBalance(ctx context.Context, input dto.BalanceOperationInput) error {
	err := w.DB.AddBalance(ctx, mapper.BalanceParams(input))
	if err != nil {
		return err
	}
	return nil
}

func (t *Repository) SubtractBalance(ctx context.Context,
	input dto.BalanceOperationInput,
) (int64, error) {
	if input.Amount.Sign() <= 0 {
		return 0, errors.New("balance operation <= 0")
	}
	rows, err := t.DB.SubtractBalance(ctx, mapper.SubtractBalanceParams(input))
	if err != nil {
		return 0, err
	}
	return rows, nil
}
