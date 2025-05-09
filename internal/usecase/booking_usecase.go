package usecase

import (
	"context"
	"time"

	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/google/uuid"
)

type BookingUseCase struct {
	repo      port.BookingRepository
	outbox    port.OutboxRepository
	scheduler port.Scheduler
}

func NewBookingUseCase(r port.BookingRepository, o port.OutboxRepository, s port.Scheduler) *BookingUseCase {
	return &BookingUseCase{repo: r, scheduler: s}
}

func (uc *BookingUseCase) GetFreeTables(ctx context.Context, from, to time.Time) ([]dto.FreeTableResponse, error) {
	return uc.repo.FreeTables(ctx, from, to)
}

func (uc *BookingUseCase) BookTable(ctx context.Context, req dto.ReservationRequest) (int, error) {
	// парсим times уже сделал в хендлере
	from, _ := time.Parse(time.RFC3339, req.ReservationTimeFrom)
	to, _ := time.Parse(time.RFC3339, req.ReservationTimeTo)

	// 1) Создаем в reservations
	reservationID, err := uc.repo.CreateReservation(ctx, req.TableID, req.UserID, from, to)
	if err != nil {
		return 0, err
	}

	// 2) Формируем сообщение и outbox
	topic, payload, err := uc.outbox.GetTablePayload(ctx, req.TableID)
	if err != nil {
		return 0, err
	}

	// используем uuid для id outbox
	id := uuid.New()
	sendTime := to.Add(-1 * time.Hour)

	if err := uc.outbox.CreateOutbox(ctx, id.String(), topic, payload, sendTime); err != nil {
		return 0, err
	}

	// 3) Планируем задачу
	if err := uc.scheduler.ScheduleSendMessageJob(ctx, id); err != nil {
		return 0, err
	}

	return reservationID, nil
}

func (uc *BookingUseCase) GetFreeTimeSlots(ctx context.Context, tableID int) ([]dto.TimeSlotResponse, error) {
	return uc.repo.FreeTimeSlots(ctx, tableID)
}
