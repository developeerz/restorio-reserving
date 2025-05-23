package handlers

import (
	"net/http"
	"strconv"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetUserReservations godoc
// @Summary Получение всех бронирований пользователя
// @Description Возвращает список всех бронирований, сделанных пользователем, включая информацию о ресторане и столике.
// @Tags User
// @Param user_id header int true "ID пользователя"
// @Produce json
// @Success 200 {array} dto.UserReservationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/user [get]
func GetUserReservations(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем user_id из заголовка
		userIDStr := c.GetHeader("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат user_id в заголовке"})
			return
		}

		// SQL-запрос
		query := `
			SELECT 
				rsv.reservation_id, 
				rsv.table_id, 
				t.table_number, 
				t.seats_number, 
				r.name AS restaurant_name,
				rsv.reservation_time_from,
				rsv.reservation_time_to
			FROM "Reservations" rsv
			JOIN "Tables" t ON rsv.table_id = t.table_id
			JOIN "Restaurants" r ON t.restaurant_id = r.restaurant_id
			WHERE rsv.user_id = $1
			ORDER BY rsv.reservation_time_from;
		`

		var reservations []dto.UserReservationResponse
		err = db.Select(&reservations, query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, reservations)
	}
}
