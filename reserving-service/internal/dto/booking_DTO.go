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
type FreeTableResponse struct {
	TableID        int    `json:"table_id" db:"table_id"`
	TableNumber    int    `json:"table_number" db:"table_number"`
	SeatsNumber    int    `json:"seats_number" db:"seats_number"`
	RestaurantName string `json:"restaurant_name" db:"restaurant_name"`
}

type TimeSlotResponse struct {
	FreeFrom  time.Time `json:"free_from"`
	FreeUntil time.Time `json:"free_until"`
}

type GetTablesByRestaurantIDRequest struct {
	RestaurantID int `form:"restaurant_id"`
}

type Table struct {
	TableID     int    `json:"table_id"`
	TableNumber string `json:"table_number"`
	SeatsNumber int    `json:"seats_number"`
	Type        string `json:"type"`
	Shape       string `json:"shape"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
}

type GetTablesByRestaurantIDResponse struct {
	Tables []Table `json:"tables"`
}
