package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserResponse struct {
	FirstName string          `json:"firstName" validate:"required,min=3,max=50"`
	LastName  string          `json:"lastName" validate:"required,min=3,max=50"`
	Email     string          `json:"email" validate:"required,email"`
	ID        int64           `json:"id" validate:"gt=0"`
	Balance   decimal.Decimal `json:"balance" validate:"required"`
	CreatedAt time.Time       `json:"createdAt"`
	DeletedAt *time.Time      `json:"deletedAt"`
}

type UserInfo struct {
	FirstName string `json:"firstName" validate:"required,min=3,max=50"`
	LastName  string `json:"lastName" validate:"required,min=3,max=50"`
	ID        int64  `json:"-" validate:"required,gt=0"`
}

type UserEmail struct {
	Email string `json:"email" validate:"required,email"`
	ID    int64  `json:"-" validate:"required,gt=0"`
}

type ChangeUserPassword struct {
	HashPassword string `json:"password" validate:"required,min=6"`
	ID           int64  `json:"-" validate:"required,gt=0"`
}
