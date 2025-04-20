package dto

import "time"

/* Responses */
type UserReservationResponse struct {
	ReservationID       int       `db:"reservation_id" json:"reservation_id"`
	TableID             int       `db:"table_id" json:"table_id"`
	TableNumber         int       `db:"table_number" json:"table_number"`
	SeatsNumber         int       `db:"seats_number" json:"seats_number"`
	RestaurantName      string    `db:"restaurant_name" json:"restaurant_name"`
	ReservationTimeFrom time.Time `db:"reservation_time_from" json:"reservation_time_from"`
	ReservationTimeTo   time.Time `db:"reservation_time_to" json:"reservation_time_to"`
}
