package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alibek2024/FundGo/internal/delivery/handlers"
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandlerRegistration(t *testing.T) {
	tests := []struct {
		name       string
		form       url.Values
		setup      func(*store.MockUserStore)
		wantStatus int
		wantCookie bool
	}{
		{
			name: "success",
			form: url.Values{
				"first_name": {"Ada"},
				"last_name":  {"Lovelace"},
				"email":      {"ada@example.com"},
				"password":   {"secret123"},
			},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "ada@example.com").
					Return(nil, store.ErrNotFound)
				mockStore.EXPECT().
					CreateUser(mock.Anything, mock.MatchedBy(func(input dto.RegistrationInput) bool {
						return input.FirstName == "Ada" &&
							input.LastName == "Lovelace" &&
							input.Email == "ada@example.com" &&
							bcrypt.CompareHashAndPassword([]byte(input.HashPassword), []byte("secret123")) == nil
					})).
					Return(zeroDeletedAtUser(42), nil)
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		{
			name: "validation error",
			form: url.Values{
				"first_name": {"A"},
				"last_name":  {"L"},
				"Email":      {"not-email"},
				"password":   {"short"},
			},
			setup:      func(mockStore *store.MockUserStore) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "duplicate email",
			form: url.Values{
				"first_name": {"Ada"},
				"last_name":  {"Lovelace"},
				"Email":      {"ada@example.com"},
				"password":   {"secret123"},
			},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "ada@example.com").
					Return(&model.User{ID: 1, Email: "ada@example.com"}, nil)
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			handler := handlers.NewAuthHandler(newAuthService(t, mockStore), schema.NewDecoder())
			rec := httptest.NewRecorder()

			handler.Registration(rec, formRequest(http.MethodPost, "/api/v1/auth/register", tt.form))

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantCookie {
				require.NotEmpty(t, findCookie(rec.Result().Cookies(), "access_token"))
				require.NotEmpty(t, findCookie(rec.Result().Cookies(), "refresh_token"))

				var body dto.UserResponse
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
				require.Equal(t, int64(42), body.ID)
				require.Equal(t, "ada@example.com", body.Email)
			}
		})
	}
}

func TestAuthHandlerAuthentication(t *testing.T) {
	tests := []struct {
		name       string
		form       url.Values
		setup      func(*store.MockUserStore)
		wantStatus int
		wantCookie bool
	}{
		{
			name: "success",
			form: url.Values{
				"id":         {"1"},
				"first_name": {"Ada"},
				"last_name":  {"Lovelace"},
				"email":      {"ada@example.com"},
				"password":   {"secret123"},
			},
			setup: func(mockStore *store.MockUserStore) {
				hashed, err := bcrypt.GenerateFromPassword([]byte("secret123"), 12)
				require.NoError(t, err)

				user := zeroDeletedAtUser(42)
				user.HashPassword = string(hashed)
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "ada@example.com").
					Return(user, nil)
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		{
			name: "bad credentials",
			form: url.Values{
				"id":         {"1"},
				"first_name": {"Ada"},
				"last_name":  {"Lovelace"},
				"email":      {"ada@example.com"},
				"password":   {"secret123"},
			},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "ada@example.com").
					Return(nil, store.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "validation error",
			form: url.Values{
				"email":    {"bad-email"},
				"password": {"short"},
			},
			setup:      func(mockStore *store.MockUserStore) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)
			tt.setup(mockStore)

			handler := handlers.NewAuthHandler(newAuthService(t, mockStore), schema.NewDecoder())
			rec := httptest.NewRecorder()

			handler.Authentication(rec, formRequest(http.MethodPost, "/api/v1/auth/login", tt.form))

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantCookie {
				require.NotEmpty(t, findCookie(rec.Result().Cookies(), "access_token"))
				require.NotEmpty(t, findCookie(rec.Result().Cookies(), "refresh_token"))
			}
		})
	}
}

func TestAuthHandlerRefresh(t *testing.T) {
	t.Run("missing refresh cookie", func(t *testing.T) {
		mockStore := store.NewMockUserStore(t)
		handler := handlers.NewAuthHandler(newAuthService(t, mockStore), schema.NewDecoder())
		rec := httptest.NewRecorder()

		handler.Refresh(rec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		assertJSONError(t, rec)
	})

	t.Run("success rotates access cookie", func(t *testing.T) {
		mockStore := store.NewMockUserStore(t)
		service := newAuthService(t, mockStore)
		tokens, err := service.GetAccessToken(42)
		require.NoError(t, err)

		handler := handlers.NewAuthHandler(service, schema.NewDecoder())
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken})
		rec := httptest.NewRecorder()

		handler.Refresh(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)
		require.NotEmpty(t, findCookie(rec.Result().Cookies(), "access_token"))
	})

	t.Run("invalid refresh cookie clears auth cookies", func(t *testing.T) {
		mockStore := store.NewMockUserStore(t)
		handler := handlers.NewAuthHandler(newAuthService(t, mockStore), schema.NewDecoder())
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "broken"})
		rec := httptest.NewRecorder()

		handler.Refresh(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.Empty(t, findCookie(rec.Result().Cookies(), "access_token"))
		require.Empty(t, findCookie(rec.Result().Cookies(), "refresh_token"))
	})
}

func TestAuthHandlerAuthenticationMapsLoginError(t *testing.T) {
	require.ErrorIs(t, contracts.ErrLogin, contracts.ErrLogin)
}
