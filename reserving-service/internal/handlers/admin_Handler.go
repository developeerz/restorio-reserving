package handlers

import (
	"database/sql"
	"net/http"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
)

// CreateTableHandler добавляет новый столик и его позицию (если указана).
// @Summary Добавить столик
// @Description Создаёт новый столик в ресторане и сохраняет его позицию, если координаты указаны
// @Tags Admin
// @Accept json
// @Produce json
// @Param table body dto.CreateTableRequest true "Информация о столике"
// @Success 201 {object} map[string]interface{} "Столик добавлен"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 500 {object} map[string]string "Ошибка при добавлении столика"
// @Router /tables [post]
func CreateTableHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTableRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		// Вставка столика с учётом shape
		var tableID int
		query := `
			INSERT INTO "Tables" (restaurant_id, table_number, seats_number, type, shape)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING table_id
		`
		err := db.QueryRow(query, req.RestaurantID, req.TableNumber, req.SeatsNumber, req.Type, req.Shape).Scan(&tableID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении столика"})
			return
		}

		// Вставка позиции, если координаты указаны
		if req.X != nil && req.Y != nil {
			posQuery := `INSERT INTO positions (table_id, x, y) VALUES ($1, $2, $3)`
			_, err = db.Exec(posQuery, tableID, *req.X, *req.Y)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении позиции столика"})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"table_id": tableID, "message": "Столик добавлен"})
	}
}
