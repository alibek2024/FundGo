package postgres

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
)

var ErrCampaignInactive = errors.New("campaign inactive")
var ErrInsufficientFunds = errors.New("insufficient funds")

func (d *Repository) CreateDonation(ctx context.Context, input dto.DonateInput) (*model.Donation, error) {
	params := mapper.ToSqlcModel(input)
	storeDonate, err := d.DB.CreateDonation(ctx, params)
	if err != nil {
		return nil, err
	}
	donate := mapper.ToDonationModel(storeDonate)
	return &donate, nil
}

func (d *Repository) GetListDonations(ctx context.Context, campaignID int64) ([]model.Donation, error) {
	if campaignID < 0 {
		return nil, errors.New("id is negative")
	}
	storeDonationList, err := d.DB.GetListDonations(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	DonationList := mapper.ToDonationsModels(storeDonationList)
	return DonationList, nil
}
