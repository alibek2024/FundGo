package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/alibek2024/FundGo/internal/delivery/handlers"
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/gorilla/schema"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandlerGetUserInfo(t *testing.T) {
	deletedAt := time.Time{}

	tests := []struct {
		name       string
		userID     string
		setup      func(*store.MockUserStore)
		wantStatus int
	}{
		{
			name:   "success",
			userID: "7",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(7)).
					Return(&model.User{
						ID:        7,
						FirstName: "Grace",
						LastName:  "Hopper",
						Email:     "grace@example.com",
						DeletedAt: &deletedAt,
					}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing user id in context",
			setup:      func(mockStore *store.MockUserStore) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "not found from service",
			userID: "404",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(404)).
					Return(nil, store.ErrNotFound)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			handler := handlers.NewUserHandler(newUserService(mockStore), schema.NewDecoder())
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
			if tt.userID != "" {
				req = authedRequest(req, tt.userID)
			}
			rec := httptest.NewRecorder()

			handler.GetUserInfo(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestUserHandlerUpdateInfo(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		form       url.Values
		setup      func(*store.MockUserStore)
		wantStatus int
	}{
		{
			name:   "success",
			userID: "7",
			form: url.Values{
				"first_name": {"Grace"},
				"last_name":  {"Hopper"},
			},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdateInfo(mock.Anything, dto.UserInfo{ID: 7, FirstName: "Grace", LastName: "Hopper"}).
					Return(&model.User{ID: 7}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "validation error",
			userID: "7",
			form: url.Values{
				"first_name": {"G"},
				"last_name":  {"H"},
			},
			setup:      func(mockStore *store.MockUserStore) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "missing user id",
			form:       url.Values{"first_name": {"Grace"}, "last_name": {"Hopper"}},
			setup:      func(mockStore *store.MockUserStore) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			handler := handlers.NewUserHandler(newUserService(mockStore), schema.NewDecoder())
			req := formRequest(http.MethodPatch, "/api/v1/users/me", tt.form)
			if tt.userID != "" {
				req = authedRequest(req, tt.userID)
			}
			rec := httptest.NewRecorder()

			handler.UpdateInfo(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestUserHandlerAccountMutations(t *testing.T) {
	deletedAt := time.Now()

	tests := []struct {
		name       string
		call       func(*handlers.UserHandler, http.ResponseWriter, *http.Request)
		form       url.Values
		setup      func(*store.MockUserStore)
		wantStatus int
	}{
		{
			name: "change email",
			call: (*handlers.UserHandler).ChangeEmail,
			form: url.Values{"email": {"new@example.com"}},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdateEmail(mock.Anything, dto.UserEmail{ID: 7, Email: "new@example.com"}).
					Return(&model.User{ID: 7, Email: "new@example.com"}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "change password",
			call: (*handlers.UserHandler).ChangePassword,
			form: url.Values{"hash_password": {"secret123"}},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					UpdatePassword(mock.Anything, dto.ChangeUserPassword{ID: 7, HashPassword: "secret123"}).
					Return(&model.User{ID: 7}, nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "deactivate account",
			call: (*handlers.UserHandler).DeactivateAccount,
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(7)).
					Return(&model.User{ID: 7, Balance: decimal.Zero}, nil)
				mockStore.EXPECT().
					SoftDeleteUser(mock.Anything, int64(7)).
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "delete account",
			call: (*handlers.UserHandler).DeleteAccount,
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByID(mock.Anything, int64(7)).
					Return(&model.User{ID: 7, Balance: decimal.Zero, DeletedAt: &deletedAt}, nil)
				mockStore.EXPECT().
					DeleteUser(mock.Anything, int64(7)).
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			handler := handlers.NewUserHandler(newUserService(mockStore), schema.NewDecoder())
			req := authedRequest(formRequest(http.MethodPatch, "/api/v1/users/me", tt.form), "7")
			rec := httptest.NewRecorder()

			tt.call(&handler, rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
