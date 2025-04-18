package main

import (
	"log"

	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/reserving-service/internal/db"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/handlers"
	"github.com/gin-gonic/gin" // Необходимо для доступа к файлам Swagger UI
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Restorio API
// @version 1.0
// @description API для бронирования столиков
// @host localhost:8082
// @BasePath /
func main() {
	// Инициализируем БД
	db.InitDB()
	defer db.DB.Close()

	// Создаём роутер
	router := gin.Default()

	// Роут для Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Эндпоинты
	router.GET("/free-tables", handlers.GetFreeTables(db.DB))
	router.POST("/reservations", handlers.BookTable(db.DB))
	router.GET("/users/:user_id/reservations", handlers.GetUserReservations(db.DB))

	// Запускаем сервер
	log.Println("Сервер запущен на порту 8082")
	router.Run(":8082")
}
