package mapper

import (
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return store.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return store.ErrAlreadyExists
		case "23503":
			return store.ErrDataConflict
		case "23502":
			return store.ErrInternalDB
		}
		return fmt.Errorf("%w: pg_code=%s", store.ErrInternalDB, pgErr.Code)
	}

	return fmt.Errorf("%w: %v", store.ErrInternalDB, err)
}
