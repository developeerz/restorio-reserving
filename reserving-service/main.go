package main

import (
	"log"

	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/reserving-service/internal/db"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/routes"
	"github.com/gin-gonic/gin" // Необходимо для доступа к файлам Swagger UI
	"github.com/jmoiron/sqlx"
)

// @title Restorio API
// @version 1.0
// @description API для бронирования столиков
// @host localhost:8082
// @BasePath /
func main() {
	// Инициализируем БД
	var DB *sqlx.DB
	DB = db.InitDB()
	defer DB.Close()

	// Создаём роутер
	router := gin.Default()

	routes.SetupRoutes(router, DB)

	// Запускаем сервер
	log.Println("Сервер запущен на порту 8082")
	router.Run(":8082")
}
