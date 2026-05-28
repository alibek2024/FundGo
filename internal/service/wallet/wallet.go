package wallet

import (
	"github.com/alibek2024/FundGo/internal/repository/store"
)

type Service struct {
	Store store.Store
}

func NewWalletService(store store.Store) *Service {
	return &Service{
		Store: store,
	}
}
