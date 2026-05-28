package mapper

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func ToTransactionModel(input generated.Transaction) *model.Transaction {
	return &model.Transaction{
		ID:            input.ID,
		UserID:        input.ID,
		DonationID:    &input.DonationID.Int64,
		Type:          model.TransactionType(input.TransactionType),
		Amount:        input.Amount,
		BalanceBefore: input.BalanceBefore,
		BalanceAfter:  input.BalanceAfter,
		CreatedAt:     input.CreatedAt.Time,
	}
}

func ToTransactionModels(args []generated.Transaction) []*model.Transaction {
	result := make([]*model.Transaction, 0)
	for i, v := range args {
		result[i] = ToTransactionModel(v)
	}
	return result
}

func ToTXPostgresParams(input model.TransactionInput) generated.CreateTransactionParams {
	return generated.CreateTransactionParams{
		UserID:          Int(input.UserID),
		DonationID:      Int(*input.DonationID),
		TransactionType: generated.TransactionType(input.Type),
		Amount:          input.Amount,
		BalanceBefore:   input.BalanceBefore,
		BalanceAfter:    input.BalanceAfter,
	}
}