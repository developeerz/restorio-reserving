package models

type PayloadTelegram struct {
	RestaurantName    string `json:"restaurant_name"`
	RestaurantAddress string `json:"restaurant_address"`
	TableNumber       string `json:"table_number"`
	ReservationTime   string `json:"reservation_time"`
	TelegramID        int    `json:"telegram_id"`
}
