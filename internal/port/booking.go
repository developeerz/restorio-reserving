package port

import (
	"context"
	"time"

	"github.com/developeerz/restorio-reserving/internal/dto"
)

// BookingRepository описывает методы для работы с бронированиями
type BookingRepository interface {
	FreeTables(ctx context.Context, from, to time.Time) ([]dto.FreeTableResponse, error)
	CreateReservation(ctx context.Context, tableID, userID int, from, to time.Time) (int, error)
	FreeTimeSlots(ctx context.Context, tableID int) ([]dto.TimeSlotResponse, error)
}
