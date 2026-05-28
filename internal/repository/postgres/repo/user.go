package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	Store store.SQLStore
}

func NewUserStore(store store.SQLStore) *UserStore {
	return &UserStore{
		Store: store,
	}
}

func (u *UserStore) CreateUser(
	ctx context.Context,
	input model.UserInput,
) (*model.UserResponse, error) {

	if err := u.CheckEmail(ctx, input); err != nil {
		return nil, err
	}

	hashPassword, err := hashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("technical error: %w", err)
	}
	strHashPassword := string(hashPassword)
	postUserInput := mapper.UserParams(input, strHashPassword)

	storeUser, err := u.Store.DB.CreateUser(ctx, postUserInput)
	if err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}

	user := mapper.UserResponse(storeUser)
	return &user, nil
}

func (u *UserStore) GetByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	storeUser, err := u.Store.DB.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user := mapper.UserResponse(storeUser)
	return &user, nil
}

func (u *UserStore) GetByID(ctx context.Context, id int64) (*model.UserResponse, error) {
	storeUser, err := u.Store.DB.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user := mapper.UserResponse(storeUser)
	return &user, nil
}

func (u *UserStore) UpdateUser(
	ctx context.Context,
	input model.UserInput,
) (*model.UserResponse, error) {
	if err := u.CheckEmail(ctx, input); err != nil {
		return nil, err
	}

	bytehashPassword, err := hashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("technical error: %w", err)
	}
	hashPassword := string(bytehashPassword)

	userParams := mapper.UpdateUserParams(input, hashPassword)
	postUser, err := u.Store.DB.UpdateUser(ctx, userParams)
	if err != nil {
		return nil, err
	}
	user := mapper.UserResponse(postUser)
	return &user, nil
}

func (u *UserStore) DeleteUser(ctx context.Context, id int64) error {
	err := u.Store.DB.SoftDeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) RestoreUser(ctx context.Context, id int64) (*model.UserResponse, error) {
	postUser, err := u.Store.DB.RestoreUser(ctx, id)
	if err != nil {
		return nil, err
	}
	user := mapper.UserResponse(postUser)
	return &user, nil
}

func CheckPassword(password string, hashPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) CheckEmail(ctx context.Context, input model.UserInput) error {
	user, err := u.Store.DB.GetByEmail(ctx, input.Email)
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

func hashPassword(password string) ([]byte, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return HashPassword, nil
}
