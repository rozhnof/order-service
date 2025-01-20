package publisher

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.InfoContext(
			c.Request.Context(),
			"incoming request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("address", c.Request.RemoteAddr),
			slog.String("duration", duration.String()),
		)
	}
}
