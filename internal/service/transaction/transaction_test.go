package transaction_test

import (
	"context"
	"testing"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/alibek2024/FundGo/internal/service/transaction"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPaymentHistory(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*store.MockStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockStore) {
				mockStore.EXPECT().
					HistoryTX(mock.Anything, int64(1)).
					Return([]model.Transaction{{ID: 1, UserID: 1}}, nil)
			},
		},
		{
			name: "user not found",
			setup: func(mockStore *store.MockStore) {
				mockStore.EXPECT().
					HistoryTX(mock.Anything, int64(1)).
					Return(nil, store.ErrNotFound)
			},
			wantErr: contracts.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			service := transaction.CreateTX(mockStore)
			got, err := service.GetPaymentHistory(context.Background(), 1)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got)
		})
	}
}
