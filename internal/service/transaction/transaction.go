package transaction

import (
	"github.com/alibek2024/FundGo/internal/repository"
)

type TxService struct {
	Store repository.Store
}

func CreateTX(db repository.Store) *TxService {
	return &TxService{
		Store: db,
	}
}
