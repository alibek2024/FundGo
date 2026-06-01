package service

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
)

var (
	ErrLogin = errors.New("Incorrect email or password")

	ErrDataConflict = errors.New("data conflict")

	ErrUserNotFound        = errors.New("user not found")
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrDonationNotFound    = errors.New("donation not found")
	ErrTransactionNotFound = errors.New("transaction not found")

	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserAlreadyExists  = errors.New("email already exists")

	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrCampaignClosed          = errors.New("campaign is closed")
	ErrCannotDonateOwnCampaign = errors.New("cannot donate to self")
)

type AuthUseCase interface {
	RegisterUser(ctx context.Context, input *dto.RegistrationInput) (*model.User, error)
	SignIn(ctx context.Context, input dto.SignIn) (string, error)

	Authenticate(tokenString string) (*dto.TokenClaims, error)
}

type UserUseCase interface {
	UpdateUserInfo(ctx context.Context, input *dto.UserInfo) error

	ChangePassword(ctx context.Context, input *dto.ChangeUserPassword) error
	ChangeEmail(ctx context.Context, input *dto.UserEmail) error

	DeactivateAccount(ctx context.Context, userID int64) error
	PurgeUserData(ctx context.Context, userID int64) error
}

type CampaignUseCase interface {
	CreateCampaign(ctx context.Context, input dto.CreateCampaignInput) (*model.Campaign, error)

	GetCampaignByID(ctx context.Context, campaignID int64) (*model.Campaign, error)

	SearchCampaign(ctx context.Context, name string) ([]model.Campaign, error)

	WrapUpCampaign(ctx context.Context, campaignID int64) error
	ForceDeleteCampaign(ctx context.Context, campaignID int64) error

	GetCampaignDonors(ctx context.Context, campaignID int64) ([]model.Donation, error)
}

type DonationUseCase interface {
	DonateToCampaign(ctx context.Context, input dto.DonateInput) (*model.Donation, error)
	RefundDonation(ctx context.Context, donationID int64) error
}

type TransactionUseCase interface {
	TopUpBalance(ctx context.Context, input dto.BalanceOperationInput) error
	WithdrawBalance(ctx context.Context, input dto.BalanceOperationInput) error

	GetPaymentHistory(ctx context.Context, userID int64) ([]model.Transaction, error)
}
