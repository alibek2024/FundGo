package donate

import (
	"errors"

	"github.com/alibek2024/FundGo/internal/repository/store"
)

var ErrCampaignInactive = errors.New("campaign inactive")
var ErrInsufficientFunds = errors.New("insufficient funds")

type Service struct {
	Store store.Store
}

func CreateDonateService(store store.Store) *Service {
	return &Service{
		Store: store,
	}
}
