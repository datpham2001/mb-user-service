package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Ratelimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
