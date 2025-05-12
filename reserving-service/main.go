package main

import (
	"context"
	"log"

	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/reserving-service/internal/config"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/db"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/kafka"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/routes"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/scheduler"
	"github.com/gin-gonic/gin" // Необходимо для доступа к файлам Swagger UI
)

// @title Restorio API
// @version 1.0
// @description API для бронирования столиков
// @host localhost:8082
// @BasePath /
func main() {
	ctx := context.Background()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error while load config: %v", err)
	}

	DB := db.InitDB()

	defer DB.Close()

	outboxRepo := postgres.NewOutboxRepository(DB)

	kafkaSender := kafka.NewKafka(config.Brokers(), config.Topic())

	sched, err := scheduler.New(ctx, kafkaSender, outboxRepo)
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
