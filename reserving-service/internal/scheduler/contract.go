package scheduler

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context, payload []byte) error
	Close() error
}
