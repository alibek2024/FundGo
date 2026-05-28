package transaction

import (
	"github.com/alibek2024/FundGo/internal/repository/store"
)

type Service struct {
	Store store.Store
}

func CreateTX(db store.Store) *Service {
	return &Service{
		Store: db,
	}
}
