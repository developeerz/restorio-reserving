package routes

import (
	"reserving-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRoutes(r *gin.Engine, db *sqlx.DB) {
	// Группа для столиков
	tableGroup := r.Group("/tables")
	{
		tableGroup.POST("/", handlers.CreateTable(db))
		tableGroup.PUT("/:table_id", handlers.UpdateTablePosition(db))
		tableGroup.DELETE("/:table_id", handlers.DeleteTable(db))
	}
	
	// Группа для бронирований
	reservationGroup := r.Group("/reservations")
	{
		reservationGroup.GET("/user", handlers.GetUserReservations(db))
	}
}
