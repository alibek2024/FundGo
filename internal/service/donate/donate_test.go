package donate_test

import (
	"context"
	"testing"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/alibek2024/FundGo/internal/service/donate"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDonateToCampaign(t *testing.T) {
	input := dto.DonateInput{
		UserID:     1,
		CampaignID: 2,
		Amount:     decimal.NewFromInt(25),
	}

	active := model.Active
	archived := model.Archived
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
					GetCampaignStatus(mock.Anything, int64(2)).
					Return(&active, nil)

				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(1)).
					Return(&balance, nil)

				mockStore.EXPECT().
					SubtractBalance(
						mock.Anything,
						dto.BalanceOperationInput{
							ID:     1,
							Amount: input.Amount,
						},
					).
					Return(nil)

				mockStore.EXPECT().
					IncreaseCampaignAmount(
						mock.Anything,
						dto.CampaignBalanceOperation{
							ID:     2,
							Amount: input.Amount,
						},
					).
					Return(&model.Campaign{
						ID:            2,
						CurrentAmount: decimal.NewFromInt(25),
						TargetAmount:  decimal.NewFromInt(100),
					}, nil)

				mockStore.EXPECT().
					CreateDonation(mock.Anything, input).
					Return(&model.Donation{
						ID:         10,
						UserID:     1,
						CampaignID: 2,
						Amount:     input.Amount,
					}, nil)

				mockStore.EXPECT().
					CreateTransaction(
						mock.Anything,
						mock.AnythingOfType("dto.TransactionInput"),
					).
					Return(&model.Transaction{ID: 1}, nil)
			},
		},
		{
			name: "campaign becomes successful",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)

				mockStore.EXPECT().
					GetCampaignStatus(mock.Anything, int64(2)).
					Return(&active, nil)

				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(1)).
					Return(&balance, nil)

				mockStore.EXPECT().
					SubtractBalance(
						mock.Anything,
						dto.BalanceOperationInput{
							ID:     1,
							Amount: input.Amount,
						},
					).
					Return(nil)

				mockStore.EXPECT().
					IncreaseCampaignAmount(
						mock.Anything,
						dto.CampaignBalanceOperation{
							ID:     2,
							Amount: input.Amount,
						},
					).
					Return(&model.Campaign{
						ID:            2,
						CurrentAmount: decimal.NewFromInt(100),
						TargetAmount:  decimal.NewFromInt(100),
					}, nil)

				mockStore.EXPECT().
					UpdateCampaignStatus(
						mock.Anything,
						int64(2),
						model.Successful,
					).
					Return(&model.Campaign{
						ID:     2,
						Status: model.Successful,
					}, nil)

				mockStore.EXPECT().
					CreateDonation(mock.Anything, input).
					Return(&model.Donation{
						ID:         10,
						UserID:     1,
						CampaignID: 2,
						Amount:     input.Amount,
					}, nil)

				mockStore.EXPECT().
					CreateTransaction(
						mock.Anything,
						mock.AnythingOfType("dto.TransactionInput"),
					).
					Return(&model.Transaction{ID: 1}, nil)
			},
		},
		{
			name: "inactive campaign",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)

				mockStore.EXPECT().
					GetCampaignStatus(mock.Anything, int64(2)).
					Return(&archived, nil)
			},
			wantErr: contracts.ErrCampaignNotActive,
		},
		{
			name: "insufficient balance",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)

				mockStore.EXPECT().
					GetCampaignStatus(mock.Anything, int64(2)).
					Return(&active, nil)

				mockStore.EXPECT().
					GetBalance(mock.Anything, int64(1)).
					Return(&balance, nil)

				mockStore.EXPECT().
					SubtractBalance(
						mock.Anything,
						dto.BalanceOperationInput{
							ID:     1,
							Amount: input.Amount,
						},
					).
					Return(store.ErrDataConflict)
			},
			wantErr: contracts.ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			service := donate.CreateDonateService(mockStore)

			err := service.DonateToCampaign(
				context.Background(),
				input,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRefundDonation(t *testing.T) {
	amount := decimal.NewFromInt(25)
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
					GetDonationByID(mock.Anything, int64(10)).
					Return(&model.Donation{
						ID:         10,
						UserID:     1,
						CampaignID: 2,
						Status:     model.DonationActive,
						Amount:     amount,
					}, nil)
				mockStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(2)).
					Return(&model.Campaign{ID: 2, Status: model.Active}, nil)
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(1)).
					Return(&model.User{ID: 1, Balance: decimal.NewFromInt(75)}, nil)
				mockStore.EXPECT().
					DecreaseCampaignBalance(mock.Anything, dto.CampaignBalanceOperation{ID: 2, Amount: amount}).
					Return(&model.Campaign{ID: 2}, nil)
				mockStore.EXPECT().
					AddBalance(mock.Anything, dto.BalanceOperationInput{ID: 1, Amount: amount}).
					Return(nil)
				mockStore.EXPECT().
					RefundDonationStatus(mock.Anything, int64(10)).
					Return(&model.Donation{ID: 10, Status: model.DonationRefund}, nil)
				mockStore.EXPECT().
					CreateTransaction(mock.Anything, mock.AnythingOfType("dto.TransactionInput")).
					Return(&model.Transaction{ID: 1}, nil)
			},
		},
		{
			name: "already refunded",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetDonationByID(mock.Anything, int64(10)).
					Return(&model.Donation{ID: 10, Status: model.DonationRefund}, nil)
			},
			wantErr: contracts.ErrDonateRefunded,
		},
		{
			name: "campaign not active",
			setup: func(mockStore *store.MockStore) {
				expectExecTx(mockStore)
				mockStore.EXPECT().
					GetDonationByID(mock.Anything, int64(10)).
					Return(&model.Donation{
						ID:         10,
						UserID:     1,
						CampaignID: 2,
						Status:     model.DonationActive,
						Amount:     amount,
					}, nil)
				mockStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(2)).
					Return(&model.Campaign{ID: 2, Status: model.Archived}, nil)
			},
			wantErr: contracts.ErrCampaignNotActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockStore(t)
			tt.setup(mockStore)

			service := donate.CreateDonateService(mockStore)
			err := service.RefundDonation(context.Background(), 10)

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
