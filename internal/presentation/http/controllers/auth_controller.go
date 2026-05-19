package controllers

import (
	"github.com/datpham2001/mb-api-gateway/internal/application/usecases"
	"github.com/datpham2001/mb-api-gateway/internal/presentation/http/middlewares"
	"github.com/datpham2001/mb-api-gateway/internal/presentation/http/request"
	"github.com/datpham2001/mb-api-gateway/internal/presentation/http/response"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService usecases.IAuthUsecase
}

func NewAuthController(authService usecases.IAuthUsecase) *AuthController {
	return &AuthController{authService: authService}
}

func (h *AuthController) RegisterRoutes(router *gin.Engine, mm *middlewares.Middlewares) {
	authV1 := router.Group("/api/v1/auth")
	{
		authV1.POST("/register", h.RegisterAccount)
		authV1.POST("/login", h.LoginAccount)
		authV1.POST("/token/refresh", h.RefreshToken)
		authV1.POST("/oauth/google/callback", h.GoogleCallback)
	}

	protectedAuthV1 := authV1.Group("")
	protectedAuthV1.Use(mm.Auth.Handle())
	{
		protectedAuthV1.POST("/logout", h.LogoutAccount)
	}
}

func (h *AuthController) RegisterAccount(c *gin.Context) {
	var req request.RegisterAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.authService.RegisterAccount(c.Request.Context(), req.ToDTO())
	if err != nil {
		response.UsecaseError(c, err)
		return
	}

	response.Created(c, "Account registered successfully", res)
}

func (h *AuthController) LoginAccount(c *gin.Context) {
	var req request.LoginAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.authService.LoginAccount(c.Request.Context(), req.ToDTO())
	if err != nil {
		response.UsecaseError(c, err)
		return
	}

	response.OK(c, "Login successfully", res)
}

func (h *AuthController) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.authService.RefreshToken(c.Request.Context(), req.ToDTO())
	if err != nil {
		response.UsecaseError(c, err)
		return
	}

	response.OK(c, "Token refreshed successfully", res)
}

func (h *AuthController) LogoutAccount(c *gin.Context) {
	var req request.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.authService.Logout(c.Request.Context(), req.ToDTO())
	if err != nil {
		response.UsecaseError(c, err)
		return
	}

	response.OK(c, "Logged out successfully", nil)
}

func (h *AuthController) GoogleCallback(c *gin.Context) {
	var req request.GoogleCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.authService.GoogleOAuthLogin(c.Request.Context(), req.ToDTO())
	if err != nil {
		response.UsecaseError(c, err)
		return
	}

	if res == nil {
		response.InternalServerError(c, "Internal server error")
		return
	}

	response.OK(c, "Google OAuth login successfully", res)
}
