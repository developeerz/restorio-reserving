package scheduler

import (
	"context"
	"time"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/outbox"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sender     Sender
	outboxRepo postgres.OutboxRepository
	scheduler  gocron.Scheduler
}

func New(ctx context.Context, sender Sender, outboxRepo postgres.OutboxRepository) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s := new(Scheduler)
	s.sender = sender
	s.outboxRepo = outboxRepo
	s.scheduler = scheduler

	s.scheduler.Start()

	err = s.scheduleDeleteSentJob(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Scheduler) ScheduleSendMessageJob(ctx context.Context, outboxMessage outbox.Entity) error {
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

func (s *Scheduler) scheduleDeleteSentJob(ctx context.Context) error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(deleteSentJob, ctx, s.outboxRepo),
	)
	if err != nil {
		return err
	}

	return nil
}
