package entity

import (
	"time"

	"github.com/google/uuid"
)

// Entity — строка таблицы outbox
type Outbox struct {
	ID         uuid.UUID `db:"id"`
	Topic      string    `db:"topic"`
	Payload    []byte    `db:"payload"`
	SendTime   time.Time `db:"send_time"`
	SendStatus bool      `db:"send_status"`
}

// NewEntity создаёт новую запись для вставки
func NewEntity(topic string, payload []byte, sendTime time.Time) *Outbox {
	return &Outbox{
		ID:         uuid.New(),
		Topic:      topic,
		Payload:    payload,
		SendTime:   sendTime,
		SendStatus: false,
	}
}
