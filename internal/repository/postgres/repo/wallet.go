package repo

import (
	"context"
	"errors"

	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/postgres/mapper"
	"github.com/alibek2024/FundGo/internal/repository/store"
	"github.com/shopspring/decimal"
)

type WalletStore struct {
	Store store.SQLStore
}

func NewWalletStore(store store.SQLStore) *WalletStore {
	return &WalletStore{
		Store: store,
	}
}

func (w *WalletStore) GetBalance(ctx context.Context, account int64) (*decimal.Decimal, error) {
	balance, err := w.Store.DB.GetBalance(ctx, account)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (w *WalletStore) TopUpBalance(ctx context.Context, input model.Amount) error {
	return w.Store.ExecTx(ctx, func(q store.SQLStore) error {
		params := mapper.BalanceParams(input)
		err := q.DB.AddBalance(ctx, params)
		if err != nil {
			return err
		}

		balance, err := q.DB.GetBalance(ctx, input.ID)
		if err != nil {
			return err
		}

		txInput := model.TransactionInput{
			UserID:        input.ID,
			DonationID:    nil,
			Type:          model.TransactionTopUp,
			Amount:        input.Amount,
			BalanceBefore: balance,
			BalanceAfter:  balance.Add(input.Amount),
		}

		_, err = q.DB.CreateTransaction(ctx, mapper.ToTXPostgresParams(txInput))
		if err != nil {
			return err
		}

		return nil
	})
}

func (w *WalletStore) WithDraw(ctx context.Context, input model.Amount) error {
	return w.Store.ExecTx(ctx, func(q store.SQLStore) error {
		rows, err := q.DB.SubtractBalance(ctx, mapper.SubtractBalanceParams(input))
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("insufficient funds")
		}

		balance, err := q.DB.GetBalance(ctx, input.ID)
		if err != nil {
			return err
		}

		txInput := model.TransactionInput{
			UserID:        input.ID,
			DonationID:    nil,
			Type:          model.TransactionWithdraw,
			Amount:        input.Amount,
			BalanceBefore: balance,
			BalanceAfter:  balance.Sub(input.Amount),
		}

		_, err = q.DB.CreateTransaction(ctx, mapper.ToTXPostgresParams(txInput))
		if err != nil {
			return err
		}

		return nil
	})
}

func (w *WalletStore) DebitUserBalance(ctx context.Context, input model.Amount) error {
	rows, err := w.Store.DB.SubtractBalance(ctx, mapper.SubtractBalanceParams(input))
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("insufficient funds")
	}
	return nil
}
