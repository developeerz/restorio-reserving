package dto

import "time"

/* Requests */
type ReservationRequest struct {
	TableID             int    `json:"table_id"`              // Идентификатор столика
	UserID              int    `json:"user_id"`               // Идентификатор пользователя
	ReservationTimeFrom string `json:"reservation_time_from"` // Время начала бронирования (RFC3339)
	ReservationTimeTo   string `json:"reservation_time_to"`   // Время окончания бронирования (RFC3339)
}

/* Responses */
type FreeTable struct {
	TableID        int    `json:"table_id" db:"table_id"`
	TableNumber    int    `json:"table_number" db:"table_number"`
	SeatsNumber    int    `json:"seats_number" db:"seats_number"`
	RestaurantName string `json:"restaurant_name" db:"restaurant_name"`
}

type TimeSlot struct {
	FreeFrom  time.Time `json:"free_from"`
	FreeUntil time.Time `json:"free_until"`
}
