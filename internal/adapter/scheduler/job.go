package scheduler

import (
	"context"

	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/google/uuid"
)

// sendMessageJob вызывается одной разовой задачей
func sendMessageJob(ctx context.Context, sender port.NotificationSender, outboxRepo port.OutboxRepository, id uuid.UUID) {
	// внутри Transaction помечаем отправку и сами отправляем
	_ = outboxRepo.Transaction(ctx, func(repo port.OutboxRepository) error {
		if err := repo.UpdateSendStatusTrueByID(ctx, id); err != nil {
			return err
		}
		// получаем из базы entity по id, но у нас outbox_entity будет читать все поля
		e := entity.NewEntityFromID(id) // вспомогательный метод, реализуйте под себя
		return sender.Send(ctx, e.Topic, e.Payload)
	})
}

// deleteSentJob — периодически чистит таблицу
func deleteSentJob(ctx context.Context, outboxRepo port.OutboxRepository) {
	_ = outboxRepo.DeleteSent(ctx)
}
