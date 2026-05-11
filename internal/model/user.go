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
	ID           int32           `json:"id" db:"id"`
	Balance      decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt    time.Time       `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt" db:"updatedAt"`
	DeletedAt    time.Time       `json:"deletedAt" db:"deletedAt"`
}

type UserInput struct {
	FirstName    string          `json:"firstName" db:"firstName"`
	LastName     string          `json:"lastName" db:"lastName"`
	Email        string          `json:"email" db:"email"`
	HashPassword string          `json:"hashPassword" db:"hashPassword"`
	UpdatedAt    time.Time       `json:"updatedAt" db:"updatedAt"`
}

type UserResponse struct {
	FirstName    string          `json:"firstName" db:"firstName"`
	LastName     string          `json:"lastName" db:"lastName"`
	Email        string          `json:"email" db:"email"`
	ID           int32           `json:"id" db:"id"`
	Balance      decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt    time.Time       `json:"createdAt" db:"createdAt"`
	DeletedAt    time.Time       `json:"deletedAt" db:"deletedAt"`
}