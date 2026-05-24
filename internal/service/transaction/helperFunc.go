package transaction

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/helperfunc"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
)

func (t *TxService) toModel(input postgres.Transaction) *model.Transaction {
	return &model.Transaction{
		ID:         input.ID,
		UserID:     input.ID,
		DonationID: &input.DonationID.Int64,
		Type:       model.TransactionType(input.TransactionType),
		Amount:     input.Amount,
		CreatedAt:  input.CreatedAt.Time,
	}
}

func (t *TxService) toModels(args []postgres.Transaction) []*model.Transaction {
	result := make([]*model.Transaction, 0)
	for i, v := range args {
		result[i] = t.toModel(v)
	}
	return result
}

func (t *TxService) ToPostgresParams(input model.TransactionInput) postgres.CreateTransactionParams {
	return postgres.CreateTransactionParams{
		UserID:          helperfunc.Int(input.UserID),
		DonationID:      helperfunc.Int(*input.DonationID),
		TransactionType: postgres.TransactionType(input.Type),
		Amount:          input.Amount,
	}
}
