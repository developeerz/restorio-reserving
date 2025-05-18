package routes

import (
	_ "github.com/developeerz/restorio-reserving/docs" // Импортируем сгенерированную документацию Swagger
	"github.com/developeerz/restorio-reserving/reserving-service/internal/handlers"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/handlers/table"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/middleware"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/scheduler"
	"github.com/gin-gonic/gin" // Необходимо для доступа к файлам Swagger UI
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, h *table.Handler, db *sqlx.DB, scheduler *scheduler.Scheduler) {
	// Сваггер с генерацией по аннотациям
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа для столиков
	tableGroup := r.Group("/tables")
	{
		tableGroup.POST("/new-table", handlers.CreateNewTable(db))
		tableGroup.GET("/:table_id/free-times", handlers.GetFreeTimeSlotsHandler(db))
		tableGroup.GET("", h.GetTablesByRestaurantID)
	}

	// Группа для бронирований
	reservationGroup := r.Group("/reservations")
	{
		reservationGroup.GET("/free-tables", handlers.GetFreeTables(db))
		reservationGroup.POST("/new-reservation", middleware.AuthUserMiddleware, handlers.BookTable(db, scheduler))
		reservationGroup.GET("/user", handlers.GetUserReservations(db))
	}
}
