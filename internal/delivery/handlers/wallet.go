package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/alibek2024/FundGo/internal/delivery/helpers"
	"github.com/alibek2024/FundGo/internal/delivery/middleware"
	"github.com/alibek2024/FundGo/internal/delivery/validation"
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/gorilla/schema"
)

type WalletHandler struct {
	Service contracts.WalletUseCase
	Decoder *schema.Decoder
}

func NewWalletHandler(wallet contracts.WalletUseCase, decoder *schema.Decoder) WalletHandler {
	return WalletHandler{
		Service: wallet,
		Decoder: decoder,
	}
}

func (wh WalletHandler) TopUpBalance(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		helpers.RespondWithError(w, helpers.Unauthorized, helpers.ErrNotFound)
		return
	}
	userID, err := strconv.ParseInt(strUserID, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	var balance dto.BalanceOperationInput

	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	if err := wh.Decoder.Decode(&balance, r.PostForm); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	ok = validation.Validate(w, balance)
	if !ok {
		return
	}
	balance.ID = userID
	wh.Service.TopUpBalance(ctx, balance)

	w.WriteHeader(http.StatusNoContent)
}

func (wh WalletHandler) WithdrawBalance(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		helpers.RespondWithError(w, helpers.Unauthorized, helpers.ErrNotFound)
		return
	}
	userID, err := strconv.ParseInt(strUserID, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	var balance dto.BalanceOperationInput

	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	if err := wh.Decoder.Decode(&balance, r.PostForm); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	ok = validation.Validate(w, balance)
	if !ok {
		return
	}
	balance.ID = userID
	wh.Service.WithdrawBalance(ctx, balance)

	w.WriteHeader(http.StatusNoContent)
}
