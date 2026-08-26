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
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type DonationHandler struct {
	Service contracts.DonationUseCase
	tx      contracts.TransactionUseCase
	Decoder *schema.Decoder
}

func NewDonationHandler(donate contracts.DonationUseCase, decoder *schema.Decoder, tx contracts.TransactionUseCase) DonationHandler {
	return DonationHandler{
		Service: donate,
		tx:      tx,
		Decoder: decoder,
	}
}

func (d DonationHandler) DonateToCampaign(w http.ResponseWriter, r *http.Request) {
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

	var dotnationInput dto.DonateInput

	if err := json.NewDecoder(r.Body).Decode(&dotnationInput); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	dotnationInput.UserID = userID

	ok = validation.Validate(w, dotnationInput)
	if !ok {
		return
	}

	err = d.Service.DonateToCampaign(ctx, dotnationInput)
	if err != nil {
		switch {
		case errors.Is(err, contracts.ErrCampaignNotFound):
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		case errors.Is(err, contracts.ErrCampaignNotActive):
			helpers.RespondWithError(w, helpers.BadRequest, helpers.ErrCampaignNotActive)
			return
		default:
			helpers.RespondWithError(w, helpers.InternalServerError, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (d DonationHandler) TransactionHistory(w http.ResponseWriter, r *http.Request) {
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

	transactions, err := d.tx.GetPaymentHistory(ctx, userID)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	helpers.Respond(w, helpers.OK, transactions)
}

func (d DonationHandler) RefundDonation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	strID := vars["id"]
	if strID == "" {
		strID = vars["donation_id"]
	}
	donationID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

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

	err = d.tx.CheckDonation(ctx, userID, donationID)
	if err != nil {
		if errors.Is(err, contracts.ErrDonationID) {
			helpers.RespondWithError(w, helpers.Unauthorized, contracts.ErrDonationID)
			return
		}
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	err = d.Service.RefundDonation(ctx, donationID)
	if err != nil {
		switch {
		case errors.Is(err, contracts.ErrCampaignNotFound):
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		case errors.Is(err, contracts.ErrCampaignNotActive):
			helpers.RespondWithError(w, helpers.BadRequest, helpers.ErrCampaignNotActive)
			return
		case errors.Is(err, contracts.ErrDonationNotFound):
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		case errors.Is(err, contracts.ErrDonateRefunded):
			helpers.RespondWithError(w, helpers.BadRequest, err)
			return
		case errors.Is(err, contracts.ErrCampaignClosed):
			helpers.RespondWithError(w, helpers.BadRequest, err)
			return
		case errors.Is(err, contracts.ErrUserNotFound):
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		default:
			helpers.RespondWithError(w, helpers.InternalServerError, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
