package repository

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

// InMemoryDB simulates a database for testing
type InMemoryDB struct {
    pings []time.Time
}

func NewInMemoryDB() *InMemoryDB {
    return &InMemoryDB{
        pings: make([]time.Time, 0),
    }
}

type InMemoryPingRepository struct {
    db *InMemoryDB
}

func NewInMemoryPingRepository() *InMemoryPingRepository {
    return &InMemoryPingRepository{
        db: NewInMemoryDB(),
    }
}

func (r *InMemoryPingRepository) StorePing(ctx context.Context, timestamp time.Time) error {
    r.db.pings = append(r.db.pings, timestamp)
    return nil
}

func (r *InMemoryPingRepository) CountPings(ctx context.Context) (int64, error) {
    return int64(len(r.db.pings)), nil
}

func TestPingRepository_StorePing(t *testing.T) {
    repo := NewInMemoryPingRepository()
    timestamp := time.Now().UTC()

    err := repo.StorePing(context.Background(), timestamp)
    assert.NoError(t, err)

    count, err := repo.CountPings(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, int64(1), count)
    assert.Equal(t, timestamp, repo.db.pings[0])
}

func TestPingRepository_CountPings(t *testing.T) {
    repo := NewInMemoryPingRepository()
    timestamps := []time.Time{
        time.Now().UTC(),
        time.Now().UTC().Add(-1 * time.Hour),
        time.Now().UTC().Add(-2 * time.Hour),
    }

    for _, ts := range timestamps {
        err := repo.StorePing(context.Background(), ts)
        assert.NoError(t, err)
    }

    count, err := repo.CountPings(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, int64(len(timestamps)), count)
}
