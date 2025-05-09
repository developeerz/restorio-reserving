package scheduler

import (
	"context"
	"time"

	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/go-co-op/gocron/v2"
)

type scheduler struct {
	sender     port.NotificationSender
	outboxRepo port.OutboxRepository
	cron       gocron.Scheduler
}

func New(ctx context.Context, sender port.NotificationSender, outboxRepo port.OutboxRepository) (port.Scheduler, error) {
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s := &scheduler{
		sender:     sender,
		outboxRepo: outboxRepo,
		cron:       cron,
	}

	return s, s.scheduleDeleteSentJob(ctx)
}

func (s *scheduler) ScheduleSendMessageJob(ctx context.Context, outboxMessage entity.Outbox) error {
	_, err := s.cron.NewJob(
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

func (s *scheduler) scheduleDeleteSentJob(ctx context.Context) error {
	_, err := s.cron.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(deleteSentJob, ctx, s.outboxRepo),
	)
	if err != nil {
		return err
	}

	return nil
}
