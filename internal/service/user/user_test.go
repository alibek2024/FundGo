package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/alibek2024/FundGo/internal/service/user"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserInfo(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*store.MockUserStore)
		wantErr string
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdateInfo(mock.Anything, dto.UserInfo{
						ID:        1,
						FirstName: "NewFirstName",
						LastName:  "NewLastName",
					}).
					Return(&model.User{ID: 1}, nil)
			},
		},
		{
			name: "store error",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdateInfo(mock.Anything, mock.AnythingOfType("dto.UserInfo")).
					Return(nil, errors.New("db unavailable"))
			},
			wantErr: "update user info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			service := user.NewUserService(mockStore)
			err := service.UpdateUserInfo(context.Background(), &dto.UserInfo{
				ID:        1,
				FirstName: "NewFirstName",
				LastName:  "NewLastName",
			})

			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestChangeEmail(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*store.MockUserStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdateEmail(mock.Anything, dto.UserEmail{ID: 1, Email: "new@example.com"}).
					Return(&model.User{ID: 1, Email: "new@example.com"}, nil)
			},
		},
		{
			name: "email already exists",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdateEmail(mock.Anything, mock.AnythingOfType("dto.UserEmail")).
					Return(nil, store.ErrNotFound)
			},
			wantErr: contracts.ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			service := user.NewUserService(mockStore)
			err := service.ChangeEmail(context.Background(), &dto.UserEmail{
				ID:    1,
				Email: "new@example.com",
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDeactivateAccount(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*store.MockUserStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(&model.User{ID: 1, Balance: decimal.Zero}, nil)
				mockStore.EXPECT().
					SoftDeleteUser(mock.Anything, int64(1)).
					Return(nil)
			},
		},
		{
			name: "remaining balance",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(&model.User{ID: 1, Balance: decimal.NewFromInt(10)}, nil)
			},
			wantErr: contracts.ErrDataConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			service := user.NewUserService(mockStore)
			err := service.DeactivateAccount(context.Background(), 1)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPurgeUserData(t *testing.T) {
	deletedAt := time.Now()

	tests := []struct {
		name    string
		setup   func(*store.MockUserStore)
		wantErr string
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(&model.User{ID: 1, Balance: decimal.Zero, DeletedAt: &deletedAt}, nil)
				mockStore.EXPECT().
					DeleteUser(mock.Anything, int64(1)).
					Return(nil)
			},
		},
		{
			name: "account is not deleted",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(&model.User{ID: 1, Balance: decimal.Zero}, nil)
			},
			wantErr: "user did not delete his account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			service := user.NewUserService(mockStore)
			err := service.PurgeUserData(context.Background(), 1)

			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
