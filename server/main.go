package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/wcygan/ping/server/config"
	"github.com/wcygan/ping/server/handler"
	"github.com/wcygan/ping/server/kafka"
	"github.com/wcygan/ping/server/logger"
	"github.com/wcygan/ping/server/repository"
	"github.com/wcygan/ping/server/service"
	"go.uber.org/zap"

	"buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "Run migrations and exit")
	flag.Parse()

	cfg := config.LoadConfig()
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)

	// Initialize logger
	logger.Init()
	log := logger.Get()

	if *migrateOnly {
		if err := runMigrations(dbURL); err != nil {
			log.Fatal("Failed to run migrations", zap.Error(err))
		}
		log.Info("Migrations completed successfully")
		return
	}

	// Connect to the database
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Failed to parse database config", zap.Error(err))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Failed to create connection pool: %v", zap.Error(err))
	}
	defer pool.Close()

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database: %v", zap.Error(err))
	}

	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	// Create repository with Redis connection
	repo := repository.NewPingRepository(pool, "ping-cache")

	pingService := service.NewPingService(repo, producer, log)
	pingHandler := handler.NewPingServiceHandler(pingService)

	mux := http.NewServeMux()
	mux.Handle(pingv1connect.NewPingServiceHandler(pingHandler))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		err := pool.Ping(r.Context())
		if err != nil {
			log.Error("Health check failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server: %v", zap.Error(err))
	}
}
