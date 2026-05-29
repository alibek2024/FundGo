package service

type ServiceUseCase interface {
	UserUseCase
	CampaignUseCase
	DonationUseCase
	TransactionUseCase
}

// type Service struct {
// 	CampaignService    *campaign.Service
// 	DonationService    *donate.Service
// 	TransactionService *transaction.Service
// 	UserService        *user.Service
// 	WalletService      *wallet.Service
// }

// func CreateService(
// 	databaseStore store.Store,
// ) *Service {
// 	return &Service{
// 		CampaignService:    campaign.NewCampaignService(databaseStore),
// 		DonationService:    donate.CreateDonateService(databaseStore),
// 		TransactionService: transaction.CreateTX(databaseStore),
// 		UserService:        user.NewUserService(databaseStore),
// 		WalletService:      wallet.NewWalletService(databaseStore),
// 	}
// }
