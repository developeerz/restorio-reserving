package port

import (
	"context"
	"time"

	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
	"github.com/google/uuid"
)

// OutboxRepository описывает, что мы можем делать с таблицей outbox
type OutboxRepository interface {
	// Transaction запускает fn в рамках одной транзакции
	Transaction(ctx context.Context, fn func(repo OutboxRepository) error) error

	// UpdateSendStatusTrueByID помечает сообщение как отправленное
	UpdateSendStatusTrueByID(ctx context.Context, id uuid.UUID) error

	// DeleteSent удаляет все записи, у которых send_status=true
	DeleteSent(ctx context.Context) error

	// CreateOutbox создает новую запись в таблице outbox
	CreateOutbox(ctx context.Context, id, topic string, payload []byte, sendTime time.Time) error

	// GetTablePayload возвращает тему и уже сериализованный JSON-пейлоад
	GetTablePayload(ctx context.Context, tableID int) (topic string, payload []byte, err error)

	// GetByID извлекает сообщение по ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Entity, error)
}
