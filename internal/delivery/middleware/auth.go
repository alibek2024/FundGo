package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/alibek2024/FundGo/internal/delivery/helpers"
	"github.com/alibek2024/FundGo/internal/service/contracts"
)

type AuthMiddleware struct {
	service contracts.AuthUseCase
}

type contextKey string

const UserIDKey contextKey = "userID"

func NewAuthMiddleware(s contracts.AuthUseCase) AuthMiddleware {
	return AuthMiddleware{
		service: s,
	}
}

func (a *AuthMiddleware) CheckAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("access_token")
		if err != nil {
			helpers.RespondWithError(w, helpers.Unauthorized, errors.New("Authorization cookie is missing."))
			return
		}

		claims, err := a.service.Authenticate(tokenString.Value)
		if err != nil {
			helpers.RespondWithError(w, helpers.Unauthorized, errors.New("Invalid or expired token."))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
