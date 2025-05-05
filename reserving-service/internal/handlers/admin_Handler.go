package handlers

import (
	"net/http"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
func CreateNewTable(db *sqlx.DB) gin.HandlerFunc {
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

/*
func UpdateTableWithPosition(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTableRequest

		tableIDStr := c.Param("table_id")
		tableID, err := strconv.Atoi(tableIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат table_id"})
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON: " + err.Error()})
			return
		}

		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка начала транзакции: " + err.Error()})
			return
		}
		defer tx.Rollback()

		// Обновление столика
		updateTableQuery := `
			UPDATE "Tables"
			SET table_number = $1, seats_number = $2, type = $3, shape = $4
			WHERE table_id = $5
		`
		_, err = tx.Exec(updateTableQuery, req.TableNumber, req.SeatsNumber, req.Type, req.Shape, tableID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления столика: " + err.Error()})
			return
		}

		// Обновление позиции (если передано)
		if req.X != nil && req.Y != nil {
			updatePositionQuery := `
				INSERT INTO "Positions" (table_id, x, y)
				VALUES ($1, $2, $3)
				ON CONFLICT (table_id) DO UPDATE
				SET x = EXCLUDED.x, y = EXCLUDED.y
			`
			_, err = tx.Exec(updatePositionQuery, tableID, *req.X, *req.Y)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления позиции: " + err.Error()})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изменений: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Столик успешно обновлён"})
	}
}


func DeleteTable(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableIDStr := c.Param("table_id")
		tableID, err := strconv.Atoi(tableIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат table_id"})
			return
		}

		query := `DELETE FROM "Tables" WHERE table_id = $1`
		res, err := db.Exec(query, tableID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления столика: " + err.Error()})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Столик с таким ID не найден"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Столик успешно удалён"})
	}
}
*/
