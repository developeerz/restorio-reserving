package scheduler

import (
	"context"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/outbox"
)

func sendMessageJob(ctx context.Context, sender Sender, repo postgres.OutboxRepository, outboxMessage outbox.Entity) error {
	return repo.Transaction(ctx, func(repo postgres.OutboxRepository) error {
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

func deleteSentJob(ctx context.Context, repo postgres.OutboxRepository) error {
	return repo.DeleteSent(ctx)
}
