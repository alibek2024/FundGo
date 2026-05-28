package service

import (
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/alibek2024/FundGo/internal/service/donate"
	"github.com/alibek2024/FundGo/internal/service/transaction"
	"github.com/alibek2024/FundGo/internal/service/user"
	"github.com/alibek2024/FundGo/internal/service/wallet"
)

type Service struct {
	CampaignService    *campaign.Service
	DonationService    *donate.Service
	TransactionService *transaction.Service
	UserService        *user.Service
	WalletService      *wallet.Service
}

func CreateService(
	databaseStore store.Store,
) *Service {
	return &Service{
		CampaignService:    campaign.NewCampaignService(databaseStore),
		DonationService:    donate.CreateDonateService(databaseStore),
		TransactionService: transaction.CreateTX(databaseStore),
		UserService:        user.NewUserService(databaseStore),
		WalletService:      wallet.NewWalletService(databaseStore),
	}
}
