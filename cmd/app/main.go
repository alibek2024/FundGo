package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alibek2024/FundGo/internal/config"
	"github.com/alibek2024/FundGo/internal/delivery"
	"github.com/alibek2024/FundGo/internal/delivery/middleware"
	"github.com/alibek2024/FundGo/internal/delivery/router"
	"github.com/alibek2024/FundGo/internal/delivery/worker"
	"github.com/alibek2024/FundGo/internal/repository/postgres"
	"github.com/alibek2024/FundGo/internal/service"
	"github.com/alibek2024/FundGo/internal/service/campaign"
	"github.com/gorilla/schema"
	"github.com/hibiken/asynq"
)

func main() {
	var decoder = schema.NewDecoder()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config parse error: %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	pgpool, err := postgres.InitDB(ctx, cfg.DB.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgpool.Close()

	redisURL, err := url.Parse(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("parse redis URL: %v", err)
	}

	redisConn := asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    redisURL.Host,
	}

	if redisURL.User != nil {
		redisConn.Username = redisURL.User.Username()

		password, ok := redisURL.User.Password()
		if ok {
			redisConn.Password = password
		}
	}

	if cfg.Redis.TLS {
		redisConn.TLSConfig = &tls.Config{}
	}
	asynqClient := asynq.NewClient(redisConn)
	defer asynqClient.Close()

	repo := postgres.NewStore(pgpool)
	service, err := service.NewService(repo, cfg.JWT, asynqClient)
	if err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}

	asynqServer := asynq.NewServer(
		redisConn,
		asynq.Config{
			Concurrency: 10,
		},
	)

	campaignWorkerHandler := worker.NewCampaignTaskHandler(service.Campaign)
	asynqMux := asynq.NewServeMux()
	asynqMux.HandleFunc(campaign.TypeCloseCampaign, campaignWorkerHandler.ProcessTask)

	go func() {
		if err := asynqServer.Run(asynqMux); err != nil {
			log.Printf("asynq server error: %v", err)
		}
	}()

	delivery := delivery.NewDelivery(&service, *decoder)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	authMiddleware := middleware.NewAuthMiddleware(&service.Auth)
	errorMiddleware := middleware.NewErrorMiddleware(logger)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)

	router := router.NewRouter(
		&delivery.UserHandler,
		&delivery.AuthHandler,
		&delivery.CampaignHandler,
		&delivery.DonationHandler,
		&delivery.WalletHandler,

		&authMiddleware,
		&errorMiddleware,
		&loggerMiddleware,
	)

	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router.R,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	asynqServer.Shutdown()

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("Server stopped successfully")
}
