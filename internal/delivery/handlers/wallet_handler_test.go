package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	tests := []struct {
		name       string
		userID     string
		form       url.Values
		setup      func(*store.MockStore)
		wantStatus int
	}{
		{
			name:   "success uses authenticated user id",
			userID: "12",
			form: url.Values{
				"id":     {"999"},
				"amount": {"25"},
			},
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(12)).
					Return(&model.User{ID: 12, Balance: decimal.NewFromInt(100)}, nil)
				mockStore.EXPECT().
					AddBalance(mock.Anything, dto.BalanceOperationInput{ID: 12, Amount: decimal.NewFromInt(25)}).
					Return(nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing user id in context",
			form:       url.Values{"ID": {"12"}, "Amount": {"25"}},
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "validation error",
			userID:     "12",
			form:       url.Values{"Amount": {"25"}},
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			handler := handlers.NewWalletHandler(newWalletService(mockStore), schema.NewDecoder())
			req := formRequest(http.MethodPost, "/api/v1/users/me/balance/top-up", tt.form)
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
	tests := []struct {
		name       string
		userID     string
		form       url.Values
		setup      func(*store.MockStore)
		wantStatus int
	}{
		{
			name:   "success uses authenticated user id",
			userID: "12",
			form: url.Values{
				"id":     {"999"},
				"amount": {"25"},
			},
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(12)).
					Return(&balance, nil)
				mockStore.EXPECT().
					SubtractBalance(mock.Anything, dto.BalanceOperationInput{ID: 12, Amount: decimal.NewFromInt(25)}).
					Return(nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing user id in context",
			form:       url.Values{"ID": {"12"}, "Amount": {"25"}},
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			handler := handlers.NewWalletHandler(newWalletService(mockStore), schema.NewDecoder())
			req := formRequest(http.MethodPost, "/api/v1/users/me/balance/withdraw", tt.form)
			if tt.userID != "" {
				req = authedRequest(req, tt.userID)
			}
			rec := httptest.NewRecorder()

			handler.WithdrawBalance(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
