package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/utilities"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetFreeTables godoc
// @Summary Получить список свободных столиков
// @Description Возвращает список доступных столиков на указанный период
// @Tags Tables
// @Produce json
// @Param reservation_time_from query string true "Начало периода бронирования (в формате RFC3339)"
// @Param reservation_time_to query string true "Конец периода бронирования (в формате RFC3339)"
// @Success 200 {array} dto.FreeTable
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /free-tables [get]
func GetFreeTables(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		reservationFrom := c.Query("reservation_time_from")
		reservationTo := c.Query("reservation_time_to")

		/* [START] PARSING ARGUMENTS */
		fromTime, err := time.Parse(time.RFC3339, reservationFrom)
		if utilities.JSONError(c, err, http.StatusBadRequest, "Invalid reservation_time_from format", "Your value: ", fromTime) {
			return
		}

		toTime, err := time.Parse(time.RFC3339, reservationTo)
		if utilities.JSONError(c, err, http.StatusBadRequest, "Invalid reservation_time_to format", "Your value: ", toTime) {
			return
		}

		//Дополнительные проверки для подстраховки
		if toTime.Before(fromTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "reservation_time_to must be after reservation_time_from"})
			return
		}
		/* [END] PARSING ARGUMENTS */

		/* SQL-query */
		query := `
			SELECT t.table_id, t.table_number, t.seats_number, r.name AS restaurant_name
			FROM tables t
			JOIN restaurants r ON t.restaurant_id = r.restaurant_id
			WHERE t.table_id NOT IN (
				SELECT table_id FROM reservations 
				WHERE NOT (reservation_time_to <= $1 OR reservation_time_from >= $2)
			);
		`

		/* request */
		var freeTables []dto.FreeTable
		err = db.Select(&freeTables, query, fromTime, toTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о свободных столиках.", "details": err.Error()})
			return
		}

		/* response */
		c.JSON(http.StatusOK, freeTables)
	}
}

//!!! Поправь ERROR RESPONSES

// BookTable godoc
// @Summary Бронирование столика
// @Description Эта функция выполняет бронирование столика для пользователя на указанный период времени.
// @Tags Reservations
// @Accept  json
// @Produce  json
// @Param reservation body dto.ReservationRequest true "Информация о бронировании"
// @Success 200 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations [post]
func BookTable(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ReservationRequest
		/* Parse input JSON */
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		/* Parse arguments */
		fromTime, err := time.Parse(time.RFC3339, req.ReservationTimeFrom)
		if utilities.JSONError(c, err, http.StatusBadRequest, "Invalid reservation_time_from format", "Your value: ", fromTime) {
			return
		}

		toTime, err := time.Parse(time.RFC3339, req.ReservationTimeTo)
		if utilities.JSONError(c, err, http.StatusBadRequest, "Invalid reservation_time_to format", "Your value: ", toTime) {
			return
		}

		//Дополнительные проверки для подстраховки
		if toTime.Before(fromTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "reservation_time_to must be after reservation_time_from"})
			return
		}

		/* SQL-query */
		query := `
			INSERT INTO reservations (table_id, user_id, reservation_time_from, reservation_time_to)
			VALUES ($1, $2, $3, $4)
			RETURNING reservation_id;
		`

		/* request */
		var reservationID int
		err = db.QueryRow(query, req.TableID, req.UserID, fromTime, toTime).Scan(&reservationID) // Используем QueryRow для вставки с возвратом идентификатора
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование: " + err.Error()})
			return
		}

		/* response */
		c.JSON(http.StatusOK, gin.H{
			"reservation_id": reservationID,
			"message":        "Столик успешно забронирован",
		})
	}
}

func GetFreeTimeSlots(db *sql.DB, tableID int) ([]dto.TimeSlot, error) {
	query := `
    WITH booked_slots AS (
        SELECT 
            reservation_time_from AS start_time, 
            reservation_time_to AS end_time
        FROM reservations
        WHERE table_id = $1
    )
    SELECT 
        lag(end_time, 1) OVER (ORDER BY start_time) AS free_from,
        start_time AS free_until
    FROM booked_slots;
    `
	rows, err := db.Query(query, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var freeTimeSlots []dto.TimeSlot
	for rows.Next() {
		var slot dto.TimeSlot
		if err := rows.Scan(&slot.FreeFrom, &slot.FreeUntil); err != nil {
			return nil, err
		}
		freeTimeSlots = append(freeTimeSlots, slot)
	}
	return freeTimeSlots, nil
}
