package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type RegistrationInput struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	HashPassword string `json:"hashPassword"`
}

type LoginInput struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	HashPassword string `json:"hashPassword"`
}

type UserResponse struct {
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Email     string          `json:"email"`
	ID        int64           `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"createdAt"`
	DeletedAt time.Time       `json:"deletedAt"`
}

type UserInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ID        int64  `json:"id"`
}

type UserEmail struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
}

type ChangeUserPassword struct {
	HashPassword string `json:"hashPassword"`
	ID           int64  `json:"id"`
}
