package port

import (
	"context"

	"github.com/developeerz/restorio-reserving/internal/dto"
)

// ReservationRepository описывает операции с бронированиями
type UserRepository interface {
	// UserReservations возвращает все бронирования для данного userID
	UserReservations(ctx context.Context, userID int) ([]dto.UserReservationResponse, error)
}
