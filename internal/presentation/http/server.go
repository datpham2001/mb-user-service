package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/datpham2001/mb-api-gateway/internal/infrastructure/configinfra"
	"github.com/datpham2001/mb-api-gateway/internal/infrastructure/loginfra"
	"github.com/datpham2001/mb-api-gateway/internal/presentation/http/controllers"
	"github.com/datpham2001/mb-api-gateway/internal/presentation/http/middlewares"
	"github.com/gin-gonic/gin"
)

const (
	SHUTDOWN_TIMEOUT = 5 * time.Second
)

type Server struct {
	router            *gin.Engine
	cfg               *configinfra.Config
	server            *http.Server
	logger            *loginfra.Logger
	controllerManager *controllers.ControllerManager
	middlewares       *middlewares.Middlewares
}

func NewServer(
	cfg *configinfra.Config,
	logger *loginfra.Logger,
	controllerManager *controllers.ControllerManager,
	middlewares *middlewares.Middlewares,
) *Server {
	if cfg.Server.Env == "local" || cfg.Server.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		cfg:               cfg,
		logger:            logger,
		router:            gin.New(),
		controllerManager: controllerManager,
		middlewares:       middlewares,
	}

	s.setupMiddlewares()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddlewares() {
	s.middlewares.SetupCommon(s.router, s.cfg, s.logger)
}

func (s *Server) setupRoutes() {
	s.controllerManager.RegisterRoutes(s.router, s.middlewares)
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler: s.router,
	}
	s.server.RegisterOnShutdown(cancel)

	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	s.logger.Infof("server is listening on %s", addr)

	go func() {
		var err error
		if s.cfg.Server.TLS.Enable {
			err = s.server.ListenAndServeTLS(s.cfg.Server.TLS.CertFile, s.cfg.Server.TLS.KeyFile)
		} else {
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("failed to start http server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		if closeErr := s.server.Close(); closeErr != nil {
			return fmt.Errorf("could not stop server gracefully and force to close: %w", closeErr)
		}

		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	return nil
}
