package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/popeskul/payment-gateway/internal/api"
	"github.com/popeskul/payment-gateway/internal/auth"
	"github.com/popeskul/payment-gateway/internal/config"
	"github.com/popeskul/payment-gateway/internal/core/services"
	"github.com/popeskul/payment-gateway/internal/hasher"
	"github.com/popeskul/payment-gateway/internal/infrastructure/acquiringbank"
	"github.com/popeskul/payment-gateway/internal/infrastructure/database/postgres"
	"github.com/popeskul/payment-gateway/internal/infrastructure/metrics"
	"github.com/popeskul/payment-gateway/internal/infrastructure/uuid"
	"github.com/popeskul/payment-gateway/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrationsPath := "file://./migrations"
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.DBName)

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		logger.Error("Error creating migrate instance", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("Error applying migrations", "error", err)
		os.Exit(1)
	}

	logger.Info("Migrations applied successfully")

	uuidGenerator := uuid.NewUUIDGenerator()

	merchantRepo := postgres.NewMerchantRepository(db, uuidGenerator)
	paymentRepo := postgres.NewPaymentRepository(db, uuidGenerator)
	refundRepo := postgres.NewRefundRepository(db, uuidGenerator)
	userRepo := postgres.NewUserRepository(db, uuidGenerator)

	tokenStore := postgres.NewPostgresTokenStore(db.Pool)

	acquiringBank := acquiringbank.NewAcquiringBankSimulator(logger, 200*time.Millisecond, 0.05)

	passwordHasher := hasher.NewBcryptPasswordHasher()

	merchantService := services.NewMerchantService(merchantRepo, logger)
	paymentService := services.NewPaymentService(paymentRepo, acquiringBank, logger)
	refundService := services.NewRefundService(refundRepo, paymentRepo, logger)
	userService := services.NewUserService(userRepo, logger, passwordHasher)

	jwtManager := auth.NewJWTManager(&cfg.Auth, tokenStore)

	metrics.InitMetrics()

	router := api.NewRouter(
		services.NewServices(merchantService, paymentService, refundService, userService),
		logger,
		jwtManager,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}
