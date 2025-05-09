package scheduler

import (
	"context"

	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sender     port.NotificationSender
	outboxRepo port.OutboxRepository
	cron       gocron.Scheduler
}

// New создаёт и запускает планировщик (но не регистрирует send-задачи)
func New(ctx context.Context, sender port.NotificationSender, outboxRepo port.OutboxRepository) (port.Scheduler, error) {
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s := &Scheduler{
		sender:     sender,
		outboxRepo: outboxRepo,
		cron:       cron,
	}

	return s, s.Start(ctx)
}

// ScheduleSendMessageJob регистрирует **однократную** задачу
func (s *Scheduler) ScheduleSendMessageJob(ctx context.Context, outboxMessage outbox.Entity) error {
	_, err := s.NewJob(
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

// Start регистрирует периодическую очистку outbox
func (s *Scheduler) Start(ctx context.Context) error {
	_, err := s.cron.Every(1).Hour().Do(func() { deleteSentJob(ctx, s.outboxRepo) })
	return err
}

// Stop останавливает планировщик
func (s *Scheduler) Stop() {
	s.cron.Stop()
}
