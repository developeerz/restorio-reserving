// internal/port/notifier.go
package port

import "context"

type NotificationSender interface {
	Send(ctx context.Context, topic string, payload []byte) error
}
