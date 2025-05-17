package table

import (
	"net/http"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// @Summary Получение столиков по id ресторана
// @Tags Tables
// @Produce  json
// @Param restaurant_id query int true "ID ресторана" example(1)
// @Success 200 {object} dto.GetTablesByRestaurantIDResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tables [get]
func (h *Handler) GetTablesByRestaurantID(ctx *gin.Context) {
	var err error
	var req dto.GetTablesByRestaurantIDRequest

	if err = ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "could not bind query"})
		return
	}

	res, code, err := h.service.GetTablesByRestaurantID(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(code, &dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
