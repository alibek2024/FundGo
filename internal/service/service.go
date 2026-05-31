package service

type ServiceUseCase interface {
	UserUseCase
	CampaignUseCase
	DonationUseCase
	TransactionUseCase
}

// type Service struct {
// 	UserService     user.Service
// 	WalletService   wallet.Service
// 	CampaignService campaign.Service
// 	DonationService donation.Service
// 	TransactionService transaction.Service
// }
