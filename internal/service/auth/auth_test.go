package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/auth"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockStore *store.MockUserStore)
		wantErr error
	}{
		{
			name: "success",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "test@example.com").
					Return(nil, store.ErrNotFound)

				mockStore.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("dto.RegistrationInput")).
					Return(&model.User{
						ID:    1,
						Email: "test@example.com",
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "test@example.com").
					Return(&model.User{ID: 1, Email: "test@example.com"}, nil)
			},
			wantErr: contracts.ErrEmailAlreadyExists,
		},
		{
			name: "create user already exists",
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "test@example.com").
					Return(nil, store.ErrNotFound)

				mockStore.EXPECT().
					CreateUser(mock.Anything, mock.AnythingOfType("dto.RegistrationInput")).
					Return(nil, store.ErrAlreadyExists)
			},
			wantErr: contracts.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err)
			mockStore := store.NewMockUserStore(t)

			tt.setup(mockStore)

			service := auth.NewUserService(
				mockStore,
				time.Hour,
				time.Hour*24,
				privateKey,
				&privateKey.PublicKey,
			)

			input := &dto.RegistrationInput{
				FirstName: "test",
				LastName:  "Test",
				Email:     "test@example.com",
				Password:  "password123",
			}
			user, _, err := service.RegisterUser(context.Background(), input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
		})
	}
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name    string
		input   dto.SignIn
		setup   func(mockStore *store.MockUserStore)
		wantErr error
	}{
		{
			name: "success login",
			input: dto.SignIn{
				Email:    "test@example.ru",
				Password: "password123",
			},
			setup: func(mockStore *store.MockUserStore) {
				hashed, _ := bcrypt.GenerateFromPassword(
					[]byte("password123"),
					12,
				)

				mockStore.EXPECT().
					GetByEmail(mock.Anything, "test@example.ru").
					Return(&model.User{
						Email:        "test@example.ru",
						HashPassword: string(hashed),
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			input: dto.SignIn{
				Email:    "tests@example.ru",
				Password: "password123",
			},
			setup: func(mockStore *store.MockUserStore) {
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "tests@example.ru").
					Return(nil, store.ErrNotFound)
			},
			wantErr: contracts.ErrLogin,
		},
		{
			name: "wrong password",
			input: dto.SignIn{
				Email:    "test@example.ru",
				Password: "password1234",
			},
			setup: func(mockStore *store.MockUserStore) {
				hashed, _ := bcrypt.GenerateFromPassword(
					[]byte("password123"),
					12,
				)
				mockStore.EXPECT().
					GetByEmail(mock.Anything, "test@example.ru").
					Return(&model.User{
						Email:        "test@example.ru",
						HashPassword: string(hashed),
					}, nil)
			},
			wantErr: contracts.ErrLogin,
		},
	}

	for _, tt := range tests {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			mockStore := store.NewMockUserStore(t)

			tt.setup(mockStore)

			service := auth.NewUserService(
				mockStore,
				time.Hour,
				time.Hour*24,
				privateKey,
				&privateKey.PublicKey,
			)
			user, _, err := service.SignIn(context.Background(), tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	wrongPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	generateTestToken := func(key *rsa.PrivateKey, claims dto.TokenClaims) string {
		TokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token, err := TokenObj.SignedString(key)
		require.NoError(t, err)
		return token
	}

	tokenID := uuid.NewString()
	store := store.NewMockUserStore(t)

	service := auth.NewUserService(
		store,
		time.Hour,
		time.Hour*30,
		privateKey,
		&privateKey.PublicKey,
	)

	tests := []struct {
		name        string
		tokenString string
		wantErr     error
	}{
		{
			name: "success",
			tokenString: generateTestToken(privateKey, dto.TokenClaims{
				UserID: "user-uuid-123",
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        tokenID,
					Subject:   "user-uuid-123",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			wantErr: nil,
		},
		{
			name: "token is expired",
			tokenString: generateTestToken(privateKey, dto.TokenClaims{
				UserID: "user-uuid-123",
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        tokenID,
					Subject:   "user-uuid-123",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			wantErr: contracts.ErrTokenExpired,
		},
		{
			name:        "Malformed Token",
			tokenString: "not-a-valid-jwt-token-string",
			wantErr:     contracts.InvalidToken,
		},
		{
			name: "invalid signature signed with wrong key",
			tokenString: generateTestToken(wrongPrivateKey, dto.TokenClaims{
				UserID: "user-uuid-123",
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        tokenID,
					Subject:   "user-uuid-123",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			wantErr: contracts.InvalidToken,
		},
		{
			name: "wrong signing method (HMAC instead of RSA)",
			tokenString: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": "user-uuid-123",
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				str, _ := token.SignedString([]byte("secret-key-attack"))
				return str
			}(),
			wantErr: contracts.InvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.Authenticate(tt.tokenString)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, claims)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, claims)
			require.Equal(t, "user-uuid-123", claims.UserID)
		})
	}

}
