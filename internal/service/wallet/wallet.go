package wallet

// type Service struct {
// 	Store store.Store
// }

// func NewWalletService(store store.Store) *Service {
// 	return &Service{
// 		Store: store,
// 	}
// }

// func (w *Service) TopUpBalance(ctx context.Context, input dto.BalanceOperationInput) error {
// 	return w.Store.ExecTx(ctx, func(q store.Store) error {
// 		params := mapper.BalanceParams(input)
// 		err := q.AddBalance(ctx, params)
// 		if err != nil {
// 			return err
// 		}

// 		balance, err := q.DB.GetBalance(ctx, input.ID)
// 		if err != nil {
// 			return err
// 		}

// 		txInput := model.TransactionInput{
// 			UserID:        input.ID,
// 			DonationID:    nil,
// 			Type:          model.TransactionTopUp,
// 			Amount:        input.Amount,
// 			BalanceBefore: balance,
// 			BalanceAfter:  balance.Add(input.Amount),
// 		}

// 		_, err = q.DB.CreateTransaction(ctx, mapper.ToTXPostgresParams(txInput))
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// func (w *Service) WithDraw(ctx context.Context, input model.Amount) error {
// 	return w.Store.ExecTx(ctx, func(q store.SQLStore) error {
// 		rows, err := q.DB.SubtractBalance(ctx, mapper.SubtractBalanceParams(input))
// 		if err != nil {
// 			return err
// 		}
// 		if rows == 0 {
// 			return errors.New("insufficient funds")
// 		}

// 		balance, err := q.DB.GetBalance(ctx, input.ID)
// 		if err != nil {
// 			return err
// 		}

// 		txInput := model.TransactionInput{
// 			UserID:        input.ID,
// 			DonationID:    nil,
// 			Type:          model.TransactionWithdraw,
// 			Amount:        input.Amount,
// 			BalanceBefore: balance,
// 			BalanceAfter:  balance.Sub(input.Amount),
// 		}

// 		_, err = q.DB.CreateTransaction(ctx, mapper.ToTXPostgresParams(txInput))
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }
