package delivery

import (
	"github.com/alibek2024/FundGo/internal/delivery/handlers"
	"github.com/alibek2024/FundGo/internal/service"
	"github.com/gorilla/schema"
)

type Delivery struct {
	AuthHandler     handlers.AuthHandler
	UserHandler     handlers.UserHandler
	CampaignHandler handlers.CampaignHandler
	DonationHandler handlers.DonationHandler
	WalletHandler   handlers.WalletHandler
}

func NewDelivery(service *service.Service, decoder schema.Decoder) *Delivery {
	return &Delivery{
		AuthHandler:     handlers.NewAuthHandler(&service.Auth, &decoder),
		UserHandler:     handlers.NewUserHandler(&service.User, &decoder),
		CampaignHandler: handlers.NewCampaignHandler(&service.Campaign, &decoder),
		DonationHandler: handlers.NewDonationHandler(&service.Donation, &decoder, &service.Transaction),
		WalletHandler:   handlers.NewWalletHandler(&service.Wallet, &decoder),
	}
}
