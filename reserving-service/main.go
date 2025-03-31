package main

import (
	"log"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/db"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализируем БД
	db.InitDB()
	defer db.DB.Close()

	// Создаём роутер
	router := gin.Default()

	router.GET("/free-tables", handlers.GetFreeTables(db.DB))

	// Эндпоинт для бронирования столика
	router.POST("/reservations", handlers.BookTable(db.DB))

	// Новый эндпоинт для получения бронирований пользователя
	router.GET("/users/:user_id/reservations", handlers.GetUserReservations(db.DB))

	// Запускаем сервер
	log.Println("Сервер запущен на порту 8082")
	router.Run(":8082")
}
