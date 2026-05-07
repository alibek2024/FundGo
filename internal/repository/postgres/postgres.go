package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct{
	DB sqlx.DB
}

func ConnectPostgres() (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := "host=localhost port=5432 user=postgres password=secret dbname=fund_go sslmode=disable"
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
