package service

import (
    "context"
    "time"
)

type PingRepository interface {
    StorePing(ctx context.Context, timestamp time.Time) error
    CountPings(ctx context.Context) (int64, error)
}

type EventProducer interface {
    SendPingEvent(ctx context.Context, timestamp time.Time) error
    Close() error
}
