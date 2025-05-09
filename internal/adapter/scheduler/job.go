package scheduler

import (
	"context"

	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
	"github.com/developeerz/restorio-reserving/internal/port"
)

func sendMessageJob(ctx context.Context, sender port.NotificationSender, outboxRepo port.OutboxRepository, outboxMessage entity.Outbox) error {
	return outboxRepo.Transaction(ctx, func(repo port.OutboxRepository) error {
		err := repo.UpdateSendStatusTrueByID(ctx, outboxMessage.ID)
		if err != nil {
			return err
		}

		err = sender.Send(ctx, outboxMessage.Topic, outboxMessage.Payload)
		if err != nil {
			return err
		}

		return nil
	})
}

func deleteSentJob(ctx context.Context, repo port.OutboxRepository) error {
	return repo.DeleteSent(ctx)
}
