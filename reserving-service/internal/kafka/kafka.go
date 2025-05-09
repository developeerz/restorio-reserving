package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	writer *kafka.Writer
}

func NewKafka(brokers []string) *Kafka {
	return &Kafka{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			MaxAttempts:  5,
			WriteTimeout: 10 * time.Second,
			Async:        false,
			RequiredAcks: kafka.RequireAll,
			BatchSize:    1,
		},
	}
}

func (k *Kafka) Send(ctx context.Context, topic string, payload []byte) error {
	return k.writer.WriteMessages(
		ctx,
		kafka.Message{
			Topic: topic,
			Value: payload,
		},
	)
}

func (k *Kafka) Close() error {
	return k.writer.Close()
}
