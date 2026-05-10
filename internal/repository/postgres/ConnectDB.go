package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASEURL")

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	log.Println("Успешное подключение к PostgreSQL!")
	return dbPool, nil
}
