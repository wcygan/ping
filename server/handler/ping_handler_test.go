package handler

import (
    "context"
    "testing"
    "time"

    pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
    "connectrpc.com/connect"
    "github.com/stretchr/testify/assert"
    "github.com/wcygan/ping/server/errors"
    "github.com/wcygan/ping/server/service"
    "go.uber.org/zap"
)

// FakeService implements the necessary service methods for testing
type FakeService struct {
    pings       []time.Time
    pingCount   int64
    shouldError bool
}

func NewFakeService() *FakeService {
    return &FakeService{
        pings: make([]time.Time, 0),
    }
}

func (f *FakeService) RecordPing(ctx context.Context, timestamp time.Time) error {
    if f.shouldError {
        return &errors.StorageError{Err: fmt.Errorf("fake error")}
    }
    f.pings = append(f.pings, timestamp)
    return nil
}

func (f *FakeService) GetPingCount(ctx context.Context) (int64, error) {
    if f.shouldError {
        return 0, &errors.StorageError{Err: fmt.Errorf("fake error")}
    }
    return f.pingCount, nil
}

func TestPingHandler_Ping(t *testing.T) {
    tests := []struct {
        name        string
        shouldError bool
        wantErr    bool
    }{
        {
            name:        "successful ping",
            shouldError: false,
            wantErr:    false,
        },
        {
            name:        "service error",
            shouldError: true,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            fakeService := NewFakeService()
            fakeService.shouldError = tt.shouldError
            handler := NewPingServiceHandler(fakeService)

            req := connect.NewRequest(&pingv1.PingRequest{
                TimestampMs: time.Now().UnixMilli(),
            })

            resp, err := handler.Ping(context.Background(), req)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Len(t, fakeService.pings, 1)
            }
        })
    }
}

func TestPingHandler_PingCount(t *testing.T) {
    tests := []struct {
        name        string
        pingCount   int64
        shouldError bool
        wantErr    bool
    }{
        {
            name:        "successful count",
            pingCount:   42,
            shouldError: false,
            wantErr:    false,
        },
        {
            name:        "service error",
            pingCount:   0,
            shouldError: true,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            fakeService := NewFakeService()
            fakeService.shouldError = tt.shouldError
            fakeService.pingCount = tt.pingCount
            handler := NewPingServiceHandler(fakeService)

            req := connect.NewRequest(&pingv1.PingCountRequest{})

            resp, err := handler.PingCount(context.Background(), req)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, tt.pingCount, resp.Msg.PingCount)
            }
        })
    }
}
