package handlers_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alibek2024/FundGo/internal/delivery/middleware"
	"github.com/alibek2024/FundGo/internal/model"
	"github.com/alibek2024/FundGo/internal/repository/store"
	authsvc "github.com/alibek2024/FundGo/internal/service/auth"
	campaignsvc "github.com/alibek2024/FundGo/internal/service/campaign"
	donatesvc "github.com/alibek2024/FundGo/internal/service/donate"
	transactionsvc "github.com/alibek2024/FundGo/internal/service/transaction"
	usersvc "github.com/alibek2024/FundGo/internal/service/user"
	walletsvc "github.com/alibek2024/FundGo/internal/service/wallet"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func jsonRequest(t *testing.T, method, urlStr string, body any) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}

	req := httptest.NewRequest(method, urlStr, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func authedRequest(req *http.Request, userID string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func expectExecTx(mockStore *store.MockStore) {
	mockStore.EXPECT().
		ExecTx(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(store.Store) error) error {
			return fn(mockStore)
		})
}

func newAuthService(t *testing.T, mockStore *store.MockUserStore) *authsvc.Service {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	return authsvc.NewUserService(
		mockStore,
		time.Minute,
		time.Hour,
		privateKey,
		&privateKey.PublicKey,
	)
}

func newUserService(mockStore store.UserStore) *usersvc.Service {
	return usersvc.NewUserService(mockStore)
}

func newCampaignService(mockStore store.CampaignStore, refundManager campaignsvc.RefundManager) *campaignsvc.Service {
	return campaignsvc.NewCampaignService(mockStore, refundManager)
}

func newDonateService(mockStore store.TransactionManager) *donatesvc.Service {
	return donatesvc.CreateDonateService(mockStore)
}

func newTransactionService(mockStore store.Store) *transactionsvc.Service {
	return transactionsvc.CreateTX(mockStore)
}

func newWalletService(mockStore store.TransactionManager) *walletsvc.Service {
	return walletsvc.NewWalletService(mockStore)
}

func zeroDeletedAtUser(id int64) *model.User {
	deletedAt := time.Time{}
	return &model.User{
		ID:        id,
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
		DeletedAt: &deletedAt,
	}
}

func findCookie(cookies []*http.Cookie, name string) string {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func assertJSONError(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	require.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	require.Contains(t, rec.Body.String(), "error")
}
