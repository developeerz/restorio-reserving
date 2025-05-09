package scheduler

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context, topic string, payload []byte) error
	Close() error
}
