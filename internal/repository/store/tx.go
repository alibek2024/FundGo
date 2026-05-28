package store

import (
	"context"
)

func (s SQLStore) ExecTx(ctx context.Context, fn func(st SQLStore) error) error {
	tx, err := s.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txStore := &SQLStore{
		DB:   s.DB.WithTx(tx),
		Conn: s.Conn,
	}
	if err := fn(*txStore); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
