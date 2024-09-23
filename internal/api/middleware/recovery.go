package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

func Recovery(logger ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic occurred", "error", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
