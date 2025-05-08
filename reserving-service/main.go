package main

import (
	"log"

	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/reserving-service/internal/db"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/kafka"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/routes"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/scheduler"
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

	outboxRepo := postgres.NewOutboxRepository(DB)

	kafkaSender := kafka.NewKafka(nil)
	sched, err := scheduler.New(kafkaSender, outboxRepo)
	if err != nil {
		log.Fatalf("scheduler init error: %v", err)
	}

	// Создаём роутер
	router := gin.Default()

	routes.SetupRoutes(router, DB, sched)

	// Запускаем сервер
	log.Println("Сервер запущен на порту 8082")
	router.Run(":8082")
}
