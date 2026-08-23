package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alibek2024/FundGo/internal/delivery/middleware"
	"github.com/alibek2024/FundGo/internal/dto"
)

type mockService struct {
	AuthenticateFunc   func(tokenString string) (*dto.TokenClaims, error)
	RegisterUserFunc   func(ctx context.Context, input *dto.RegistrationInput) (*dto.UserResponse, *dto.AuthTokens, error)
	SignInFunc         func(ctx context.Context, input dto.SignIn) (*dto.UserResponse, *dto.AuthTokens, error)
	GetAccessTokenFunc func(userID int64) (*dto.AuthTokens, error)
}

func (m *mockService) RegisterUser(ctx context.Context, input *dto.RegistrationInput) (*dto.UserResponse, *dto.AuthTokens, error) {
	return m.RegisterUserFunc(ctx, input)
}
func (m *mockService) SignIn(ctx context.Context, input dto.SignIn) (*dto.UserResponse, *dto.AuthTokens, error) {
	return m.SignInFunc(ctx, input)
}
func (m *mockService) Authenticate(tokenString string) (*dto.TokenClaims, error) {
	return m.AuthenticateFunc(tokenString)
}
func (m *mockService) GetAccessToken(userID int64) (*dto.AuthTokens, error) {
	return m.GetAccessTokenFunc(userID)
}

func TestAuthMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	tc := []struct {
		name           string
		auth           mockService
		makeRequest    func() *http.Request
		expectedStatus int
	}{
		{
			name: "success auth",
			auth: mockService{
				AuthenticateFunc: func(tokenString string) (*dto.TokenClaims, error) {
					return &dto.TokenClaims{}, nil
				},
			},
			makeRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/protected", nil)
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "valid_token_value",
				})
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "The token cookie is missing",
			auth: mockService{},
			makeRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/protected", nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "The cookie exists, value is empty",
			auth: mockService{
				AuthenticateFunc: func(tokenString string) (*dto.TokenClaims, error) {
					return nil, errors.New("empty token")
				},
			},
			makeRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/protected", nil)
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "",
				})
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid Token Error",
			auth: mockService{
				AuthenticateFunc: func(tokenString string) (*dto.TokenClaims, error) {
					return nil, errors.New("invalid token")
				},
			},
			makeRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/protected", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "bad_token"})
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			middleware := middleware.NewAuthMiddleware(&tt.auth)
			handler := middleware.CheckAuthorization(nextHandler)
			rr := httptest.NewRecorder()
			req := tt.makeRequest()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.expectedStatus)
			}
		})
	}
}

func TestErrorHandlerMiddleware(t *testing.T) {
	tc := []struct {
		name           string
		nextHandler    http.HandlerFunc
		expectedStatus int
		expectLog      bool
		logContains    string
	}{
		{
			name: "Success - no panic",
			nextHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			},
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Recover from string panic",
			nextHandler: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
			logContains:    "something went wrong",
		},
		{
			name: "Recover from error nil pointer panic",
			nextHandler: func(w http.ResponseWriter, r *http.Request) {
				var ptr *int
				*ptr = 42
			},
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
			logContains:    "runtime error: invalid memory address or nil pointer dereference",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, nil))
			mw := middleware.NewErrorMiddleware(logger)

			handler := mw.ErrorHandlerMiddleware(tt.nextHandler)
			req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			if rec.Code != tt.expectedStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.expectedStatus)
			}

			logOutput := logBuf.String()

			if tt.expectLog {
				if logOutput == "" {
					t.Errorf("expected log output, got empty")
				}
				if !strings.Contains(logOutput, tt.logContains) {
					t.Errorf("expected log to contain %q, got %q", tt.logContains, logOutput)
				}
				if !strings.Contains(logOutput, "panic recovered") {
					t.Errorf("expected log to contain 'panic recovered'")
				}
			} else {

				if logOutput != "" {
					t.Errorf("expected no log output, got %q", logOutput)
				}
			}
		})
	}
}
func TestLoggerMiddleware(t *testing.T) {
	tc := []struct {
		name           string
		handlerStatus  int
		responseBody   string
		expectedLevel  string
		expectedStatus string
	}{
		{
			name:           "Status 200 OK - INFO level",
			handlerStatus:  http.StatusOK,
			responseBody:   "Success payload",
			expectedLevel:  "level=INFO",
			expectedStatus: "status=200",
		},
		{
			name:           "Status 404 Not Found - WARN level",
			handlerStatus:  http.StatusNotFound,
			responseBody:   "Not Found",
			expectedLevel:  "level=WARN",
			expectedStatus: "status=404",
		},
		{
			name:           "Status 500 Internal Error - ERROR level",
			handlerStatus:  http.StatusInternalServerError,
			responseBody:   "Internal Error",
			expectedLevel:  "level=ERROR",
			expectedStatus: "status=500",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, nil))

			mw := middleware.NewLoggerMiddleware(logger)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			})

			handler := mw.LogMiddleware(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/resource", nil)
			req.Header.Set("User-Agent", "TestAgent/1.0")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.handlerStatus {
				t.Errorf("got HTTP status %d, want %d", rec.Code, tt.handlerStatus)
			}
			if rec.Body.String() != tt.responseBody {
				t.Errorf("got body %q, want %q", rec.Body.String(), tt.responseBody)
			}

			logOutput := logBuf.String()

			expectedFields := []string{
				"msg=\"HTTP request handled\"",
				tt.expectedLevel,
				tt.expectedStatus,
				"method=GET",
				"path=/api/v1/resource",
				"user_agent=TestAgent/1.0",
			}

			for _, field := range expectedFields {
				if !strings.Contains(logOutput, field) {
					t.Errorf("expected log output to contain %q, got: %s", field, logOutput)
				}
			}
		})
	}
}
