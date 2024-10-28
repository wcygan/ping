package repository

import (
    "context"
    "fmt"
    "strconv"
    "time"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

type PingRepository struct {
    db    *pgxpool.Pool
    redis *redis.Client
}

func NewPingRepository(db *pgxpool.Pool, redisHost string) *PingRepository {
    rdb := redis.NewClient(&redis.Options{
        Addr: redisHost + ":6379",
    })
    
    return &PingRepository{
        db:    db,
        redis: rdb,
    }
}

func (r *PingRepository) StorePing(ctx context.Context, timestamp time.Time) error {
    _, err := r.db.Exec(ctx, "INSERT INTO pings (pinged_at) VALUES ($1)", timestamp)
    if err != nil {
        return fmt.Errorf("failed to store ping: %v", err)
    }
    return nil
}

func (r *PingRepository) CountPings(ctx context.Context) (int64, error) {
    // Try Redis first
    val, err := r.redis.HGet(ctx, "ping:counters", "total").Result()
    if err == nil {
        // Redis had the count
        count, err := strconv.ParseInt(val, 10, 64)
        if err == nil {
            return count, nil
        }
    }

    // Fallback to database if Redis fails
    var count int64
    err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pings").Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count pings: %v", err)
    }
    return count, nil
}
