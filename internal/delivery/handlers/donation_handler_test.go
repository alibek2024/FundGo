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
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDonationHandlerDonateToCampaign(t *testing.T) {
	amount := decimal.NewFromInt(25)
	active := model.Active
	balance := decimal.NewFromInt(100)

	tests := []struct {
		name       string
		userID     string
		body       dto.DonateInput
		setup      func(*store.MockStore)
		wantStatus int
	}{
		{
			name:   "success uses authenticated user id",
			userID: "12",
			body: dto.DonateInput{
				UserID:     999,
				CampaignID: 2,
				Amount:     amount,
			},
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)

				mockStore.EXPECT().
					GetCampaignStatus(mock.Anything, int64(2)).
					Return(&active, nil)

				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(12)).
					Return(&balance, nil)

				mockStore.EXPECT().
					SubtractBalance(
						mock.Anything,
						dto.BalanceOperationInput{
							ID:     12,
							Amount: amount,
						},
					).
					Return(nil)

				mockStore.EXPECT().
					IncreaseCampaignAmount(
						mock.Anything,
						dto.CampaignBalanceOperation{
							ID:     2,
							Amount: amount,
						},
					).
					Return(&model.Campaign{
						ID:            2,
						CurrentAmount: decimal.NewFromInt(25),
						TargetAmount:  decimal.NewFromInt(100),
					}, nil)

				mockStore.EXPECT().
					CreateDonation(
						mock.Anything,
						dto.DonateInput{
							UserID:     12,
							CampaignID: 2,
							Amount:     amount,
						},
					).
					Return(&model.Donation{
						ID:         10,
						UserID:     12,
						CampaignID: 2,
						Amount:     amount,
					}, nil)

				mockStore.EXPECT().
					CreateTransaction(
						mock.Anything,
						mock.AnythingOfType("dto.TransactionInput"),
					).
					Return(&model.Transaction{ID: 1}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing authenticated user",
			body:       dto.DonateInput{UserID: 1, CampaignID: 2, Amount: amount},
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "validation error",
			userID:     "12",
			body:       dto.DonateInput{CampaignID: 2, Amount: decimal.Zero}, // Невалидный Amount (<= 0)
			setup:      func(mockStore *store.MockStore) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			handler := handlers.NewDonationHandler(
				newDonateService(mockStore),
				schema.NewDecoder(),
				newTransactionService(mockStore),
			)

			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns/2/donate", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.userID != "" {
				req = authedRequest(req, tt.userID)
			}
			req = mux.SetURLVars(req, map[string]string{"id": "2"})

			rec := httptest.NewRecorder()

			handler.DonateToCampaign(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestDonationHandlerTransactionHistory(t *testing.T) {
	mockStore := store.NewMockStore(t)
	mockStore.EXPECT().
		HistoryTX(mock.Anything, int64(12)).
		Return([]model.Transaction{{ID: 1, UserID: 12}}, nil)

	handler := handlers.NewDonationHandler(
		newDonateService(mockStore),
		schema.NewDecoder(),
		newTransactionService(mockStore),
	)
	req := authedRequest(httptest.NewRequest(http.MethodGet, "/api/v1/users/me/donations", nil), "12")
	rec := httptest.NewRecorder()

	handler.TransactionHistory(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"UserID":12`)
}

func TestDonationHandlerRefundDonationRejectsWrongOwner(t *testing.T) {
	mockStore := store.NewMockStore(t)
	differentDonationID := int64(999) 

	mockStore.EXPECT().
		HistoryTX(mock.Anything, int64(12)).
		Return([]model.Transaction{
			{ID: 1, UserID: 12, DonationID: &differentDonationID},
		}, nil)

	handler := handlers.NewDonationHandler(
		newDonateService(mockStore),
		schema.NewDecoder(),
		newTransactionService(mockStore),
	)

	req := authedRequest(httptest.NewRequest(http.MethodPost, "/api/v1/donations/10/refund", nil), "12")
	req = mux.SetURLVars(req, map[string]string{"donation_id": "10"})

	rec := httptest.NewRecorder()

	handler.RefundDonation(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
