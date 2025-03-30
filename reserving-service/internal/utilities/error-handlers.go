package utilities

import (
	"github.com/gin-gonic/gin"
)

// JSONError отправляет ответ с указанным HTTP-статусом и сообщением об ошибке,
// а затем прерывает дальнейшую обработку запроса.
func JSONError(c *gin.Context, status int, errMsg string) {
	c.JSON(status, gin.H{"error": errMsg})
	c.Abort()
}
