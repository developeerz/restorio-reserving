package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID         uuid.UUID
	Topic      string
	Payload    []byte
	SendTime   time.Time
	SendStatus bool
}

func NewOutboxEntity(topic string, payload []byte, sentTime time.Time) *Entity {
	return &Entity{
		ID:         uuid.New(),
		Topic:      topic,
		Payload:    payload,
		SendTime:   sentTime,
		SendStatus: false,
	}
}
