package middlewares

import (
	"strings"
	"time"

	"slices"

	"github.com/datpham2001/mb-user-service/internal/infrastructure/authinfra"
	"github.com/datpham2001/mb-user-service/internal/presentation/http/response"
	"github.com/gin-gonic/gin"
)

const (
	BEARER_TOKEN_PREFIX  = "Bearer"
	HEADER_AUTHORIZATION = "Authorization"
)

type AuthMiddleware struct {
	jwtSvc authinfra.IJWTService
}

func NewAuthMiddleware(jwtSvc authinfra.IJWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtSvc: jwtSvc}
}

func (am *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(HEADER_AUTHORIZATION)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == BEARER_TOKEN_PREFIX) {
			response.Unauthorized(c, "Authorization header must be Bearer token")
			c.Abort()
			return
		}

		claims, err := am.jwtSvc.ValidateAccessToken(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			response.Unauthorized(c, "Token is expired")
			c.Abort()
			return
		}

		c.Set(CONTEXT_USER_ID, claims.UserID)
		c.Set(CONTEXT_USER_ROLE, claims.Role)
		c.Next()
	}
}

func (am *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString(CONTEXT_USER_ROLE)
		if slices.Contains(roles, userRole) {
			c.Next()
			return
		}

		response.Forbidden(c, "You don't have permission to perform this action")
		c.Abort()
	}
}
