package scheduler

import (
	"context"

	"github.com/google/uuid"
)

type Sender interface {
	Send(ctx context.Context, topic string, payload []byte) error
	Close() error
}

type OutboxRepository interface {
	UpdateSendStatusTrueByID(ctx context.Context, id uuid.UUID) error
}
