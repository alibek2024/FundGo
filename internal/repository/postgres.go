package repository

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func ConnectDB() (*sqlx.DB, error) {
  	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	dsn := "host=localhost port=5432 user=postgres password=secret dbname=fund_go sslmode=disable"
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil 
}