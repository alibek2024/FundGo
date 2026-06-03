package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
)

type Service struct {
	Store store.UserStore
}

func NewUserService(store store.UserStore) *Service {
	return &Service{
		Store: store,
	}
}

func (u *Service) UpdateUserInfo(ctx context.Context, input *dto.UserInfo) error {
	_, err := u.Store.UpdateInfo(ctx, *input)
	if err != nil {
		return fmt.Errorf("update user info: %w", err)
	}
	return nil
}

func (u *Service) ChangePassword(ctx context.Context, input *dto.ChangeUserPassword) error {
	_, err := u.Store.UpdatePassword(ctx, *input)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	return nil
}

func (u *Service) ChangeEmail(ctx context.Context, input *dto.UserEmail) error {
	_, err := u.Store.UpdateEmail(ctx, *input)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return contracts.ErrEmailAlreadyExists
		}
		return fmt.Errorf("check email: %w", err)
	}
	return nil
}

func (u *Service) DeactivateAccount(ctx context.Context, userID int64) error {
	user, err := u.Store.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user.Balance.Sign() > 0 {
		return fmt.Errorf("user has remaining balance (%d): %w", user.Balance, contracts.ErrDataConflict)
	}

	err = u.Store.SoftDeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	return nil
}

func (u *Service) PurgeUserData(ctx context.Context, userID int64) error {
	user, err := u.Store.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user.Balance.Sign() > 0 {
		return fmt.Errorf("user has remaining balance (%d): %w", user.Balance, contracts.ErrDataConflict)
	}

	if user.DeletedAt == nil {
		return fmt.Errorf("user did not delete his account")
	}

	err = u.Store.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
