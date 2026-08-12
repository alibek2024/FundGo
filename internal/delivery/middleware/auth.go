package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/alibek2024/FundGo/internal/delivery"
	"github.com/alibek2024/FundGo/internal/service/contracts"
)

type authMiddleware struct {
	service contracts.AuthUseCase
}

type contextKey string
const userIDKey contextKey = "userID"

func NewAuthMiddleware(s contracts.AuthUseCase) authMiddleware {
	return authMiddleware{
		service: s,
	}
}

func (a *authMiddleware) CheckAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("access_token")
		if err != nil {
		  	RespondWithError(w, delivery.Unauthorized, errors.New("Authorization cookie is missing."))
		  	return
		}

		claims, err := a.service.Authenticate(tokenString.Value)
		if err != nil {
		  	RespondWithError(w, delivery.Unauthorized, errors.New("Invalid or expired token."))
		  	return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
