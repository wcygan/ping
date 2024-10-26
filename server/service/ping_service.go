package service

import (
    "context"
    "fmt"
    "time"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/wcygan/ping/server/kafka"
)

type PingService struct {
    db       *pgxpool.Pool
    producer *kafka.Producer
}

func NewPingService(db *pgxpool.Pool, producer *kafka.Producer) *PingService {
    return &PingService{
        db:       db,
        producer: producer,
    }
}

func (s *PingService) RecordPing(ctx context.Context, timestamp time.Time) error {
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback(ctx)

    if _, err := tx.Exec(ctx, "INSERT INTO pings (pinged_at) VALUES ($1)", timestamp); err != nil {
        return fmt.Errorf("failed to store ping: %v", err)
    }

    if err := s.producer.SendPingEvent(ctx, timestamp); err != nil {
        return fmt.Errorf("failed to send to Kafka: %v", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    return nil
}

func (s *PingService) GetPingCount(ctx context.Context) (int64, error) {
    var count int64
    err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM pings").Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count pings: %v", err)
    }
    return count, nil
}
