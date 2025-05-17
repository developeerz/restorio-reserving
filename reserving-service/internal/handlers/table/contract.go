package table

import (
	"context"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
)

type Service interface {
	GetTablesByRestaurantID(
		ctx context.Context,
		req *dto.GetTablesByRestaurantIDRequest,
	) (*dto.GetTablesByRestaurantIDResponse, int, error)
}
