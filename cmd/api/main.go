package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/datpham2001/mb-user-service/internal/application/usecases"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/authinfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/configinfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/loginfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/persistence/postgresinfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/persistence/postgresinfra/repositories"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/persistence/redisinfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/traceinfra"
	pkgHttp "github.com/datpham2001/mb-user-service/internal/presentation/http"
	"github.com/datpham2001/mb-user-service/internal/presentation/http/controllers"
	"github.com/datpham2001/mb-user-service/internal/presentation/http/middlewares"
	"github.com/datpham2001/mb-user-service/pkg/httpclient"
)

type ServerManager struct {
	HTTPServer *pkgHttp.Server
}

var (
	cfg              *configinfra.Config = &configinfra.Config{}
	logger           *loginfra.Logger
	dbInstance       *postgresinfra.DB
	redisCacheClient *redisinfra.Client
	authInfra        *authinfra.JwtService
	googleOAuthInfra authinfra.IGoogleOAuthService
)

func init() {
	var err error
	if err = configinfra.Load(cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger = loginfra.New(cfg)
	httpClient := httpclient.New()
	authInfra = authinfra.NewJWTService(cfg.JwtAuth)
	googleOAuthInfra = authinfra.NewGoogleOAuthService(cfg.OAuth2.Google, httpClient)

	dbInstance, err = postgresinfra.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	redisCacheClient, err = redisinfra.NewConnection(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}

func initControllers() *controllers.ControllerManager {
	userRepo := repositories.NewUserRepository(dbInstance.GetDB())
	tokenRepo := repositories.NewTokenRepository(dbInstance.GetDB())
	authService := usecases.NewAuthUsecase(userRepo, tokenRepo, authInfra, googleOAuthInfra, cfg.JwtAuth, logger)

	healthController := controllers.NewHealthController(map[string]controllers.Checker{
		"database": dbInstance,
		"redis":    redisCacheClient,
	})
	authController := controllers.NewAuthController(authService)

	controller := controllers.New(
		healthController,
		authController,
	)

	return controller
}

func initMiddlewares() *middlewares.Middlewares {
	authMiddleware := middlewares.NewAuthMiddleware(authInfra)
	return middlewares.New(authMiddleware)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tp, err := traceinfra.InitTracer(cfg.Server.Env, cfg.Server.ServiceName)
	if err != nil {
		logger.Fatalf("failed to initialize tracer: %v", err)
	}

	mw := initMiddlewares()
	controllerManager := initControllers()
	serverManager := &ServerManager{}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		if err = traceinfra.Shutdown(shutdownCtx, tp); err != nil {
			logger.Errorf("failed to shutdown tracer: %v", err)
		}
		if err = serverManager.shutdownHTTPServer(shutdownCtx); err != nil {
			logger.Fatal(err)
		}
		if err = dbInstance.Close(); err != nil {
			logger.Fatalf("failed to close database connection: %v", err)
		}
		if err = redisCacheClient.Close(); err != nil {
			logger.Fatalf("failed to close redis connection: %v", err)
		}

		logger.Info("server shutdown gracefully")
	}()

	go func() {
		logger.Info("starting Winsku server...")
		if err := serverManager.startHTTPServer(ctx, controllerManager, mw); err != nil {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
}

func (s *ServerManager) startHTTPServer(
	ctx context.Context,
	controllerManager *controllers.ControllerManager,
	mw *middlewares.Middlewares,
) error {
	s.HTTPServer = pkgHttp.NewServer(cfg, logger, controllerManager, mw)
	return s.HTTPServer.Start(ctx)
}

func (s *ServerManager) shutdownHTTPServer(ctx context.Context) error {
	if s.HTTPServer == nil {
		return nil
	}

	return s.HTTPServer.Shutdown(ctx)
}
