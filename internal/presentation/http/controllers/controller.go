package controllers

import (
	"github.com/datpham2001/mb-user-service/internal/presentation/http/middlewares"
	"github.com/gin-gonic/gin"
)

type ControllerRegistry interface {
	RegisterRoutes(router *gin.Engine, mm *middlewares.Middlewares)
}

type ControllerManager struct {
	controllers []ControllerRegistry
}

func New(controllers ...ControllerRegistry) *ControllerManager {
	return &ControllerManager{controllers: controllers}
}

func (cm *ControllerManager) RegisterRoutes(router *gin.Engine, mm *middlewares.Middlewares) {
	for _, controller := range cm.controllers {
		controller.RegisterRoutes(router, mm)
	}
}
