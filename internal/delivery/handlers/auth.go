package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/alibek2024/FundGo/internal/delivery/helpers"
	"github.com/alibek2024/FundGo/internal/delivery/validation"
	"github.com/alibek2024/FundGo/internal/dto"
	"github.com/alibek2024/FundGo/internal/service/contracts"
	"github.com/gorilla/schema"
)

type AuthHandler struct {
	auth    contracts.AuthUseCase
	Decoder *schema.Decoder
}

func NewAuthHandler(auth contracts.AuthUseCase, decoder *schema.Decoder) AuthHandler {
	return AuthHandler{
		auth:    auth,
		Decoder: decoder,
	}
}

func (a *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var regInput dto.RegistrationInput

	if err := json.NewDecoder(r.Body).Decode(&regInput); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	if !validation.Validate(w, regInput) {
		return
	}

	user, tokens, err := a.auth.RegisterUser(ctx, &regInput)
	if err != nil {
		if errors.Is(err, contracts.ErrEmailAlreadyExists) {
			helpers.RespondWithError(w, helpers.BadRequest, err)
			return
		}

		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	a.setAuthCookies(w, tokens)
	helpers.Respond(w, helpers.OK, user)
}

func (a *AuthHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var signInput dto.SignIn

	if err := json.NewDecoder(r.Body).Decode(&signInput); err != nil {
		helpers.RespondWithError(w, helpers.BadRequest, err)
		return
	}

	if !validation.Validate(w, signInput) {
		return
	}

	user, tokens, err := a.auth.SignIn(ctx, signInput)
	if err != nil {
		if errors.Is(err, contracts.ErrLogin) {
			helpers.RespondWithError(w, helpers.Unauthorized, err)
			return
		}

		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	a.setAuthCookies(w, tokens)
	helpers.Respond(w, helpers.OK, user)
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		helpers.RespondWithError(w, helpers.Unauthorized, helpers.ErrAuth)
		return
	}

	refreshToken := cookie.Value

	claims, err := a.auth.Authenticate(refreshToken)
	if err != nil {
		a.clearAuthCookies(w)
		helpers.RespondWithError(
			w,
			helpers.Unauthorized,
			errors.New("invalid or expired token"),
		)
		return
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	tokens, err := a.auth.GetAccessToken(userID)
	if err != nil {
		helpers.RespondWithError(w, helpers.InternalServerError, err)
		return
	}

	a.setAccessCookies(w, tokens)
	w.WriteHeader(http.StatusNoContent)
}

func (a *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (a *AuthHandler) setAuthCookies(w http.ResponseWriter, tokens *dto.AuthTokens) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/v1/auth/refresh",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (a *AuthHandler) setAccessCookies(w http.ResponseWriter, tokens *dto.AuthTokens) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
