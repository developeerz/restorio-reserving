package handlers

import (
	"net/http"
	"strconv"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetUserReservations(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем user_id из параметров URL
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат user_id"})
			return
		}

		/* SQL-query */
		query := `
			SELECT 
				rsv.reservation_id, 
				rsv.table_id, 
				t.table_number, 
				t.seats_number, 
				r.name AS restaurant_name,
				rsv.reservation_time_from,
				rsv.reservation_time_to
			FROM reservations rsv
			JOIN tables t ON rsv.table_id = t.table_id
			JOIN restaurants r ON t.restaurant_id = r.restaurant_id
			WHERE rsv.user_id = $1
			ORDER BY rsv.reservation_time_from;
		`

		/* request */
		var reservations []dto.UserReservation
		err = db.Select(&reservations, query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
			return
		}

		/* response */
		c.JSON(http.StatusOK, reservations)
	}
}
