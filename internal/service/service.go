package service

import (
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/alibek2024/FundGo/internal/service/donate"
	"github.com/alibek2024/FundGo/internal/service/transaction"
	"github.com/alibek2024/FundGo/internal/service/user"
	"github.com/alibek2024/FundGo/internal/service/wallet"
)

type Service struct {
	contracts.UserUseCase
	contracts.CampaignUseCase
	contracts.DonationUseCase
	contracts.TransactionUseCase
	contracts.WalletUseCase
}

func CreateService(store store.Store) *Service {
	return &Service{
		UserUseCase: user.NewUserService(store),
		CampaignUseCase: campaign.CreateService(
			store,
			*campaign.NewRefundManager(
				store, store,
			),
		),
		DonationUseCase:    donate.CreateDonateService(store),
		TransactionUseCase: transaction.CreateTX(store),
		WalletUseCase:      wallet.NewWalletService(store),
	}
}
