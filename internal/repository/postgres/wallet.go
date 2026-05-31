package postgres

import (
	"context"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/shopspring/decimal"
)

func (w *Repository) GetBalance(ctx context.Context, account int64) (*decimal.Decimal, error) {
	balance, err := w.DB.GetBalance(ctx, account)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	return &balance, nil
}

func (w *Repository) AddBalance(ctx context.Context, input dto.BalanceOperationInput) error {
	err := w.DB.AddBalance(ctx, mapper.BalanceParams(input))
	if err != nil {
		return mapper.MapDBError(err)
	}
	return nil
}

func (t *Repository) SubtractBalance(ctx context.Context,
	input dto.BalanceOperationInput,
) (int64, error) {
	rows, err := t.DB.SubtractBalance(ctx, mapper.SubtractBalanceParams(input))
	if err != nil {
		return 0, mapper.MapDBError(err)
	}
	return rows, nil
}
