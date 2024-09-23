package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

func Logger(logger ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		logger.Info("Request",
			"status", c.Writer.Status(),
			"method", param.Method,
			"path", path,
			"query", raw,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"latency", param.Latency,
		)
	}
}
