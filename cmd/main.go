package main

import (
	"context"
	"log"

	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/internal/adapter/http"
	"github.com/developeerz/restorio-reserving/internal/adapter/kafka"
	"github.com/developeerz/restorio-reserving/internal/adapter/postgres"
	"github.com/developeerz/restorio-reserving/internal/adapter/scheduler"
	"github.com/developeerz/restorio-reserving/internal/db"
	"github.com/developeerz/restorio-reserving/internal/usecase"
	"github.com/gin-gonic/gin" // Необходимо для доступа к файлам Swagger UI
)

// @title Restorio API
// @version 1.0
// @description API для бронирования столиков
// @host localhost:8082
// @BasePath /
func main() {
	// 1) Конфиг, DB, Kafka, Outbox
	dbConn := db.InitDB()
	defer dbConn.Close()
	kafkaSender := kafka.NewKafka(nil) //[]string{"broker1:9092"}
	outboxRepo := postgres.NewOutboxRepo(dbConn)

	// 2) Scheduler
	sched, err := scheduler.New(context.Background(), kafkaSender, outboxRepo)
	if err != nil {
		log.Fatalf("scheduler init error: %v", err)
	}
	defer sched.Stop()

	// 3) Репозитории
	adminRepo := postgres.NewPostgresAdminRepository(dbConn) // port.AdminRepository
	bookingRepo := postgres.NewBookingRepo(dbConn)           // port.BookingRepository
	reservationRepo := postgres.NewUserRepo(dbConn)          // port.ReservationRepository

	// 4) UseCase
	bookingUC := usecase.NewBookingUseCase(bookingRepo, sched)

	// 5) HTTP
	r := gin.Default()
	http.SetupRoutes(r, adminRepo, bookingUC, reservationRepo)

	// 6) Запускаем сервер
	log.Println("Сервер запущен на порту 8082")
	r.Run(":8082")
}
