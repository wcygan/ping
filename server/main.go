package main

import (
	"buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)

// PingServiceServer implements the PingService interface.
type PingServiceServer struct {
	db *pgxpool.Pool
}

// Ping handles the Ping RPC.
func (s *PingServiceServer) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
	timestamp := time.Unix(0, req.Msg.TimestampMs*int64(time.Millisecond)).UTC()
	log.Printf("Received a ping at %s (UTC)", timestamp.Format(time.RFC3339))

	_, err := s.db.Exec(ctx, "INSERT INTO pings (pinged_at) VALUES ($1)", timestamp)
	if err != nil {
		log.Printf("Failed to insert ping: %v", err)
		return nil, fmt.Errorf("failed to store ping: %v", err)
	}

	return connect.NewResponse(&pingv1.PingResponse{}), nil
}

// PingCount handles the PingCount RPC.
func (s *PingServiceServer) PingCount(ctx context.Context, req *connect.Request[pingv1.PingCountRequest]) (*connect.Response[pingv1.PingCountResponse], error) {
	var count int64
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM pings").Scan(&count)
	if err != nil {
		log.Printf("Failed to count pings: %v", err)
		return nil, fmt.Errorf("failed to count pings: %v", err)
	}

	return connect.NewResponse(&pingv1.PingCountResponse{PingCount: count}), nil
}

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

func main() {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "pinguser"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "ping123"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "pingdb"
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbName)

	// Run migrations
	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Connect to the database
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	server := &PingServiceServer{db: pool}

	// Register the PingService with the Connect server.
	mux.Handle(pingv1connect.NewPingServiceHandler(server))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
