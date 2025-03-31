package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type FreeTable struct {
	TableID        int    `json:"table_id" db:"table_id"`
	TableNumber    int    `json:"table_number" db:"table_number"`
	SeatsNumber    int    `json:"seats_number" db:"seats_number"`
	RestaurantName string `json:"restaurant_name" db:"restaurant_name"`
}

// ReservationRequest – структура запроса для бронирования столика.
type ReservationRequest struct {
	TableID             int    `json:"table_id"`              // Идентификатор столика
	UserID              int    `json:"user_id"`               // Идентификатор пользователя
	ReservationTimeFrom string `json:"reservation_time_from"` // Время начала бронирования (RFC3339)
	ReservationTimeTo   string `json:"reservation_time_to"`   // Время окончания бронирования (RFC3339)
}

// UserReservation – структура, описывающая бронирование с деталями столика и ресторана.
type UserReservation struct {
	ReservationID       int       `db:"reservation_id" json:"reservation_id"`
	TableID             int       `db:"table_id" json:"table_id"`
	TableNumber         int       `db:"table_number" json:"table_number"`
	SeatsNumber         int       `db:"seats_number" json:"seats_number"`
	RestaurantName      string    `db:"restaurant_name" json:"restaurant_name"`
	ReservationTimeFrom time.Time `db:"reservation_time_from" json:"reservation_time_from"`
	ReservationTimeTo   time.Time `db:"reservation_time_to" json:"reservation_time_to"`
}

func GetFreeTables(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		reservationFrom := c.Query("reservation_time_from")
		reservationTo := c.Query("reservation_time_to")

		// Парсим время
		fromTime, err := time.Parse(time.RFC3339, reservationFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation_time_from format"})
			return
		}

		toTime, err := time.Parse(time.RFC3339, reservationTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation_time_to format"})
			return
		}

		// SQL-запрос
		query := `
			SELECT t.table_id, t.table_number, t.seats_number, r.name AS restaurant_name
			FROM tables t
			JOIN restaurants r ON t.restaurant_id = r.restaurant_id
			WHERE t.table_id NOT IN (
				SELECT table_id FROM reservations 
				WHERE reservation_time_from < $1
				AND reservation_time_to > $2
			);
		`

		var freeTables []FreeTable
		err = db.Select(&freeTables, query, fromTime, toTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
			return
		}

		c.JSON(http.StatusOK, freeTables)
	}
}

// BookTable – обработчик для бронирования столика.
func BookTable(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ReservationRequest
		// Разбираем JSON из запроса
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Парсинг временных меток
		fromTime, err := time.Parse(time.RFC3339, req.ReservationTimeFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат reservation_time_from"})
			return
		}

		toTime, err := time.Parse(time.RFC3339, req.ReservationTimeTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат reservation_time_to"})
			return
		}

		// SQL-запрос для бронирования
		query := `
			INSERT INTO reservations (table_id, user_id, reservation_time_from, reservation_time_to)
			VALUES ($1, $2, $3, $4)
			RETURNING reservation_id;
		`

		var reservationID int
		// Используем QueryRow для вставки с возвратом идентификатора
		err = db.QueryRow(query, req.TableID, req.UserID, fromTime, toTime).Scan(&reservationID)
		if err != nil {
			// Если возникла ошибка (например, конфликт бронирования или другая ошибка БД)
			if err == sql.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить бронирование"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка бронирования: " + err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reservation_id": reservationID,
			"message":        "Столик успешно забронирован",
		})
	}
}

// GetUserReservations – обработчик, возвращающий все бронирования для заданного пользователя.
func GetUserReservations(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем user_id из параметров URL
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат user_id"})
			return
		}

		// SQL-запрос для получения бронирований с деталями столика и ресторана
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

		var reservations []UserReservation
		err = db.Select(&reservations, query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, reservations)
	}
}
