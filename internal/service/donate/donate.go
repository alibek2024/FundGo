package donate

// var ErrCampaignInactive = errors.New("campaign inactive")
// var ErrInsufficientFunds = errors.New("insufficient funds")

// type Service struct {
// 	Store store.Store
// }

// func CreateDonateService(store store.Store) *Service {
// 	return &Service{
// 		Store: store,
// 	}
// }

// func (d *Service) DonateToCampaign(ctx context.Context, input dto.DonateInput) error {
// 	return d.Store.ExecTx(ctx, func(q store.SQLStore) error {
// 		status, err := q.DB.GetCampaignStatus(ctx, input.CampaignID)
// 		if err != nil {
// 			return err
// 		}
// 		if !mapper.CheckStatusCampaign(status) {
// 			return ErrCampaignInactive
// 		}

// 		rows, err := q.DB.SubtractBalance(ctx, generated.SubtractBalanceParams{
// 			ID:      input.UserID,
// 			Balance: input.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		if rows == 0 {
// 			return ErrInsufficientFunds
// 		}

// 		_, err = q.DB.IncreaseCampaignAmount(ctx, generated.IncreaseCampaignAmountParams{
// 			ID:            input.CampaignID,
// 			CurrentAmount: input.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		err = d.CreateDonation(ctx, q.DB, generated.CreateDonationParams{
// 			UserID:     mapper.Int(input.UserID),
// 			CampaignID: input.CampaignID,
// 			Amount:     input.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }
