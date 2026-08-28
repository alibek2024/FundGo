package postgres

import (
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB   *generated.Queries
	conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *Repository {
	return &Repository{
		DB:   generated.New(conn),
		conn: conn,
	}
}
