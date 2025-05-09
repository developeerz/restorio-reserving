package http

import (
	_ "github.com/developeerz/restorio-reserving/docs" // swagger docs
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/developeerz/restorio-reserving/internal/usecase"
)

// SetupRoutes регистрирует все HTTP-роуты
func SetupRoutes(
	r *gin.Engine,
	adminRepo port.AdminRepository,
	bookingUC *usecase.BookingUseCase,
	userRepo port.UserRepository,
) {
	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Admin: tables
	adminGroup := r.Group("/tables")
	{
		adminHandler := CreateTableHandler(adminRepo)
		adminGroup.POST("/new-table", adminHandler)
		adminGroup.GET("/:table_id/free-times", GetFreeTimeSlotsHandler(bookingUC))
	}

	// Booking: free tables & new reservation
	bookingGroup := r.Group("/reservations")
	{
		bookingGroup.GET("/free-tables", GetFreeTablesHandler(bookingUC))
		bookingGroup.POST("/new-reservation", BookTableHandler(bookingUC))
	}

	// User: own reservations
	userGroup := r.Group("/reservations")
	{
		userGroup.GET("/user", GetUserReservationsHandler(userRepo))
	}
}
