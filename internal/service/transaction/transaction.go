package transaction

import (
	"context"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
)

type TxService struct {
	Store repository.Store
}

func CreateTX(db repository.Store) *TxService {
	return &TxService{
		Store: db,
	}
}

func (t *TxService) TxHistory(
	ctx context.Context,
	userID int32,
) ([]*model.Transaction, error) {
	params := helperfunc.Int(userID)
	postgresArgs, err := t.Store.HistoryTX(ctx, params)
	if err != nil {
		return nil, err
	}
	transaction := t.toModels(postgresArgs)
	return transaction, nil
}

func (t *TxService) AddTxInHistory(
	ctx context.Context, 
	input model.TransactionInput,
) (*model.Transaction, error) {
	params := t.ToPostgresParams(input)
	postgTx, err := t.Store.CreateTransaction(ctx, params)
	if err != nil {
	  return nil, err
	}
	transaction := t.toModel(postgTx)
	return transaction, nil
}

