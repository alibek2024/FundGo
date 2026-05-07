package repository

import (
	"context"

	"github.com/alibek2024/FundGo/internal/model"
)

type Repository interface {
	User
	Campaigns
	Donations
}

type User interface {
	CreateUser(ctx context.Context, userInput model.UserInput) (model.User, error)
	UpdateUser(ctx context.Context, userInput model.UserInput) (model.User, error)
	DeleteUser(ctx context.Context, userID int32) (error)
	TopUp(ctx context.Context, userID int32, amount float64) (model.User, error)
	Withdraw(ctx context.Context, userID int32, amount float64) (model.User, error)
	GetByID(ctx context.Context, userID int32) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type Campaigns interface {
	CreateCampaign(ctx context.Context,campaignInput model.CreateCmapaign) (model.Campaign, error)
	GetCurrentAmount(ctx context.Context, campaignsID int32) (model.Campaign, error)
	GetListDonations(ctx context.Context, campaignsID int32) ([]model.Donation, error)
	DeleteCampaign(ctx context.Context, campaignsID int32) (error)
}

type Donations interface {
	DonateToCampaign(ctx context.Context, donateInput model.DonateInput) (model.Donation, error)
}