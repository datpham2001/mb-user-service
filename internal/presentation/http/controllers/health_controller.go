package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/datpham2001/mb-user-service/internal/presentation/http/middlewares"
	"github.com/gin-gonic/gin"
)

type Checker interface {
	HealthCheck() error
}

type componentStatus struct {
	Status    string  `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	Error     *string `json:"error,omitempty"`
}

type readinessResponse struct {
	Status     string                     `json:"status"`
	Timestamp  string                     `json:"timestamp"`
	UptimeMs   float64                    `json:"uptime_ms"`
	Components map[string]componentStatus `json:"components"`
}

type HealthController struct {
	startTime time.Time
	checkers  map[string]Checker
}

func NewHealthController(checkers map[string]Checker) *HealthController {
	return &HealthController{
		startTime: time.Now(),
		checkers:  checkers,
	}
}

func (h *HealthController) RegisterRoutes(router *gin.Engine, _ *middlewares.Middlewares) {
	router.GET("/api/health", h.liveness)
	router.GET("/api/health/ready", h.readiness)
}

func (h *HealthController) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"uptime_ms": float64(time.Since(h.startTime).Milliseconds()),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthController) readiness(c *gin.Context) {
	const timeout = 3 * time.Second

	ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
	defer cancel()

	type result struct {
		name   string
		status componentStatus
	}

	results := make(chan result, len(h.checkers))
	var wg sync.WaitGroup

	for name, checker := range h.checkers {
		wg.Add(1)
		go func(name string, checker Checker) {
			defer wg.Done()

			start := time.Now()
			err := runWithContext(ctx, checker.HealthCheck)
			latency := float64(time.Since(start).Milliseconds())

			cs := componentStatus{
				Status:    "healthy",
				LatencyMs: latency,
			}
			if err != nil {
				msg := err.Error()
				cs.Status = "unhealthy"
				cs.Error = &msg
			}

			results <- result{name: name, status: cs}
		}(name, checker)
	}

	wg.Wait()
	close(results)

	components := make(map[string]componentStatus, len(h.checkers))
	overall := "healthy"
	for r := range results {
		components[r.name] = r.status
		if r.status.Status != "healthy" {
			overall = "unhealthy"
		}
	}

	statusCode := http.StatusOK
	if overall == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, readinessResponse{
		Status:     overall,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		UptimeMs:   float64(time.Since(h.startTime).Milliseconds()),
		Components: components,
	})
}

func runWithContext(ctx context.Context, fn func() error) error {
	done := make(chan error, 1)
	go func() { done <- fn() }()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
