package wallet_test

import (
	"context"
	"testing"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/alibek2024/FundGo/internal/service/wallet"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTopUpBalance(t *testing.T) {
	input := dto.BalanceOperationInput{
		ID:     1,
		Amount: decimal.NewFromInt(25),
	}

	tests := []struct {
		name    string
		setup   func(*store.MockStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(&model.User{ID: 1, Balance: decimal.NewFromInt(10)}, nil)
				mockStore.EXPECT().
					AddBalance(mock.Anything, input).
					Return(nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
		},
		{
			name: "user not found",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(nil, store.ErrNotFound)
			},
			wantErr: contracts.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			service := wallet.NewWalletService(mockStore)
			err := service.TopUpBalance(context.Background(), input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestWithdrawBalance(t *testing.T) {
	input := dto.BalanceOperationInput{
		ID:     1,
		Amount: decimal.NewFromInt(25),
	}
	balance := decimal.NewFromInt(100)

	tests := []struct {
		name    string
		setup   func(*store.MockStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(1)).
					Return(&balance, nil)
				mockStore.EXPECT().
					SubtractBalance(mock.Anything, input).
					Return(nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
		},
		{
			name: "insufficient balance",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(1)).
					Return(&balance, nil)
				mockStore.EXPECT().
					SubtractBalance(mock.Anything, input).
					Return(store.ErrDataConflict)
			},
			wantErr: contracts.ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			service := wallet.NewWalletService(mockStore)
			err := service.WithdrawBalance(context.Background(), input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func expectExecTx(mockStore *store.MockStore) {
	mockStore.EXPECT().
		ExecTx(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(store.Store) error) error {
			return fn(mockStore)
		})
}
