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
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservations/new-reservation [post]
func BookTable(db *sqlx.DB, sched *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var req dto.ReservationRequest

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

		payload := mapper.ToPayload(payloadEntities[0], fromTime.Local().String(), -1)

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
		err = tx.QueryRow(query, req.TableID, -1, fromTime, toTime, "reserved", time.Now()).Scan(&reservationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование: " + err.Error()})
			return
		}

		createOutboxMessageQuery := `
			INSERT INTO outbox
				(id, topic, payload, send_time)
			VALUES
				($1, $2, $3, $4);
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
			c.Status(http.StatusInternalServerError)
			return
		}

		if err = sched.ScheduleSendMessageJob(*outboxMessage); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		/* response */
		c.JSON(http.StatusOK, gin.H{
			"reservation_id": reservationID,
			"message":        "Столик успешно забронирован",
		})
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
// @Router /tables/:table_id/free-times [get]
func GetFreeTimeSlotsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableIDStr := c.Param("table_id")
		tableID, err := strconv.Atoi(tableIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный table_id"})
			return
		}

		timeSlots, err := GetFreeTimeSlots(db, tableID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении свободных слотов: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, timeSlots)
	}
}

func GetFreeTimeSlots(db *sqlx.DB, tableID int) ([]dto.TimeSlotResponse, error) {
	query := `
    WITH "Booked_slots" AS (
        SELECT 
            reservation_time_from AS start_time, 
            reservation_time_to AS end_time
        FROM reservations
        WHERE table_id = $1
    )
    SELECT 
        lag(end_time, 1) OVER (ORDER BY start_time) AS free_from,
        start_time AS free_until
    FROM "Booked_slots";
    `
	rows, err := db.Query(query, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var freeTimeSlots []dto.TimeSlotResponse
	for rows.Next() {
		var slot dto.TimeSlotResponse
		if err := rows.Scan(&slot.FreeFrom, &slot.FreeUntil); err != nil {
			return nil, err
		}
		freeTimeSlots = append(freeTimeSlots, slot)
	}
	return freeTimeSlots, nil
}
