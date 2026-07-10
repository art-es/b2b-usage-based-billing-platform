package kafka

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(ctx context.Context) (*Producer, error) {
	url := os.Getenv("KAFKA_URL")
	if url == "" {
		return nil, errors.New("KAFKA_URL required")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(url),
		Async:        true,
		WriteTimeout: 10 * time.Second,
	}

	// simulating ping
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := writer.WriteMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Produce(ctx context.Context, msgs []broker.ProduceMessage) error {
	kafkaMsgs := make([]kafka.Message, 0, len(msgs))
	for _, msg := range msgs {
		kafkaMsgs = append(kafkaMsgs, kafka.Message{
			Topic: msg.Topic,
			Key:   msg.Key,
			Value: msg.Value,
		})
	}

	return p.writer.WriteMessages(ctx, kafkaMsgs...)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
