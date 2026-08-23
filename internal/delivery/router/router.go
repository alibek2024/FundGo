package router

import (
	"net/http"

	"github.com/alibek2024/FundGo/internal/delivery/handlers"
	"github.com/alibek2024/FundGo/internal/delivery/middleware"
	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router

	userHandler     *handlers.UserHandler
	authHandler     *handlers.AuthHandler
	campaignHandler *handlers.CampaignHandler
	donationHandler *handlers.DonationHandler
	walletHandler   *handlers.WalletHandler

	authMiddleware   middleware.AuthMiddleware
	errorMiddleware  middleware.ErrorMiddleware
	loggerMiddleware middleware.LoggerMiddleware
}

func NewRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	campaignHandler *handlers.CampaignHandler,
	donationHandler *handlers.DonationHandler,
	walletHandler *handlers.WalletHandler,
	authMiddleware *middleware.AuthMiddleware,
	errorMiddleware *middleware.ErrorMiddleware,
	loggerMiddleware *middleware.LoggerMiddleware,
) *Router {
	router := &Router{
		R:               mux.NewRouter(),
		userHandler:     userHandler,
		authHandler:     authHandler,
		campaignHandler: campaignHandler,
		donationHandler: donationHandler,
		walletHandler:   walletHandler,

		authMiddleware:   *authMiddleware,
		errorMiddleware:  *errorMiddleware,
		loggerMiddleware: *loggerMiddleware,
	}

	router.setupRoutes()

	return router
}

func (r *Router) setupRoutes() {
	protected := r.R.PathPrefix("/api/v1").Subrouter()
	protected.Use(r.authMiddleware.CheckAuthorization)
	protected.Use(r.errorMiddleware.ErrorHandlerMiddleware)
	protected.Use(r.loggerMiddleware.LogMiddleware)

	r.setupAuthRoutes()
	r.setupHealthRoutes()

	r.setupUserRoutes(protected)
	r.setupCampaignRoutes(protected)
	r.setupDonationRoutes(protected)
	r.setupWalletRoutes(protected)
}

func (r *Router) setupAuthRoutes() {
	r.R.HandleFunc("/api/v1/auth/register", r.authHandler.Registration).Methods(http.MethodPost)
	r.R.HandleFunc("/api/v1/auth/login", r.authHandler.Authentication).Methods(http.MethodPost)
	r.R.HandleFunc("/api/v1/auth/refresh", r.authHandler.Refresh).Methods(http.MethodPost)
}

func (r *Router) setupUserRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/{id}",
		r.userHandler.GetUserInfo,
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/users/{id}",
		r.userHandler.UpdateInfo,
	).Methods(http.MethodPatch)

	router.HandleFunc(
		"/users/{id}/email",
		r.userHandler.ChangeEmail,
	).Methods(http.MethodPatch)
	router.HandleFunc(
		"/users/{id}/password",
		r.userHandler.ChangePassword,
	).Methods(http.MethodPatch)
	router.HandleFunc(
		"/users/{id}/deactivate",
		r.userHandler.DeactivateAccount,
	).Methods(http.MethodDelete)
	router.HandleFunc(
		"/users/{id}",
		r.userHandler.DeleteAccount,
	).Methods(http.MethodDelete)
}

func (r *Router) setupCampaignRoutes(router *mux.Router) {
	router.HandleFunc(
		"/campaigns",
		r.campaignHandler.CreateCampaign,
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/campaigns/{id}",
		r.campaignHandler.SearchCampaign,
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/campaigns/{id}",
		r.campaignHandler.ForceDeleteCampaign,
	).Methods(http.MethodDelete)
}

func (r *Router) setupDonationRoutes(router *mux.Router) {
	router.HandleFunc(
		"/campaigns/{id}/donate",
		r.donationHandler.DonateToCampaign,
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/donations/{id}/refund",
		r.donationHandler.RefundDonation,
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/users/me/donations",
		r.donationHandler.TransactionHistory,
	).Methods(http.MethodGet)
}

func (r *Router) setupWalletRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/{id}/balance/top-up",
		r.walletHandler.TopUpBalance,
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/users/{id}/balance/withdraw",
		r.walletHandler.WithdrawBalance,
	).Methods(http.MethodPost)
}

func (r *Router) setupHealthRoutes() {
	r.R.HandleFunc(
		"/health",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		},
	).Methods(http.MethodGet)
}
