package campaign_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCampaign(t *testing.T) {
	input := dto.CreateCampaignInput{
		Title:        "Clean Water",
		CreatorID:    1,
		Description:  "help",
		TargetAmount: decimal.NewFromInt(100),
	}

	tests := []struct {
		name    string
		setup   func(*store.MockCampaignStore)
		wantErr string
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					CreateCampaign(mock.Anything, input).
					Return(&model.Campaign{ID: 1, Title: input.Title}, nil)
			},
		},
		{
			name: "store error",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					CreateCampaign(mock.Anything, input).
					Return(nil, errors.New("db unavailable"))
			},
			wantErr: "create camapaign",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockCampaignStore(t)
			tt.setup(mockStore)

			service := campaign.NewCampaignService(mockStore, campaign.RefundManager{})
			got, err := service.CreateCampaign(context.Background(), input)

			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestGetCampaignByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*store.MockCampaignStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(1)).
					Return(&model.Campaign{ID: 1, Status: model.Active}, nil)
			},
		},
		{
			name: "not found",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(1)).
					Return(nil, store.ErrNotFound)
			},
			wantErr: contracts.ErrCampaignNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockCampaignStore(t)
			tt.setup(mockStore)

			service := campaign.NewCampaignService(mockStore, campaign.RefundManager{})
			got, err := service.GetCampaignByID(context.Background(), 1)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestWrapUpCampaign(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*store.MockCampaignStore, *store.MockStore)
		wantErr error
	}{
		{
			name: "successful campaign",
			setup: func(campaignStore *store.MockCampaignStore, _ *store.MockStore) {
				campaignStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(1)).
					Return(&model.Campaign{
						ID:            1,
						Status:        model.Active,
						CurrentAmount: decimal.NewFromInt(100),
						TargetAmount:  decimal.NewFromInt(100),
					}, nil)
				campaignStore.EXPECT().
					UpdateCampaignStatus(mock.Anything, int64(1), model.Successful).
					Return(&model.Campaign{ID: 1, Status: model.Successful}, nil)
			},
		},
		{
			name: "inactive campaign",
			setup: func(campaignStore *store.MockCampaignStore, _ *store.MockStore) {
				campaignStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(1)).
					Return(&model.Campaign{ID: 1, Status: model.Archived}, nil)
			},
			wantErr: contracts.ErrCampaignNotActive,
		},
		{
			name: "failed campaign refunds donations",
			setup: func(campaignStore *store.MockCampaignStore, txStore *store.MockStore) {
				campaignStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(1)).
					Return(&model.Campaign{
						ID:            1,
						Status:        model.Active,
						CurrentAmount: decimal.NewFromInt(50),
						TargetAmount:  decimal.NewFromInt(100),
					}, nil)
				txStore.EXPECT().
					ExecTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(store.Store) error) error {
						return fn(txStore)
					})
				txStore.EXPECT().
					GetListDonations(mock.Anything, int64(1)).
					Return([]model.Donation{}, nil)
				campaignStore.EXPECT().
					UpdateCampaignStatus(mock.Anything, int64(1), model.Failed).
					Return(&model.Campaign{ID: 1, Status: model.Failed}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaignStore := store.NewMockCampaignStore(t)
			txStore := store.NewMockStore(t)
			tt.setup(campaignStore, txStore)

			refundManager := campaign.NewRefundManager(txStore)
			service := campaign.NewCampaignService(campaignStore, *refundManager)
			err := service.WrapUpCampaign(context.Background(), 1)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
