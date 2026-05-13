package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) defaultParams(
	input model.UserInput,
	hashPassword string,
) (
	string, string,
	pgtype.Text, pgtype.Text,
) {
	return input.Email,
		hashPassword,
		pgtype.Text{String: input.FirstName, Valid: input.FirstName != ""},
		pgtype.Text{String: input.LastName, Valid: input.FirstName != ""}
}

func (u *UserService) userParams(input model.UserInput, hashPassword string) postgres.CreateUserParams {
	email, hash, first, last := u.defaultParams(input, hashPassword)
	return postgres.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		FirstName:    first,
		LastName:     last,
	}
}

func (u *UserService) updateUserParams(input model.UserInput, hashPassword string) postgres.UpdateUserParams {
	email, hash, first, last := u.defaultParams(input, hashPassword)
	return postgres.UpdateUserParams{
		Email:        email,
		PasswordHash: hash,
		FirstName:    first,
		LastName:     last,
		ID:           input.ID,
	}
}

func (u *UserService) UserResponse(input postgres.User) model.UserResponse {
	return model.UserResponse{
		FirstName: input.FirstName.String,
		LastName:  input.LastName.String,
		ID:        input.ID,
		Balance:   input.Balance,
		CreatedAt: input.CreatedAt.Time,
		DeletedAt: input.DeletedAt.Time,
	}
}

func (u *UserService) hashPassword(password string) ([]byte, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return HashPassword, nil
}

func (u *UserService) CheckPassword(password string, hashPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) CheckEmail(ctx context.Context, input model.UserInput) error {
	user, err := u.Store.GetByEmail(ctx, input.Email)
	if err == nil {
		if input.ID != user.ID {
			return fmt.Errorf("This %s email is taken", input.Email)
		}
		return nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("technical error: %w", err)
	}
	return nil
}


