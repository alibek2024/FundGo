package mapper

import (
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

func DefaultUserParams(
	input model.UserInput,
	hashPassword string,
) (
	string, string,
	pgtype.Text, pgtype.Text,
) {
	firstname := Text(input.FirstName)
	lastname := Text(input.LastName)

	return input.Email,
		hashPassword,
		firstname,
		lastname
}

func UserParams(input model.UserInput, hashPassword string) generated.CreateUserParams {
	email, hash, first, last := DefaultUserParams(input, hashPassword)
	return generated.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		FirstName:    first,
		LastName:     last,
	}
}

func UpdateUserParams(input model.UserInput, hashPassword string) generated.UpdateUserParams {
	email, hash, first, last := DefaultUserParams(input, hashPassword)
	return generated.UpdateUserParams{
		Email:        email,
		PasswordHash: hash,
		FirstName:    first,
		LastName:     last,
		ID:           input.ID,
	}
}

func UserResponse(input generated.User) model.UserResponse {
	return model.UserResponse{
		FirstName: input.FirstName.String,
		LastName:  input.LastName.String,
		Email:     input.Email,
		ID:        input.ID,
		Balance:   input.Balance,
		CreatedAt: input.CreatedAt.Time,
		DeletedAt: input.DeletedAt.Time,
	}
}
