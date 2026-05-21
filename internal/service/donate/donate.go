package donate

import (
	"context"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository"
)

type DonateService struct {
	Store repository.Store
}

func CreateDonateService(store repository.Store) *DonateService {
	return &DonateService{
		Store: store,
	}
}



func (d *DonateService) createDonate(ctx context.Context, input model.DonateInput) (*model.Donation, error) {
	params := d.toParams(input)
	res, err := d.Store.CreateDonation(ctx, params)
	if err != nil {
		return nil, err
	}
	donation := d.toDonate(res)
	return &donation, nil
}

func (d *DonateService) GetListDonate(ctx context.Context, campaignID int32) ([]*model.Donation, error) {
	rows, err := d.Store.GetListDonations(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	donations := d.toDonationSlice(rows)
	return donations, nil
}

func (d *DonateService) Donate() {
	
}