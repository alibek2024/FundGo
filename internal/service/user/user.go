package user

import (
	"context"
	"fmt"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
)

type UserService struct {
	Store repository.Store
}

func NewUserService(store repository.Store) *UserService {
	return &UserService{
		Store: store,
	}
}

func (u *UserService) CreateUser(
	ctx context.Context,
	input model.UserInput,
) (*model.UserResponse, error) {

	if err := u.CheckEmail(ctx, input); err != nil {
		return nil, err
	}

	hashPassword, err := u.hashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("technical error: %w", err)
	}
	strHashPassword := string(hashPassword)
	postUserInput := u.userParams(input, strHashPassword)

	storeUser, err := u.Store.CreateUser(ctx, postUserInput)
	if err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}

	user := u.UserResponse(storeUser)
	return &user, nil
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	storeUser, err := u.Store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user := u.UserResponse(storeUser)
	return &user, nil
}

func (u *UserService) GetByID(ctx context.Context, id int32) (*model.UserResponse, error) {
	storeUser, err := u.Store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user := u.UserResponse(storeUser)
	return &user, nil
}

func (u *UserService) UpdateUser(
	ctx context.Context,
	input model.UserInput,
) (*model.UserResponse, error) {
	if err := u.CheckEmail(ctx, input); err != nil {
		return nil, err
	}

	bytehashPassword, err := u.hashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("technical error: %w", err)
	}
	hashPassword := string(bytehashPassword)

	userParams := u.updateUserParams(input, hashPassword)
	postUser, err := u.Store.UpdateUser(ctx, userParams)
	if err != nil {
	  return nil, err
	}
	user := u.UserResponse(postUser)
	return &user, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int32) error {
	err := u.Store.SoftDeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) RestoreUser(ctx context.Context, id int32) (*model.UserResponse, error) {
	postUser, err := u.Store.RestoreUser(ctx, id)
	if err != nil {
	  return nil, err
	}
	user := u.UserResponse(postUser)
	return &user, nil
}

