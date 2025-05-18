package middleware

import (
	"net/http"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	var err error
	var authHeader dto.AuthHeader

	if err = ctx.ShouldBindHeader(&authHeader); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Next()
}
