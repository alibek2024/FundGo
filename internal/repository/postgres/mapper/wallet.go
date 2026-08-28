package mapper

import (
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func BalanceParams(input dto.BalanceOperationInput) generated.AddBalanceParams {
	return generated.AddBalanceParams{
		ID:      input.ID,
		Balance: input.Amount,
	}
}

func SubtractBalanceParams(input dto.BalanceOperationInput) generated.SubtractBalanceParams {
	return generated.SubtractBalanceParams{
		ID:      input.ID,
		Balance: input.Amount,
	}
}
