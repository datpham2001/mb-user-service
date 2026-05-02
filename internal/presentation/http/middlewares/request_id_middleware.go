package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	HEADER_X_REQUEST_ID = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(HEADER_X_REQUEST_ID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header(HEADER_X_REQUEST_ID, requestID)
		c.Set(CONTEXT_REQUEST_ID, requestID)

		c.Next()
	}
}
