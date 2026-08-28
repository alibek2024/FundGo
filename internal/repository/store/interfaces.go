package store

import (
	"context"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/shopspring/decimal"
)

type Store interface {
	UserStore
	CampaignStore
	DonationStore
	TrasactionStore
	WalletStore
	TransactionManager
}

type TransactionManager interface {
	ExecTx(ctx context.Context, fn func(Store) error) error
}

type UserStore interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	RestoreUser(ctx context.Context, id int64) (*model.User, error)
	SoftDeleteUser(ctx context.Context, id int64) error
	UpdateEmail(ctx context.Context, input dto.UserEmail) (*model.User, error)
	UpdateInfo(ctx context.Context, input dto.UserInfo) (*model.User, error)
	UpdatePassword(ctx context.Context, input dto.UpdateUserPassword) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateUser(ctx context.Context, input dto.RegistrationInput) (*model.User, error)
	GetByIDForPurge(ctx context.Context, id int64) (*model.User, error)
}

type CampaignStore interface {
	UpdateCampaignStatus(ctx context.Context, id int64, status model.CampaignStatus) (*model.Campaign, error)
	CreateCampaign(ctx context.Context, input dto.CreateCampaignInput) (*model.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	GetCampaignByID(ctx context.Context, id int64) (*model.Campaign, error)
	GetCampaignByTitle(ctx context.Context, title string) ([]*model.Campaign, error)
	GetCampaignStatus(ctx context.Context, id int64) (*model.CampaignStatus, error)
	GetCurrentAmount(ctx context.Context, id int64) (*decimal.Decimal, error)
	IncreaseCampaignAmount(ctx context.Context, input dto.CampaignBalanceOperation) (*model.Campaign, error)
	DecreaseCampaignBalance(ctx context.Context, input dto.CampaignBalanceOperation) (*model.Campaign, error)
	ListCampaigns(ctx context.Context, input dto.PaginationParams) ([]*model.Campaign, error)
}

type DonationStore interface {
	CreateDonation(ctx context.Context, input dto.DonateInput) (*model.Donation, error)
	GetDonationByID(ctx context.Context, id int64) (*model.Donation, error)
	GetListDonations(ctx context.Context, campaignID int64) ([]model.Donation, error)
	UpdateDonationStatus(ctx context.Context, input dto.UpdateDonationStatus) (*model.Donation, error)
	RefundDonationStatus(ctx context.Context, input int64) (*model.Donation, error)
}

type TrasactionStore interface {
	CreateTransaction(ctx context.Context, input dto.TransactionInput) (*model.Transaction, error)
	HistoryTX(ctx context.Context, userID int64) ([]model.Transaction, error)
}

type WalletStore interface {
	AddBalance(ctx context.Context, input dto.BalanceOperationInput) error
	GetBalance(ctx context.Context, id int64) (*decimal.Decimal, error)
	SubtractBalance(ctx context.Context, input dto.BalanceOperationInput) error
}
