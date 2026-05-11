package user

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgtype"
)

func UserParams(input model.UserInput, hashPassword string) postgres.CreateUserParams {
	return postgres.CreateUserParams{
		Email:        input.Email,
		PasswordHash: hashPassword,
		FirstName: pgtype.Text{
			String: input.FirstName,
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: input.LastName,
			Valid:  true,
		},
	}
}

func UserResponse(input postgres.User) model.UserResponse {
	return model.UserResponse{
		FirstName: input.FirstName.String,
		LastName:  input.LastName.String,
		ID:        input.ID,
		Balance:   input.Balance,
		CreatedAt: input.CreatedAt.Time,
		DeletedAt: input.DeletedAt.Time,
	}
}
