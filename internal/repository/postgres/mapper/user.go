package mapper

import (
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/generated"
)

func CreateUserParams(input dto.RegistrationInput) generated.CreateUserParams {
	return generated.CreateUserParams{
		Email:        input.Email,
		PasswordHash: input.HashPassword,
		FirstName:    Text(input.FirstName),
		LastName:     Text(input.LastName),
	}
}

func UpdateUserInfoParams(input dto.UserInfo) generated.UpdateInfoParams {
	return generated.UpdateInfoParams{
		FirstName: Text(input.FirstName),
		LastName:  Text(input.LastName),
		ID:        input.ID,
	}
}
func UpdateUserEmailParams(input dto.UserEmail) generated.UpdateEmailParams {
	return generated.UpdateEmailParams{
		Email: input.Email,
		ID:    input.ID,
	}
}
func UpdateUserPasswordParams(input dto.ChangeUserPassword) generated.UpdatePasswordParams {
	return generated.UpdatePasswordParams{
		PasswordHash: input.HashPassword,
		ID:           input.ID,
	}
}

func ToModel(input generated.User) model.User {
	return model.User{
		ID:           input.ID,
		FirstName:    input.FirstName.String,
		LastName:     input.LastName.String,
		Email:        input.Email,
		HashPassword: input.PasswordHash,
		Balance:      input.Balance,
		CreatedAt:    input.CreatedAt.Time,
		UpdatedAt:    input.UpdatedAt.Time,
		DeletedAt:    input.DeletedAt.Time,
	}
}
