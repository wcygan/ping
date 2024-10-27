package service

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"
)

// FakePingRepository implements PingRepository interface for testing
type FakePingRepository struct {
    pings    []time.Time
    failNext bool
}

func NewFakePingRepository() *FakePingRepository {
    return &FakePingRepository{
        pings: make([]time.Time, 0),
    }
}

func (f *FakePingRepository) StorePing(ctx context.Context, timestamp time.Time) error {
    if f.failNext {
        return &errors.StorageError{Err: fmt.Errorf("fake storage error")}
    }
    f.pings = append(f.pings, timestamp)
    return nil
}

func (f *FakePingRepository) CountPings(ctx context.Context) (int64, error) {
    if f.failNext {
        return 0, &errors.StorageError{Err: fmt.Errorf("fake storage error")}
    }
    return int64(len(f.pings)), nil
}

// FakeEventProducer implements EventProducer interface for testing
type FakeEventProducer struct {
    events   []time.Time
    failNext bool
}

func NewFakeEventProducer() *FakeEventProducer {
    return &FakeEventProducer{
        events: make([]time.Time, 0),
    }
}

func (f *FakeEventProducer) SendPingEvent(ctx context.Context, timestamp time.Time) error {
    if f.failNext {
        return &errors.KafkaError{Err: fmt.Errorf("fake kafka error")}
    }
    f.events = append(f.events, timestamp)
    return nil
}

func (f *FakeEventProducer) Close() error {
    return nil
}

func TestPingService_RecordPing(t *testing.T) {
    tests := []struct {
        name          string
        repoFail     bool
        producerFail bool
        wantErr      bool
        errType      interface{}
    }{
        {
            name:          "successful ping",
            repoFail:     false,
            producerFail: false,
            wantErr:      false,
        },
        {
            name:          "repository failure",
            repoFail:     true,
            producerFail: false,
            wantErr:      true,
            errType:      &errors.StorageError{},
        },
        {
            name:          "producer failure",
            repoFail:     false,
            producerFail: true,
            wantErr:      true,
            errType:      &errors.KafkaError{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := NewFakePingRepository()
            producer := NewFakeEventProducer()
            logger := zap.NewNop()

            repo.failNext = tt.repoFail
            producer.failNext = tt.producerFail

            service := NewPingService(repo, producer, logger)
            timestamp := time.Now().UTC()

            err := service.RecordPing(context.Background(), timestamp)

            if tt.wantErr {
                assert.Error(t, err)
                assert.IsType(t, tt.errType, err)
            } else {
                assert.NoError(t, err)
                assert.Len(t, repo.pings, 1)
                assert.Equal(t, timestamp, repo.pings[0])
                assert.Len(t, producer.events, 1)
                assert.Equal(t, timestamp, producer.events[0])
            }
        })
    }
}

func TestPingService_GetPingCount(t *testing.T) {
    tests := []struct {
        name      string
        pingCount int
        repoFail  bool
        wantErr   bool
    }{
        {
            name:      "successful count",
            pingCount: 3,
            repoFail:  false,
            wantErr:   false,
        },
        {
            name:      "repository failure",
            pingCount: 0,
            repoFail:  true,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := NewFakePingRepository()
            producer := NewFakeEventProducer()
            logger := zap.NewNop()

            // Add test pings
            for i := 0; i < tt.pingCount; i++ {
                repo.pings = append(repo.pings, time.Now().UTC())
            }

            repo.failNext = tt.repoFail
            service := NewPingService(repo, producer, logger)

            count, err := service.GetPingCount(context.Background())

            if tt.wantErr {
                assert.Error(t, err)
                assert.IsType(t, &errors.StorageError{}, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, int64(tt.pingCount), count)
            }
        })
    }
}
