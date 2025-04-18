package dto

/* Responses */

/* Requests */
type CreateTableRequest struct {
	RestaurantID int    `json:"restaurant_id" binding:"required"`
	TableNumber  string `json:"table_number"`
	SeatsNumber  int    `json:"seats_number" binding:"required"`
	Type         string `json:"type" binding:"required"`  // ENUM TABLE_TYPE
	Shape        string `json:"shape" binding:"required"` // ENUM TABLE_SHAPE
	X            *int   `json:"x,omitempty"`              // Координаты опционально
	Y            *int   `json:"y,omitempty"`
}
