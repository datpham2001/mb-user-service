package middlewares

import (
	"github.com/datpham2001/be-winsku/internal/infrastructure/configinfra"
	"github.com/datpham2001/be-winsku/internal/infrastructure/loginfra"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	CONTEXT_REQUEST_ID = "request_id"
	CONTEXT_USER_ID    = "user_id"
	CONTEXT_USER_ROLE  = "user_role"
)

type Middlewares struct {
	Auth *AuthMiddleware
}

func New(auth *AuthMiddleware) *Middlewares {
	return &Middlewares{Auth: auth}
}

func (m *Middlewares) SetupCommon(router *gin.Engine, cfg *configinfra.Config, l *loginfra.Logger) {
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(cfg.Server.ServiceName))
	router.Use(RequestID())
	router.Use(Logger(l))
	router.Use(CORS(cfg.CORS))
	router.Use(Ratelimit())
}
