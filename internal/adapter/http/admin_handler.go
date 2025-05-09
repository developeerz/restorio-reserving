package http

import (
	"net/http"

	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/port"

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
// @Router /tables/new-table [post]
func CreateTableHandler(repo port.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTableRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		tableID, err := repo.CreateTable(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating table"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"table_id": tableID,
			"message":  "Table created",
		})
	}
}
