package store

import (
	"context"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/shopspring/decimal"
)

type UserStore interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	RestoreUser(ctx context.Context, id int64) (*model.User, error)
	SoftDeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, input model.UserInput) (*model.User, error)
	UserResponce(ctx context.Context, id int64) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error)
}

type CampaignStore interface {
	CreateCampaign(ctx context.Context, input model.CreateCampaignInput) (*model.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	GetCampaignByID(ctx context.Context, id int64) (*model.Campaign, error)
	GetCampaignByTitle(ctx context.Context, title string) (*model.Campaign, error)
	GetCampaignStatus(ctx context.Context, id int64) (*model.CampaignStatus, error)
	GetCurrentAmount(ctx context.Context, id int64) (decimal.Decimal, error)
	IncreaseCampaignAmount(ctx context.Context, input model.IncreaseCampaignAmount) (*model.Campaign, error)
	ListCampaigns(ctx context.Context, input model.PaginationParams) ([]*model.Campaign, error)
}

type DonationStore interface {
	CreateDonation(ctx context.Context, input model.DonateInput) (model.Donation, error)
	GetListDonations(ctx context.Context, campaignID int64) ([]model.Donation, error)
}

type TrasactionStore interface {
	AddBalance(ctx context.Context, input model.Amount) error
	CreateTransaction(ctx context.Context, input model.TransactionInput) (model.Transaction, error)
	GetBalance(ctx context.Context, id int64) (decimal.Decimal, error)
	SubtractBalance(ctx context.Context, input model.Amount) (int64, error)
}
