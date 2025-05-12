package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	writer *kafka.Writer
	topic  string
}

func NewKafka(brokers []string, topic string) *Kafka {
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
		topic: topic,
	}
}

func (k *Kafka) Send(ctx context.Context, payload []byte) error {
	return k.writer.WriteMessages(
		ctx,
		kafka.Message{
			Topic: k.topic,
			Value: payload,
		},
	)
}

func (k *Kafka) Close() error {
	return k.writer.Close()
}
