package scheduler

import (
	"context"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/outbox"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sender     Sender
	outboxRepo OutboxRepository
	scheduler  gocron.Scheduler
}

func New(sender Sender, outboxRepo OutboxRepository) (*Scheduler, error) {
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

func sendMessageJob(ctx context.Context, sender Sender, repo OutboxRepository, outboxMessage outbox.Entity) error {
	err := sender.Send(ctx, outboxMessage.Topic, outboxMessage.Payload)
	if err != nil {
		return err
	}

	return repo.UpdateSendStatusTrueByID(ctx, outboxMessage.ID)
}
