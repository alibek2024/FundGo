package repository

import (
	"context"

	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	postgres.Querier
	ExecTx(ctx context.Context, fn func(postgres.Querier) error) error
}

type SQLStore struct {
	*postgres.Queries
	Conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		Queries: postgres.New(conn),
		Conn:    conn,
	}
}
