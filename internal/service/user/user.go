package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/auth"
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
		if errors.Is(err, store.ErrNotFound) {
			return contracts.ErrUserNotFound
		}
		return fmt.Errorf("update user info: %w", err)
	}
	return nil
}

func (u *Service) UserInfo(ctx context.Context, id int64) (*dto.UserResponse, error) {
	userDB, err := u.Store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, contracts.ErrUserNotFound
		}
		return nil, fmt.Errorf("get by id: %w", err)
	}
	user := dto.UserResponse{
		FirstName: userDB.FirstName,
		LastName:  userDB.LastName,
		Email:     userDB.Email,
		ID:        userDB.ID,
		Balance:   userDB.Balance,
		CreatedAt: userDB.CreatedAt,
		DeletedAt: userDB.DeletedAt,
	}
	return &user, nil
}

func (u *Service) ChangePassword(ctx context.Context, input *dto.ChangeUserPassword) error {
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = u.Store.UpdatePassword(ctx, dto.UpdateUserPassword{
		ID:           input.ID,
		PasswordHash: string(hash),
	})
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
	user, err := u.Store.GetByIDForPurge(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	if user.Balance.Sign() > 0 {
		return fmt.Errorf(
			"user has remaining balance (%s): %w",
			user.Balance.String(),
			contracts.ErrDataConflict,
		)
	}

	if user.DeletedAt == nil {
		return fmt.Errorf(
			"user account is not deactivated: %w",
			contracts.ErrDataConflict,
		)
	}

	if err := u.Store.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
