package postgres

import (
	"context"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
)

func (u *Repository) CreateUser(
	ctx context.Context,
	input dto.RegistrationInput,
) (*model.User, error) {
	params := mapper.CreateUserParams(input)

	storeUser, err := u.DB.CreateUser(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}

	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	storeUser, err := u.DB.GetByEmail(ctx, email)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	storeUser, err := u.DB.GetByID(ctx, id)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}

	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) UpdateInfo(
	ctx context.Context,
	input dto.UserInfo,
) (*model.User, error) {
	params := mapper.UpdateUserInfoParams(input)
	storeUser, err := u.DB.UpdateInfo(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) UpdateEmail(
	ctx context.Context,
	input dto.UserEmail,
) (*model.User, error) {
	params := mapper.UpdateUserEmailParams(input)
	storeUser, err := u.DB.UpdateEmail(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) UpdatePassword(
	ctx context.Context,
	input dto.UpdateUserPassword,
) (*model.User, error) {
	params := mapper.UpdateUserPasswordParams(input)
	storeUser, err := u.DB.UpdatePassword(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) SoftDeleteUser(ctx context.Context, id int64) error {
	if err := u.DB.SoftDeleteUser(ctx, id); err != nil {
		return mapper.MapDBError(err)
	}
	return nil
}

func (u *Repository) DeleteUser(ctx context.Context, id int64) error {
	if err := u.DB.DeleteUser(ctx, id); err != nil {
		return mapper.MapDBError(err)
	}
	return nil
}

func (u *Repository) RestoreUser(ctx context.Context, id int64) (*model.User, error) {
	storeUser, err := u.DB.RestoreUser(ctx, id)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	user := mapper.ToModel(storeUser)
	return &user, nil
}

func (u *Repository) GetByIDForPurge(ctx context.Context, id int64) (*model.User, error) {
	storeUser, err := u.DB.GetByIDForPurge(ctx, id)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	user := mapper.ToModel(storeUser)

	return &user, nil
}
