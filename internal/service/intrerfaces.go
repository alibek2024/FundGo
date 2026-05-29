package service

import (
	"context"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, input dto.RegistrationInput) (*model.User, error)
	LoginUser(ctx context.Context, input dto.LoginInput) (*model.User, error)

	UpdateUserInfo(ctx context.Context, input dto.UserInfo) (*model.User, error)

	ChangePassword(ctx context.Context, input dto.ChangeUserPassword) error
	ChangeEmail(ctx context.Context, input dto.UserEmail) error
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
