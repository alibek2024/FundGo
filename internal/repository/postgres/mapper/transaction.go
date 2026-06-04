package mapper

import (
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func ToTransactionModel(input generated.Transaction) *model.Transaction {
	if input.DonationID.Valid == false {
		return &model.Transaction{
			ID:            input.ID,
			UserID:        input.UserID.Int64,
			DonationID:    nil,
			Type:          model.TransactionType(input.TransactionType),
			Amount:        input.Amount,
			BalanceBefore: input.BalanceBefore,
			BalanceAfter:  input.BalanceAfter,
			CreatedAt:     input.CreatedAt.Time,
		}
	}
	return &model.Transaction{
		ID:            input.ID,
		UserID:        input.UserID.Int64,
		DonationID:    &input.DonationID.Int64,
		Type:          model.TransactionType(input.TransactionType),
		Amount:        input.Amount,
		BalanceBefore: input.BalanceBefore,
		BalanceAfter:  input.BalanceAfter,
		CreatedAt:     input.CreatedAt.Time,
	}
}

func ToTransactionModels(args []generated.Transaction) []model.Transaction {
	result := make([]model.Transaction, len(args))
	for i, v := range args {
		result[i] = *ToTransactionModel(v)
	}
	return result
}

func ToTXPostgresParams(input dto.TransactionInput) generated.CreateTransactionParams {
	if input.DonationID != nil {
		return generated.CreateTransactionParams{
			UserID:          Int(input.UserID),
			DonationID:      Int(*input.DonationID),
			TransactionType: generated.TransactionTypeDonation,
			Amount:          input.Amount,
			BalanceBefore:   input.BalanceBefore,
			BalanceAfter:    input.BalanceAfter,
		}
	}
	return generated.CreateTransactionParams{
		UserID:          Int(input.UserID),
		TransactionType: generated.TransactionType(input.Type),
		Amount:          input.Amount,
		BalanceBefore:   input.BalanceBefore,
		BalanceAfter:    input.BalanceAfter,
	}
}
