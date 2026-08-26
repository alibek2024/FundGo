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

type UserHandler struct {
	Service contracts.UserUseCase
	Decoder *schema.Decoder
}

func NewUserHandler(user contracts.UserUseCase, decoder *schema.Decoder) UserHandler {
	return UserHandler{
		Service: user,
		Decoder: decoder,
	}
}

func (u *UserHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
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
	var updateInput dto.UserInfo

	if err := json.NewDecoder(r.Body).Decode(&updateInput); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	updateInput.ID = userID
	ok = validation.Validate(w, updateInput)
	if !ok {
		return
	}

	err = u.Service.UpdateUserInfo(ctx, &updateInput)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
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

	user, err := u.Service.UserInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, contracts.ErrUserNotFound) {
			helpers.RespondWithError(w, helpers.NotFound, err)
			return
		}
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	helpers.Respond(w, helpers.OK, user)
}

func (u *UserHandler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
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
	var updateInput dto.UserEmail

	if err := json.NewDecoder(r.Body).Decode(&updateInput); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	updateInput.ID = userID
	ok = validation.Validate(w, updateInput)
	if !ok {
		return
	}

	err = u.Service.ChangeEmail(ctx, &updateInput)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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
	var updateInput dto.ChangeUserPassword

	if err := json.NewDecoder(r.Body).Decode(&updateInput); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	updateInput.ID = userID
	ok = validation.Validate(w, updateInput)
	if !ok {
		return
	}

	err = u.Service.ChangePassword(ctx, &updateInput)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
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

	err = u.Service.PurgeUserData(ctx, userID)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHandler) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
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

	err = u.Service.DeactivateAccount(ctx, userID)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
