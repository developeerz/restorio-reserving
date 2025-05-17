package table

import (
	"context"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/table"
)

type TableRepository interface {
	GetTablesByRestaurantID(ctx context.Context, restaurantID int) ([]table.Table, error)
}
