package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Store store.UserStore
}

func NewUserService(store store.UserStore) *Service {
	return &Service{
		Store: store,
	}
}

func (u Service) RegisterUser(ctx context.Context, input *dto.RegistrationInput) (*model.User, error) {
	if err := u.CheckEmail(ctx, input.Email); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return nil, service.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("check email: %w", err)
	}
	hashPassword, err := hashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	input.HashPassword = string(hashPassword)

	user, err := u.Store.CreateUser(ctx, *input)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return nil, service.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (u Service) SignIn(ctx context.Context, input dto.SignIn) (string, error) {
	user, err := u.Store.GetByID(ctx, input.ID)
	if err != nil {
	  if err == store.ErrNotFound {
		return "", service.ErrUserNotFound
	  }
	  return "", fmt.Errorf("get user by id: %w", err)
	}

	
}

func (u Service) UpdateUserInfo(ctx context.Context, input *dto.UserInfo) error {
	_, err := u.Store.UpdateInfo(ctx, *input)
	if err != nil {
		return fmt.Errorf("update user info: %w", err)
	}
	return nil
}

func (u Service) ChangePassword(ctx context.Context, input *dto.ChangeUserPassword) error {
	_, err := u.Store.UpdatePassword(ctx, *input)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	return nil
}

func (u Service) ChangeEmail(ctx context.Context, input *dto.UserEmail) error {
	_, err := u.Store.UpdateEmail(ctx, *input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return service.ErrEmailAlreadyExists
		}
		return fmt.Errorf("check email: %w", err)
	}
	return nil
}

func (u Service) DeactivateAccount(ctx context.Context, userID int64) error {
	user, err := u.Store.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user.Balance.Sign() > 0 {
		return fmt.Errorf("user has remaining balance (%d): %w", user.Balance, service.ErrDataConflict)
	}

	err = u.Store.SoftDeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	return nil
}
func (u Service) PurgeUserData(ctx context.Context, userID int64) error {
	user, err := u.Store.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user.Balance.Sign() > 0 {
		return fmt.Errorf("user has remaining balance (%d): %w", user.Balance, service.ErrDataConflict)
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

func (u *Service) CheckEmail(ctx context.Context, Email string) error {
	_, err := u.Store.GetByEmail(ctx, Email)
	if err != nil {
		return fmt.Errorf("check email: %w", err)
	}
	if err == nil {
		return service.ErrEmailAlreadyExists
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
