package http

import (
	"net/http"
	"strconv"

	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/port"

	"github.com/gin-gonic/gin"
)

// GetUserReservationsHandler godoc
// @Summary Получение всех бронирований пользователя
// @Description Возвращает список всех бронирований, сделанных пользователем, включая информацию о ресторане и столике.
// @Tags User
// @Param user_id header int true "ID пользователя"
// @Produce json
// @Success 200 {array} dto.UserReservationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/user [get]
func GetUserReservationsHandler(repo port.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Неверный формат user_id в заголовке"})
			return
		}

		resvs, err := repo.UserReservations(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Ошибка получения данных"})
			return
		}

		c.JSON(http.StatusOK, resvs)
	}
}
