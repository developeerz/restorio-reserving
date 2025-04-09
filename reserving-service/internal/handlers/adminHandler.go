package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTableRequest struct {
	RestaurantID int    `json:"restaurant_id" binding:"required"`
	TableNumber  string `json:"table_number"`
	SeatsNumber  int    `json:"seats_number" binding:"required"`
	Type         string `json:"type"`
}

func CreateTableHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateTableRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		var tableID int
		query := `INSERT INTO tables (restaurant_id, table_number, seats_number, type)
		          VALUES ($1, $2, $3, $4) RETURNING table_id`
		err := db.QueryRow(query, req.RestaurantID, req.TableNumber, req.SeatsNumber, req.Type).Scan(&tableID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении столика"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"table_id": tableID, "message": "Столик добавлен"})
	}
}
