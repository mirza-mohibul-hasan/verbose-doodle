package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a Gin middleware that logs each request using structured slog.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Build log attributes
		attrs := []slog.Attr{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Duration("latency", latency),
			slog.Int("body_size", c.Writer.Size()),
		}

		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		// Convert to []any for slog
		args := make([]any, len(attrs))
		for i, attr := range attrs {
			args[i] = attr
		}

		// Log based on status code
		status := c.Writer.Status()
		switch {
		case status >= 500:
			slog.Error("Server error", args...)
		case status >= 400:
			slog.Warn("Client error", args...)
		default:
			slog.Info("Request", args...)
		}
	}
}
