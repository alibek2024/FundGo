package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alibek2024/FundGo/internal/delivery/handlers"
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/gorilla/schema"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWalletHandlerTopUpBalance(t *testing.T) {
	amount := decimal.NewFromInt(25)

	tests := []struct {
		name       string
		userID     string
		body       dto.BalanceOperationInput
		setup      func(*store.MockStore)
		wantStatus int
	}{
		{
			name:   "success uses authenticated user id",
			userID: "12",
			body: dto.BalanceOperationInput{
				ID:     999,
				Amount: amount,
			},
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(12)).
					Return(&model.User{ID: 12, Balance: decimal.NewFromInt(100)}, nil)
				mockStore.EXPECT().
					AddBalance(mock.Anything, dto.BalanceOperationInput{ID: 12, Amount: amount}).
					Return(nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing user id in context",
			body:       dto.BalanceOperationInput{ID: 12, Amount: amount},
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "validation error",
			userID: "12",
			body: dto.BalanceOperationInput{
				Amount: decimal.Zero,
			},
			setup: func(mockStore *store.MockStore) {
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "validation error negative amount",
			userID: "12",
			body: dto.BalanceOperationInput{
				Amount: decimal.NewFromInt(-25),
			},
			setup: func(mockStore *store.MockStore) {
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			handler := handlers.NewWalletHandler(newWalletService(mockStore), schema.NewDecoder())

			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/me/balance/top-up", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.userID != "" {
				req = authedRequest(req, tt.userID)
			}
			rec := httptest.NewRecorder()

			handler.TopUpBalance(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestWalletHandlerWithdrawBalance(t *testing.T) {
	balance := decimal.NewFromInt(100)
	amount := decimal.NewFromInt(25)

	tests := []struct {
		name       string
		userID     string
		body       dto.BalanceOperationInput
		setup      func(*store.MockStore)
		wantStatus int
	}{
		{
			name:   "success uses authenticated user id",
			userID: "12",
			body: dto.BalanceOperationInput{
				ID:     999,
				Amount: amount,
			},
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(12)).
					Return(&balance, nil)
				mockStore.EXPECT().
					SubtractBalance(mock.Anything, dto.BalanceOperationInput{ID: 12, Amount: amount}).
					Return(nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing user id in context",
			body:       dto.BalanceOperationInput{ID: 12, Amount: amount},
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "validation error",
			userID: "12",
			body: dto.BalanceOperationInput{
				Amount: decimal.Zero,
			},
			setup: func(mockStore *store.MockStore) {
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "validation error negative amount",
			userID: "12",
			body: dto.BalanceOperationInput{
				Amount: decimal.NewFromInt(-25),
			},
			setup: func(mockStore *store.MockStore) {
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			handler := handlers.NewWalletHandler(newWalletService(mockStore), schema.NewDecoder())

			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/me/balance/withdraw", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.userID != "" {
				req = authedRequest(req, tt.userID)
			}
			rec := httptest.NewRecorder()

			handler.WithdrawBalance(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
