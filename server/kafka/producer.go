package kafka

import (
	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) SendPingEvent(ctx context.Context, timestamp time.Time) error {
	event := &pingv1.PingRequest{
		TimestampMs: timestamp.UnixNano() / int64(time.Millisecond),
	}

	data, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "ping-events",
		Value: sarama.StringEncoder(data),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
