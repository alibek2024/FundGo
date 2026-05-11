package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store repository.SQLStore
}

func (u UserService) CreateUser(
	ctx context.Context,
	input model.UserInput,
) (*model.UserResponse, error) {
	_, err := u.Store.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, fmt.Errorf("This %s email is taken", input.Email)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("technical error: %w", err)
	}

	hashPassword, err := HashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("technical error: %w", err)
	}
	strHashPassword := string(hashPassword)
	postUserInput := UserParams(input, strHashPassword)

	storeUser, err := u.Store.CreateUser(ctx, postUserInput)
	if err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}

	user := UserResponse(storeUser)
	return &user, nil
}

func HashPassword(password string) ([]byte, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return []byte(password), err
	}
	return HashPassword, nil
}

func CheckPassword(password string, hashPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}
