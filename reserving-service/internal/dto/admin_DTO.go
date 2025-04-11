package dto

/* Responses */
type CreateTableRequest struct {
	RestaurantID int    `json:"restaurant_id" binding:"required"`
	TableNumber  string `json:"table_number"`
	SeatsNumber  int    `json:"seats_number" binding:"required"`
	Type         string `json:"type"`
}
