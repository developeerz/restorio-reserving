package scheduler

import (
	"context"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/outbox"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sender     Sender
	outboxRepo postgres.OutboxRepository
	scheduler  gocron.Scheduler
}

func New(sender Sender, outboxRepo postgres.OutboxRepository) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s := new(Scheduler)
	s.sender = sender
	s.outboxRepo = outboxRepo
	s.scheduler = scheduler

	s.scheduler.Start()

	return s, nil
}

func (s *Scheduler) ScheduleJob(ctx context.Context, outboxMessage outbox.Entity) error {
	_, err := s.scheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(outboxMessage.SendTime),
		),
		gocron.NewTask(sendMessageJob, ctx, s.sender, s.outboxRepo, outboxMessage),
	)
	if err != nil {
		return err
	}

	return nil
}

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
