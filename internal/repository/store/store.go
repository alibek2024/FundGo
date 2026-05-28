package store

import (
	"context"

	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	UserStore
	CampaignStore
	DonationStore
	TrasactionStore
	ExecTx(ctx context.Context, fn func(Store) error) error
}

type SQLStore struct {
	DB   *generated.Queries
	Conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		DB:   generated.New(conn),
		Conn: conn,
	}
}
