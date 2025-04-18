package utilities

import (
	"github.com/gin-gonic/gin"
)

/*
 * JSONError отправляет ответ с указанным HTTP-статусом и сообщением об ошибке.
 * Возвращает true отправке ошибки.
 */
func JSONError(c *gin.Context, err error, status int, errMsg string, details ...any) bool {
	if err != nil {
		response := gin.H{"error": errMsg}

		// Если переданы дополнительные аргументы, добавляем их в ответ
		if len(details) > 0 {
			response["details"] = details
		}

		c.JSON(status, response)
		return true
	}
	return false
}
