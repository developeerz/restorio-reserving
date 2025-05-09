package postgres

import (
	"context"

	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/port"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) port.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) UserReservations(ctx context.Context, userID int) ([]dto.UserReservationResponse, error) {
	const q = `
        SELECT 
            rsv.reservation_id, 
            rsv.table_id, 
            t.table_number, 
            t.seats_number, 
            r.name AS restaurant_name,
            rsv.reservation_time_from,
            rsv.reservation_time_to
        FROM reservations rsv
        JOIN tables t ON rsv.table_id = t.table_id
        JOIN restaurants r ON t.restaurant_id = r.restaurant_id
        WHERE rsv.user_id = $1
        ORDER BY rsv.reservation_time_from;
    `
	var out []dto.UserReservationResponse
	if err := r.db.SelectContext(ctx, &out, q, userID); err != nil {
		return nil, err
	}
	return out, nil
}
