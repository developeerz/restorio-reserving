package postgres

import (
	"context"
	"time"

	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/jmoiron/sqlx"
)

// BookingRepo — реализация порт-репозитория для бронирований
type BookingRepo struct {
	db *sqlx.DB
}

// NewBookingRepo создаёт новый экземпляр BookingRepo
func NewBookingRepo(db *sqlx.DB) port.BookingRepository {
	return &BookingRepo{db: db}
}

// FreeTables возвращает все свободные столики в период [from, to)
func (r *BookingRepo) FreeTables(ctx context.Context, from, to time.Time) ([]dto.FreeTableResponse, error) {
	var out []dto.FreeTableResponse
	query := `SELECT t.table_id, t.table_number, t.seats_number, r.name AS restaurant_name
			  FROM tables t
			  JOIN restaurants r ON t.restaurant_id = r.restaurant_id
			  WHERE t.table_id NOT IN (
			    SELECT table_id FROM reservations
				WHERE NOT (reservation_time_to <= $1 OR reservation_time_from >= $2)
			  )`
	err := r.db.SelectContext(ctx, &out, query, from, to)
	return out, err
}

// CreateReservation вставляет бронирование и возвращает его ID
func (r *BookingRepo) CreateReservation(ctx context.Context, tableID, userID int, from, to time.Time) (int, error) {
	var id int
	query := `INSERT INTO reservations
				(table_id, user_id, reservation_time_from, reservation_time_to, status, created_at)
			  VALUES ($1, $2, $3, $4, 'reserved', NOW())
			  RETURNING reservation_id`
	err := r.db.QueryRowContext(ctx, query, tableID, userID, from, to).Scan(&id)
	return id, err
}

// FreeTimeSlots возвращает интервалы свободного времени конкретного столика
func (r *BookingRepo) FreeTimeSlots(ctx context.Context, tableID int) ([]dto.TimeSlotResponse, error) {
	var slots []dto.TimeSlotResponse
	query := `WITH booked AS (
				SELECT reservation_time_from AS start_time, reservation_time_to AS end_time
				FROM reservations WHERE table_id = $1
			  )
			  SELECT lag(end_time) OVER (ORDER BY start_time) AS free_from,
					 start_time AS free_until
			  FROM booked`
	err := r.db.SelectContext(ctx, &slots, query, tableID)
	return slots, err
}
