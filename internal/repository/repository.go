package repository

import (
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLStore struct {
	*postgres.Queries
	Conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		Queries: postgres.New(conn),
		Conn: conn,
	}
}