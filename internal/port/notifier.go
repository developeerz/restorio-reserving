// internal/port/notifier.go
package port

import "context"

// NotificationSender — порт для отправки уведомлений (через Kafka, Email и т.п.)
type NotificationSender interface {
	// Send публикует raw-пейлоад в указанный топик
	Send(ctx context.Context, topic string, payload []byte) error
}
