package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/usecase"

	"github.com/gin-gonic/gin"
)

// GetFreeTables godoc
// @Summary Получить список свободных столиков
// @Description Возвращает список доступных столиков на указанный период
// @Tags Booking
// @Produce json
// @Param reservation_time_from query string true "Начало периода бронирования (в формате RFC3339)"
// @Param reservation_time_to query string true "Конец периода бронирования (в формате RFC3339)"
// @Success 200 {array} dto.FreeTableResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/free-tables [get]
func GetFreeTablesHandler(uc *usecase.BookingUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		fromStr := c.Query("reservation_time_from")
		toStr := c.Query("reservation_time_to")

		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: fmt.Sprintf("Invalid reservation_time_from.\n Details: %s", err.Error())})
			return
		}
		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: fmt.Sprintf("Invalid reservation_time_to.\n Details: %s", err.Error())})
			return
		}
		if to.Before(from) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "reservation_time_to must be after reservation_time_from"})
			return
		}

		list, err := uc.GetFreeTables(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Error fetching free tables.\n Details: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

// BookTable godoc
// @Summary Бронирование столика
// @Description Эта функция выполняет бронирование столика для пользователя на указанный период времени.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param reservation body dto.ReservationRequest true "Информация о бронировании"
// @Success 200 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/new-reservation [post]
func BookTableHandler(uc *usecase.BookingUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ReservationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}

		id, err := uc.BookTable(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Error booking table.\n Details: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reservation_id": id, "message": "Столик успешно забронирован"})
	}
}

// GetFreeTimeSlotsHandler godoc
// @Summary Получить свободные временные интервалы для столика
// @Description Возвращает список свободных временных интервалов для бронирования указанного столика.
// @Tags Booking
// @Param table_id path int true "ID столика"
// @Produce json
// @Success 200 {array} dto.TimeSlotResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/:table_id/free-times [get]
func GetFreeTimeSlotsHandler(uc *usecase.BookingUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		tidStr := c.Param("table_id")
		tid, err := strconv.Atoi(tidStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid table_id"})
			return
		}
		slots, err := uc.GetFreeTimeSlots(c.Request.Context(), tid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Error fetching free time slots.\n Details: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, slots)
	}
}
