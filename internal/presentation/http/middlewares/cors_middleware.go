package middlewares

import (
	"github.com/datpham2001/mb-user-service/internal/infrastructure/configinfra"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	MAX_AGE = 12 * 60 * 60 // 12 hours
)

func CORS(cfg configinfra.CORSConfig) gin.HandlerFunc {
	if !cfg.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           MAX_AGE,
	})
}
