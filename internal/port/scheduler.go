package port

import (
	"context"

	"github.com/google/uuid"
)

// Scheduler описывает наш планировщик
type Scheduler interface {
	// ScheduleSendMessageJob планирует задачу отправки конкретного сообщения из outbox
	ScheduleSendMessageJob(ctx context.Context, outboxID uuid.UUID) error
	// Stop останавливает все задачи
	Stop()
}
