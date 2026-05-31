package postgres

import (
	"context"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
)

func (d *Repository) CreateDonation(ctx context.Context, input dto.DonateInput) (*model.Donation, error) {
	params := mapper.ToSqlcModel(input)
	storeDonate, err := d.DB.CreateDonation(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	donate := mapper.ToDonationModel(storeDonate)
	return &donate, nil
}

func (d *Repository) GetListDonations(ctx context.Context, campaignID int64) ([]model.Donation, error) {
	storeDonationList, err := d.DB.GetListDonations(ctx, campaignID)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	DonationList := mapper.ToDonationsModels(storeDonationList)
	return DonationList, nil
}
