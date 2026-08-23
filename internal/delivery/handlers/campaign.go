package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alibek2024/FundGo/internal/delivery/helpers"
	"github.com/alibek2024/FundGo/internal/delivery/validation"
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/gorilla/schema"
)

type CampaignHandler struct {
	Service contracts.CampaignUseCase
	decoder *schema.Decoder
}

func NewCampaignHandler(campaign contracts.CampaignUseCase, decoder *schema.Decoder) CampaignHandler {
	return CampaignHandler{
		Service: campaign,
		decoder: decoder,
	}
}

func (c *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var campaignInput dto.CreateCampaignInput

	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	if err := c.decoder.Decode(&campaignInput, r.PostForm); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}
	ok := validation.Validate(w, campaignInput)
	if !ok {
		return
	}

	campaign, err := c.Service.CreateCampaign(ctx, campaignInput)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	helpers.Respond(w, helpers.OK, campaign)
}

func (c *CampaignHandler) SearchCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	query := r.URL.Query().Get("q")
	campaignName := strings.TrimSpace(query)

	if campaignName == "" {
		helpers.RespondWithError(w, helpers.InternalServerError, helpers.ErrQueryEmpty)
		return
	}

	campaign, err := c.Service.SearchCampaign(ctx, campaignName)
	if err != nil {
		if errors.Is(err, contracts.ErrCampaignNotFound) {
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		}
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	helpers.Respond(w, helpers.OK, campaign)
}

func (c *CampaignHandler) WrapUpCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strID := r.PathValue("id")
	campaignID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	err = c.Service.WrapUpCampaign(ctx, campaignID)
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

func (c *CampaignHandler) ForceDeleteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strID := r.PathValue("id")
	campaignID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	err = c.Service.ForceDeleteCampaign(ctx, campaignID)
	if err != nil {
		if errors.Is(err, contracts.ErrCampaignNotFound) {
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		}
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
