package repository

import (
    "context"
    "fmt"
    "time"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PingRepository struct {
    db *pgxpool.Pool
}

func NewPingRepository(db *pgxpool.Pool) *PingRepository {
    return &PingRepository{db: db}
}

func (r *PingRepository) StorePing(ctx context.Context, timestamp time.Time) error {
    _, err := r.db.Exec(ctx, "INSERT INTO pings (pinged_at) VALUES ($1)", timestamp)
    if err != nil {
        return fmt.Errorf("failed to store ping: %v", err)
    }
    return nil
}

func (r *PingRepository) CountPings(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pings").Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count pings: %v", err)
    }
    return count, nil
}
