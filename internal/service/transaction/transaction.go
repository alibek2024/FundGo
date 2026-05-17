package transaction

import (
	"context"
	"time"

	"github.com/alibek2024/FundGo/internal/repository"
)

type TxService struct {
	Store repository.Store
}

type CtxKey struct{}

type Settings interface {
	EnrichBy(external Settings) Settings
	CtxKey() CtxKey
	Propahgation()
	Cancelable() bool
	TimeOutOrNil() *time.Duration
}

type Manager interface {
	DO(context.Context, func(context.Context) error) error

	DoWithSettings(
		context.Context,
		Settings,
		func(context.Context) error,
	)
}

func CreateTX(db repository.Store) *TxService {
	return &TxService{
		Store: db,
	}
}
