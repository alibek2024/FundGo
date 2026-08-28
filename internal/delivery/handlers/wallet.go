package handlers

import (
	"context"
	"encoding/json"
	"errors"
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

	if err := json.NewDecoder(r.Body).Decode(&balance); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	balance.ID = userID

	ok = validation.Validate(w, balance)
	if !ok {
		return
	}
	if err := wh.Service.TopUpBalance(ctx, balance); err != nil {
		if errors.Is(err, contracts.ErrUserNotFound) {
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		}
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

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

	if err := json.NewDecoder(r.Body).Decode(&balance); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	balance.ID = userID

	ok = validation.Validate(w, balance)
	if !ok {
		return
	}
	if err := wh.Service.WithdrawBalance(ctx, balance); err != nil {
		if errors.Is(err, contracts.ErrUserNotFound) || errors.Is(err, contracts.ErrInsufficientBalance) {
			helpers.RespondWithError(w, helpers.BadRequest, err)
			return
		}
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
