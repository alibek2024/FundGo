package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
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
) Service {
	return Service{
		privateKey: privateKey,
		publicKey:  publicKey,
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
		Store:      store,
	}
}

func (u *Service) RegisterUser(ctx context.Context, input *dto.RegistrationInput) (*model.User, error) {
	if err := u.CheckEmail(ctx, input.Email); err != nil {
		if errors.Is(err, contracts.ErrEmailAlreadyExists) {
			return nil, contracts.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("check email: %w", err)
	}
	hashPassword, err := hashPassword(input.HashPassword)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	input.HashPassword = string(hashPassword)

	user, err := u.Store.CreateUser(ctx, *input)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return nil, contracts.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (u *Service) SignIn(ctx context.Context, input dto.SignIn) (*dto.AuthTokens, error) {
	user, err := u.Store.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, fmt.Errorf("get user by email: %w", err)
		}
		return nil, err
	}

	if err := u.CheckPassword(input.HashPassword, user.HashPassword); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errors.New("Incorrect login or password")
		}
		return nil, fmt.Errorf("compare password hash: %w", err)
	}

	tokens, err := u.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("error generate token pair: %w", err)
	}
	return tokens, nil
}

func (s *Service) Authenticate(tokenString string) (*dto.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if claims, ok := token.Claims.(*dto.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
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
