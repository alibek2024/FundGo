package postgres

import (
	"context"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/alibek2024/FundGo/internal/repository/store"
)

func (t *Repository) HistoryTX(
	ctx context.Context,
	userID int64,
) ([]model.Transaction, error) {
	params := mapper.Int(userID)
	postgresArgs, err := t.DB.HistoryTX(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	transaction := mapper.ToTransactionModels(postgresArgs)
	return transaction, nil
}

func (t *Repository) CreateTransaction(
	ctx context.Context,
	input dto.TransactionInput,
) (*model.Transaction, error) {
	params := mapper.ToTXPostgresParams(input)
	postgTx, err := t.DB.CreateTransaction(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	transaction := mapper.ToTransactionModel(postgTx)
	return transaction, nil
}

func (s Repository) ExecTx(ctx context.Context, fn func(store.Store) error) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return mapper.MapDBError(err)
	}
	defer tx.Rollback(ctx)
	txRepo := Repository{
		DB:   s.DB.WithTx(tx),
		conn: s.conn,
	}
	if err := fn(&txRepo); err != nil {
		return mapper.MapDBError(err)
	}
	return tx.Commit(ctx)
}
