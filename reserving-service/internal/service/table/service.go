package table

import (
	"context"
	"fmt"
	"net/http"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/service/table/mapper"
)

type Service struct {
	tableRepo TableRepository
}

func New(tableRepo TableRepository) *Service {
	return &Service{tableRepo: tableRepo}
}

func (s *Service) GetTablesByRestaurantID(
	ctx context.Context,
	req *dto.GetTablesByRestaurantIDRequest,
) (
	*dto.GetTablesByRestaurantIDResponse,
	int,
	error,
) {
	tables, err := s.tableRepo.GetTablesByRestaurantID(ctx, req.RestaurantID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("get tables by restaurant id(%d): %v", req.RestaurantID, err)
	}

	tablesDTO := mapper.MapToTables(tables)

	return &dto.GetTablesByRestaurantIDResponse{Tables: tablesDTO}, http.StatusOK, nil
}
