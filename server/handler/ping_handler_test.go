package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/wcygan/ping/server/service"
)

// MockPingService implements service.PingService for testing
type MockPingService struct {
	shouldError bool
}

func NewMockPingService(shouldError bool) *service.PingService {
	mock := &MockPingService{
		shouldError: shouldError,
	}
	return &service.PingService{} // Return an empty service for the handler
}

func (m *MockPingService) RecordPing(ctx context.Context, timestamp time.Time) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockPingService) GetPingCount(ctx context.Context) (int64, error) {
	if m.shouldError {
		return 0, fmt.Errorf("mock error")
	}
	return 42, nil
}

func TestPingHandler_Ping(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		wantErr     bool
	}{
		{
			name:        "successful ping",
			shouldError: false,
			wantErr:     false,
		},
		{
			name:        "service error",
			shouldError: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockPingService(tt.shouldError)
			handler := NewPingServiceHandler(mockService)

			req := connect.NewRequest(&pingv1.PingRequest{
				TimestampMs: time.Now().UnixMilli(),
			})

			resp, err := handler.Ping(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestPingHandler_PingCount(t *testing.T) {
	tests := []struct {
		name        string
		pingCount   int64
		shouldError bool
		wantErr     bool
	}{
		{
			name:        "successful count",
			pingCount:   42,
			shouldError: false,
			wantErr:     false,
		},
		{
			name:        "service error",
			pingCount:   0,
			shouldError: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockPingService(tt.shouldError)
			handler := NewPingServiceHandler(mockService)

			req := connect.NewRequest(&pingv1.PingCountRequest{})

			resp, err := handler.PingCount(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
