package middleware

import (
	"net/http"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/gin-gonic/gin"
)

func authMiddleware(ctx *gin.Context, auth string) {
	var err error
	var authHeader dto.AuthHeader

	if err = ctx.ShouldBindHeader(&authHeader); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if authHeader.Auths != auth {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Set("userID", authHeader.UserID)

	ctx.Next()
}

func AuthUserMiddleware(ctx *gin.Context) {
	authMiddleware(ctx, "USER")
	ctx.Next()
}

func AuthAdminMiddleware(ctx *gin.Context) {
	authMiddleware(ctx, "ADMIN")
	ctx.Next()
}
