package service

import (
	"fmt"

	"github.com/alibek2024/FundGo/internal/config"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/alibek2024/FundGo/internal/service/auth"
	"github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/alibek2024/FundGo/internal/service/donate"
	"github.com/alibek2024/FundGo/internal/service/transaction"
	"github.com/alibek2024/FundGo/internal/service/user"
	"github.com/alibek2024/FundGo/internal/service/wallet"
	"github.com/hibiken/asynq"
)

type Service struct {
	Auth        auth.Service
	User        user.Service
	Campaign    campaign.CampaignScheduler
	Wallet      wallet.Service
	Donation    donate.Service
	Transaction transaction.Service
}

func NewService(
	store store.Store,
	cfg config.JWTConfig,
	client *asynq.Client,
) (Service, error) {
	privateKey, pubKey, err := cfg.RsaKeys()
	if err != nil {
		return Service{}, fmt.Errorf("load RSA keys: %w", err)
	}

	refund := campaign.NewRefundManager(store)

	campaignService := campaign.NewCampaignService(store, *refund)
	campaignScheduler := campaign.NewCampaignScheduler(*campaignService, client)

	return Service{
		Auth: *auth.NewUserService(
			store,
			cfg.AccessTokenTTL,
			cfg.RefreshTokenTTL,
			privateKey,
			pubKey,
		),
		User:        *user.NewUserService(store),
		Campaign:    *campaignScheduler,
		Wallet:      *wallet.NewWalletService(store),
		Donation:    *donate.CreateDonateService(store),
		Transaction: *transaction.CreateTX(store),
	}, nil
}
