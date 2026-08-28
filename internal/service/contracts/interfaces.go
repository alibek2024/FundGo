package contracts

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
)

var (
	ErrTokenExpired = errors.New("token expired")
	InvalidToken    = errors.New("invalid token")
	ErrLogin        = errors.New("Incorrect email or password")

	ErrDataConflict = errors.New("data conflict")

	ErrUserNotFound        = errors.New("user not found")
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrDonationNotFound    = errors.New("donation not found")
	ErrTransactionNotFound = errors.New("transaction not found")

	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserAlreadyExists  = errors.New("email already exists")

	ErrCampaignNotActive = errors.New("campaign is not active")
	ErrCampaignClosed    = errors.New("campaign is closed")

	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrCannotDonateOwnCampaign = errors.New("cannot donate to self")

	ErrDonateRefunded = errors.New("Donation error: Refunded")
	ErrDonationID     = errors.New("Donation error: wrong donation id")
)

type AuthUseCase interface {
	RegisterUser(ctx context.Context, input *dto.RegistrationInput) (*dto.UserResponse, *dto.AuthTokens, error)
	SignIn(ctx context.Context, input dto.SignIn) (*dto.UserResponse, *dto.AuthTokens, error)

	Authenticate(tokenString string) (*dto.TokenClaims, error)
	GetAccessToken(userID int64) (*dto.AuthTokens, error)
}

type UserUseCase interface {
	UpdateUserInfo(ctx context.Context, input *dto.UserInfo) error
	UserInfo(ctx context.Context, id int64) (*dto.UserResponse, error)

	ChangePassword(ctx context.Context, input *dto.ChangeUserPassword) error
	ChangeEmail(ctx context.Context, input *dto.UserEmail) error

	DeactivateAccount(ctx context.Context, userID int64) error
	PurgeUserData(ctx context.Context, userID int64) error
}

type CampaignUseCase interface {
	CreateCampaign(ctx context.Context, input dto.CreateCampaignInput) (*model.Campaign, error)

	GetCampaignByID(ctx context.Context, campaignID int64) (*model.Campaign, error)
	SearchCampaign(ctx context.Context, name string) ([]*model.Campaign, error)

	WrapUpCampaign(ctx context.Context, campaignID int64) error
	ForceDeleteCampaign(ctx context.Context, campaignID int64) error
}

type DonationUseCase interface {
	DonateToCampaign(ctx context.Context, input dto.DonateInput) error
	RefundDonation(ctx context.Context, donationID int64) error
}

type TransactionUseCase interface {
	GetPaymentHistory(ctx context.Context, userID int64) ([]model.Transaction, error)
	CheckDonation(ctx context.Context, userID, donationID int64) error
}

type WalletUseCase interface {
	TopUpBalance(ctx context.Context, input dto.BalanceOperationInput) error
	WithdrawBalance(ctx context.Context, input dto.BalanceOperationInput) error
}
