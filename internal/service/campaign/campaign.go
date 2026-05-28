package campaign

import (
	"github.com/alibek2024/FundGo/internal/repository/store"
)

type Service struct {
	Store store.Store
}

func NewCampaignService(store store.Store) *Service {
	return &Service{
		Store: store,
	}
}
