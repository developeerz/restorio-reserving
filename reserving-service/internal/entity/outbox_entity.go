package entity

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEntity struct {
	ID         uuid.UUID
	Topic      string
	Payload    []byte
	SendTime   time.Time
	SendStatus bool
}

func NewOutboxEntity(topic string, payload []byte, sentTime time.Time) *OutboxEntity {
	return &OutboxEntity{
		ID:         uuid.New(),
		Topic:      topic,
		Payload:    payload,
		SendTime:   sentTime,
		SendStatus: false,
	}
}
