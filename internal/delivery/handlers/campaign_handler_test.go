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
	campaignsvc "github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/gorilla/schema"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCampaignHandlerSearchCampaign(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setup      func(*store.MockCampaignStore)
		wantStatus int
	}{
		{
			name:  "success",
			query: "school",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					GetCampaignByTitle(mock.Anything, "school").
					Return([]*model.Campaign{{ID: 10, Title: "school supplies"}}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty query",
			query:      "   ",
			setup:      func(mockStore *store.MockCampaignStore) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "not found",
			query: "missing",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					GetCampaignByTitle(mock.Anything, "missing").
					Return(nil, store.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockCampaignStore(t)
			tt.setup(mockStore)

			handler := handlers.NewCampaignHandler(
				newCampaignService(mockStore, campaignsvc.RefundManager{}),
				schema.NewDecoder(),
			)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns?q="+url.QueryEscape(tt.query), nil)
			rec := httptest.NewRecorder()

			handler.SearchCampaign(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestCampaignHandlerWrapUpCampaign(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		setup      func(*store.MockCampaignStore)
		wantStatus int
	}{
		{
			name:   "success",
			pathID: "10",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(10)).
					Return(&model.Campaign{
						ID:            10,
						Status:        model.Active,
						CurrentAmount: decimal.NewFromInt(100),
						TargetAmount:  decimal.NewFromInt(100),
					}, nil)
				mockStore.EXPECT().
					UpdateCampaignStatus(mock.Anything, int64(10), model.Successful).
					Return(&model.Campaign{ID: 10, Status: model.Successful}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:   "not found",
			pathID: "10",
			setup: func(mockStore *store.MockCampaignStore) {
				mockStore.EXPECT().
					GetCampaignByID(mock.Anything, int64(10)).
					Return(nil, store.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bad path id",
			pathID:     "bad",
			setup:      func(mockStore *store.MockCampaignStore) {},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockCampaignStore(t)
			tt.setup(mockStore)

			handler := handlers.NewCampaignHandler(
				newCampaignService(mockStore, campaignsvc.RefundManager{}),
				schema.NewDecoder(),
			)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns/"+tt.pathID+"/wrap-up", nil)
			req.SetPathValue("id", tt.pathID)
			rec := httptest.NewRecorder()

			handler.WrapUpCampaign(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestCampaignHandlerForceDeleteCampaign(t *testing.T) {
	mockStore := store.NewMockStore(t)
	expectExecTx(mockStore)
	mockStore.EXPECT().
		GetCampaignByID(mock.Anything, int64(10)).
		Return(&model.Campaign{ID: 10, Status: model.Active}, nil)
	mockStore.EXPECT().
		GetListDonations(mock.Anything, int64(10)).
		Return([]model.Donation{}, nil)
	mockStore.EXPECT().
		UpdateCampaignStatus(mock.Anything, int64(10), model.Failed).
		Return(&model.Campaign{ID: 10, Status: model.Failed}, nil)

	refundManager := campaignsvc.RefundManager{Tx: mockStore, DonationStore: mockStore}
	handler := handlers.NewCampaignHandler(newCampaignService(mockStore, refundManager), schema.NewDecoder())
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/campaigns/10", nil)
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	handler.ForceDeleteCampaign(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCampaignHandlerCreateCampaignPanicsOnInvalidValidationTags(t *testing.T) {
	mockStore := store.NewMockCampaignStore(t)
	handler := handlers.NewCampaignHandler(
		newCampaignService(mockStore, campaignsvc.RefundManager{}),
		schema.NewDecoder(),
	)
	req := formRequest(http.MethodPost, "/api/v1/campaigns", url.Values{
		"title":         {"Go"},
		"creator_id":    {"0"},
		"description":   {"bad"},
		"target_amount": {"100"},
	})
	rec := httptest.NewRecorder()

	require.Panics(t, func() {
		handler.CreateCampaign(rec, req)
	})
}

func TestCampaignHandlerCreateCampaignSuccessPanicsOnInvalidValidationTags(t *testing.T) {
	mockStore := store.NewMockCampaignStore(t)
	input := dto.CreateCampaignInput{
		Title:        "Medical support",
		CreatorID:    7,
		Description:  "Help with treatment",
		TargetAmount: decimal.NewFromInt(100),
	}

	handler := handlers.NewCampaignHandler(
		newCampaignService(mockStore, campaignsvc.RefundManager{}),
		schema.NewDecoder(),
	)
	req := formRequest(http.MethodPost, "/api/v1/campaigns", url.Values{
		"title":         {input.Title},
		"creator_id":    {"7"},
		"description":   {input.Description},
		"target_amount": {"100"},
	})
	rec := httptest.NewRecorder()

	require.Panics(t, func() {
		handler.CreateCampaign(rec, req)
	})
}
