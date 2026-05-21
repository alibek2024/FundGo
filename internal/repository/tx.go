package repository

import (
	"context"

	"github.com/alibek2024/FundGo/internal/repository/postgres"
)

func (s SQLStore) ExecTx(ctx context.Context, fn func(postgres.Querier) error) error {
	tx, err := s.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := s.Queries.WithTx(tx)
	if err := fn(q); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
