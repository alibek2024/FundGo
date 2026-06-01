package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) generateTokenPair(user *model.User) (*dto.AuthTokens, error) {
	userID := strconv.FormatInt(user.ID, 10)
	accessTokenID := uuid.NewString()
	refreshTokenID := uuid.NewString()

	accessToken, err := s.generateToken(userID, accessTokenID, s.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(userID, refreshTokenID, s.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &dto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) generateToken(userID, tokenID string, tokenTTL time.Duration) (string, error) {
	Claims := dto.TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	TokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)
	Token, err := TokenObj.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return Token, nil
}
