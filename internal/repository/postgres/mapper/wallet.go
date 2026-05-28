package mapper

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func BalanceParams(input model.Amount) generated.AddBalanceParams {
	return generated.AddBalanceParams{
		ID: input.ID,
		Balance: input.Amount,
	}
}

func SubtractBalanceParams(input model.Amount) generated.SubtractBalanceParams {
	return generated.SubtractBalanceParams{
		ID: input.ID,
		Balance: input.Amount,
	}
}