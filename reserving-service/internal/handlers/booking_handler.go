package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/mapper"
	entity "github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/outbox"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/scheduler"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/utilities"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
		var freeTables []dto.FreeTableResponse
		err = db.Select(&freeTables, query, fromTime, toTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о свободных столиках.", "details": err.Error()})
			return
		}

		/* response */
		c.JSON(http.StatusOK, freeTables)
	}
}

// BookTable godoc
// @Summary Бронирование столика
// @Description Эта функция выполняет бронирование столика для пользователя на указанный период времени.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param reservation body dto.ReservationRequest true "Информация о бронировании"
// @Success 200 {object} dto.ReservationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/new-reservation [post]
func BookTable(db *sqlx.DB, sched *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var req dto.ReservationRequest

		userID := c.GetInt("userID")

		/* Parse input JSON */
		if err = c.ShouldBindJSON(&req); err != nil {
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

		getPayloadQuery := `
			SELECT 
				restaurants.name,
				restaurants.address,
				tables.table_number
			FROM restaurants
			JOIN tables ON tables.restaurant_id = restaurants.restaurant_id
			WHERE
				tables.table_id = $1;
		`
		var payloadEntities []entity.Payload

		err = db.Select(&payloadEntities, getPayloadQuery, req.TableID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(payloadEntities) == 0 {
			c.Status(http.StatusInternalServerError)
			return
		}

		payload := mapper.ToPayload(payloadEntities[0], fromTime.Local().String(), userID)

		payloadByte, err := json.Marshal(&payload)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		tx, err := db.BeginTx(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		defer tx.Rollback()

		/* SQL-query */
		query := `
			INSERT INTO reservations (table_id, user_id, reservation_time_from, reservation_time_to, status, created_at)
			VALUES ($1, $2, $3, $4, 'booked', NOW()::TIMESTAMP)
			RETURNING reservation_id;
		`

		var reservationID int
		err = tx.QueryRow(query, req.TableID, userID, fromTime, toTime).Scan(&reservationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование: " + err.Error()})
			return
		}

		createOutboxMessageQuery := `
			INSERT INTO outbox
				(id, topic, payload, send_time, send_status)
			VALUES
				($1, $2, $3, $4, $5);
		`

		outboxMessage := entity.NewOutboxEntity("", payloadByte, fromTime.Local().Add(-1*time.Hour).Truncate(60*time.Second))

		_, err = tx.ExecContext(
			c.Request.Context(),
			createOutboxMessageQuery,
			outboxMessage.ID,
			outboxMessage.Topic,
			outboxMessage.Payload,
			outboxMessage.SendTime,
			outboxMessage.SendStatus,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование: " + err.Error()})
			return
		}

		if err = sched.ScheduleSendMessageJob(*outboxMessage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование: " + err.Error()})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование: " + err.Error()})
			return
		}

		/* response */
		c.JSON(http.StatusOK, dto.ReservationResponse{ReservationID: reservationID, Message: "Столик забронирован"})
	}
}

// GetFreeTimeSlotsHandler godoc
// @Summary Получить свободные временные интервалы для столика
// @Description Возвращает список свободных временных интервалов для бронирования указанного столика.
// @Tags Booking
// @Param table_id path int true "ID столика" example(1)
// @Param start    query string true "Начало интервала (формат RFC3339)" example(2025-03-26T08:00:00Z)
// @Param end      query string true "Конец интервала (формат RFC3339)" example(2025-03-26T22:00:00Z)
// @Produce json
// @Success 200 {array} dto.TimeSlotResponse "Пример:\n[  {\"free_from\": {\"Time\":\"2025-03-26T08:00:00Z\",\"Valid\":true},\"free_until\": {\"Time\":\"2025-03-26T10:00:00Z\",\"Valid\":true}},  {\"free_from\": {\"Time\":\"2025-03-26T12:00:00Z\",\"Valid\":true},\"free_until\": {\"Time\":\"2025-03-26T18:00:00Z\",\"Valid\":true}}]"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tables/{table_id}/free-times [get]
func GetFreeTimeSlotsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableIDStr := c.Param("table_id")
		tableID, err := strconv.Atoi(tableIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный table_id"})
			return
		}

		startStr := c.Query("start")
		endStr := c.Query("end")

		if startStr == "" || endStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Параметры `start` и `end` обязательны"})
			return
		}

		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра `start`, ожидается RFC3339"})
			return
		}

		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра `end`, ожидается RFC3339"})
			return
		}

		if !start.Before(end) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "`start` должен быть раньше `end`"})
			return
		}

		timeSlots, err := GetFreeTimeSlots(db, tableID, start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении свободных слотов: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, timeSlots)
	}
}

func GetFreeTimeSlots(db *sqlx.DB, tableID int, dayStart, dayEnd time.Time) ([]dto.TimeSlotResponse, error) {
	query := `
	WITH params AS (
		SELECT 
			$1::int AS table_id,
			$2::timestamp AS day_start,
			$3::timestamp AS day_end
	),
	booked_slots AS (
		SELECT 
			reservation_time_from AS start_time,
			reservation_time_to AS end_time
		FROM reservations r
		JOIN params p ON r.table_id = p.table_id
		WHERE reservation_time_from >= p.day_start AND reservation_time_to <= p.day_end
	),
	ordered_slots AS (
		SELECT
			start_time,
			end_time,
			LAG(end_time) OVER (ORDER BY start_time) AS prev_end
		FROM booked_slots
	),
	free_slots AS (
		-- Интервалы между бронированиями
		SELECT
			prev_end AS free_from,
			start_time AS free_until
		FROM ordered_slots
		WHERE prev_end IS NOT NULL

		UNION ALL

		-- Интервал до первой брони
		SELECT
			p.day_start AS free_from,
			(SELECT MIN(start_time) FROM booked_slots) AS free_until
		FROM params p

		UNION ALL

		-- Интервал после последней брони
		SELECT
			(SELECT MAX(end_time) FROM booked_slots) AS free_from,
			p.day_end AS free_until
		FROM params p
	)
	SELECT *
	FROM free_slots
	WHERE free_from < free_until
	ORDER BY free_from;
	`

	var slots []dto.TimeSlotResponse
	err := db.Select(&slots, query, tableID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	return slots, nil
}
