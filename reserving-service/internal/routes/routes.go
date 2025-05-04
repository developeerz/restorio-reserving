package routes

import (
	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/reserving-service/internal/handlers"
	"github.com/gin-gonic/gin" // Необходимо для доступа к файлам Swagger UI
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, db *sqlx.DB) {
	// Сваггер с генерацией по аннотациям
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа для столиков
	tableGroup := r.Group("/tables")
	{
		tableGroup.POST("/new-table", handlers.CreateNewTable(db))
		tableGroup.GET("/:table_id/free-times", handlers.GetFreeTimeSlotsHandler(db))
	}

	// Группа для бронирований
	reservationGroup := r.Group("/reservations")
	{
		reservationGroup.GET("/free-tables", handlers.GetFreeTables(db))
		reservationGroup.POST("/new-reservation", handlers.BookTable(db))
		reservationGroup.GET("/user", handlers.GetUserReservations(db))
	}
}
