package repo

import (
	"context"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/alibek2024/FundGo/internal/repository/store"
)

type TransactionStore struct {
	Store store.SQLStore
}

func CreateTXStore(db store.SQLStore) *TransactionStore {
	return &TransactionStore{
		Store: db,
	}
}

func (t *TransactionStore) TxHistory(
	ctx context.Context,
	userID int64,
) ([]*model.Transaction, error) {
	params := mapper.Int(userID)
	postgresArgs, err := t.Store.DB.HistoryTX(ctx, params)
	if err != nil {
		return nil, err
	}
	transaction := mapper.ToTransactionModels(postgresArgs)
	return transaction, nil
}

func (t *TransactionStore) AddTxInHistory(
	ctx context.Context,
	input model.TransactionInput,
) (*model.Transaction, error) {
	params := mapper.ToTXPostgresParams(input)
	postgTx, err := t.Store.DB.CreateTransaction(ctx, params)
	if err != nil {
		return nil, err
	}
	transaction := mapper.ToTransactionModel(postgTx)
	return transaction, nil
}
