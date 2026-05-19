package middlewares

import (
	"net/http"
	"time"

	"github.com/datpham2001/mb-api-gateway/internal/infrastructure/loginfra"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(l *loginfra.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		latency := time.Since(start)

		entry := l.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
		})

		if requestID := c.GetString(CONTEXT_REQUEST_ID); requestID != "" {
			entry = entry.WithField(CONTEXT_REQUEST_ID, requestID)
		}

		if userID := c.GetString(CONTEXT_USER_ID); userID != "" {
			entry = entry.WithField(CONTEXT_USER_ID, userID)
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			entry.Error("HTTP Request")
		} else if c.Writer.Status() >= http.StatusBadRequest {
			entry.Warn("HTTP Request")
		} else {
			entry.Info("HTTP Request")
		}
	}
}
