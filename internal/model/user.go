package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	FirstName    string          `json:"firstName" db:"firstName"`
	LastName     string          `json:"lastName" db:"lastName"`
	Email        string          `json:"email" db:"email"`
	HashPassword string          `json:"hashPassword" db:"hashPassword"`
	ID           int64           `json:"id" db:"id"`
	Balance      decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt    time.Time       `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt" db:"updatedAt"`
	DeletedAt    time.Time       `json:"deletedAt" db:"deletedAt"`
}
