package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	tokenTTL   time.Duration
	refreshTTL time.Duration
	Store      store.UserStore
}

func NewUserService(
	store store.UserStore,
	tokenTTL, refreshTTL time.Duration,
	privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey,
) *Service {
	return &Service{
		privateKey: privateKey,
		publicKey:  publicKey,
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
		Store:      store,
	}
}

func (u *Service) RegisterUser(ctx context.Context, input *dto.RegistrationInput) (*dto.UserResponse, *dto.AuthTokens, error) {
	if err := u.CheckEmail(ctx, input.Email); err != nil {
		if errors.Is(err, contracts.ErrEmailAlreadyExists) {
			return nil, nil, contracts.ErrEmailAlreadyExists
		}
		return nil, nil, fmt.Errorf("check email: %w", err)
	}
	hashPassword, err := hashPassword(input.HashPassword)
	if err != nil {
		return nil, nil, fmt.Errorf("hash password: %w", err)
	}

	input.HashPassword = string(hashPassword)

	modelUser, err := u.Store.CreateUser(ctx, *input)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return nil, nil, contracts.ErrUserAlreadyExists
		}
		return nil, nil, fmt.Errorf("create user: %w", err)
	}
	tokens, err := u.generateTokenPair(modelUser.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("error generate token pair: %w", err)
	}
	user := dto.UserResponse{
		FirstName: modelUser.FirstName,
		LastName:  modelUser.LastName,
		Email:     modelUser.Email,
		ID:        modelUser.ID,
		Balance:   modelUser.Balance,
		CreatedAt: modelUser.CreatedAt,
		DeletedAt: modelUser.DeletedAt,
	}
	return &user, tokens, nil
}

func (u *Service) SignIn(ctx context.Context, input dto.SignIn) (*dto.UserResponse, *dto.AuthTokens, error) {
	modelUser, err := u.Store.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, nil, contracts.ErrLogin
		}
		return nil, nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := u.CheckPassword(input.Password, modelUser.HashPassword); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil, contracts.ErrLogin
		}
		return nil, nil, fmt.Errorf("compare password hash: %w", err)
	}

	tokens, err := u.generateTokenPair(modelUser.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("error generate token pair: %w", err)
	}
	user := dto.UserResponse{
		FirstName: modelUser.FirstName,
		LastName:  modelUser.LastName,
		Email:     modelUser.Email,
		ID:        modelUser.ID,
		Balance:   modelUser.Balance,
		CreatedAt: modelUser.CreatedAt,
		DeletedAt: modelUser.DeletedAt,
	}
	return &user, tokens, nil
}

func (s *Service) Authenticate(tokenString string) (*dto.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, contracts.ErrTokenExpired
		}
		return nil, contracts.InvalidToken
	}

	claims, ok := token.Claims.(*dto.TokenClaims)
	if !ok || !token.Valid {
		return nil, contracts.InvalidToken
	}

	return claims, nil
}

func (u *Service) GetAccessToken(userID int64) (*dto.AuthTokens, error) {
	tokens, err := u.generateTokenPair(userID)
	if err != nil {
		return nil, fmt.Errorf("error generate token pair: %w", err)
	}
	return tokens, nil
}

func (u *Service) CheckEmail(ctx context.Context, Email string) error {
	user, err := u.Store.GetByEmail(ctx, Email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("get user by email: %w", err)
	}
	if user != nil {
		return contracts.ErrEmailAlreadyExists
	}
	return nil
}

func (u *Service) CheckPassword(inputPassword, hashpassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(inputPassword))
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) ([]byte, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return HashPassword, nil
}
