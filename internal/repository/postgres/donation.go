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

func (d *Repository) GetDonationByID(ctx context.Context, id int64) (*model.Donation, error) {
	storeDonate, err := d.DB.GetDonationByID(ctx, id)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	donation := mapper.ToDonationModel(storeDonate)
	return &donation, nil
}

func (d *Repository) UpdateDonationStatus(
	ctx context.Context,
	input dto.UpdateDonationStatus,
) (*model.Donation, error) {
	params := mapper.ToDonationUpdateStatus(input)
	storeDonate, err := d.DB.UpdateDonationStatus(ctx, params)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	donation := mapper.ToDonationModel(storeDonate)
	return &donation, nil
}

func (d *Repository) RefundDonationStatus(
	ctx context.Context,
	donationID int64,
) (*model.Donation, error) {
	storeDonate, err := d.DB.DonationRefunded(ctx, donationID)
	if err != nil {
		return nil, mapper.MapDBError(err)
	}
	donation := mapper.ToDonationModel(storeDonate)
	return &donation, nil
}
