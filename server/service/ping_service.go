package service

import (
	"context"
	"github.com/wcygan/ping/server/errors"
	"go.uber.org/zap"
	"time"
)

type PingService struct {
	repo     PingRepository
	producer EventProducer
	logger   *zap.Logger
}

func NewPingService(repo PingRepository, producer EventProducer, logger *zap.Logger) *PingService {
	return &PingService{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

func (s *PingService) RecordPing(ctx context.Context, timestamp time.Time) error {
	if err := s.repo.StorePing(ctx, timestamp); err != nil {
		s.logger.Error("failed to store ping", zap.Error(err))
		return &errors.StorageError{Err: err}
	}

	if err := s.producer.SendPingEvent(ctx, timestamp); err != nil {
		s.logger.Error("failed to send event", zap.Error(err))
		return &errors.KafkaError{Err: err}
	}

	s.logger.Info("ping recorded successfully",
		zap.Time("timestamp", timestamp))
	return nil
}

func (s *PingService) GetPingCount(ctx context.Context) (int64, error) {
	count, err := s.repo.CountPings(ctx)
	if err != nil {
		s.logger.Error("failed to count pings", zap.Error(err))
		return 0, &errors.StorageError{Err: err}
	}

	s.logger.Info("ping count retrieved",
		zap.Int64("count", count))
	return count, nil
}
