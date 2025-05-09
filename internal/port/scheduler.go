package port

import (
	"context"

	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
)

type Scheduler interface {
	ScheduleSendMessageJob(ctx context.Context, outboxMessage entity.Outbox) error
}
