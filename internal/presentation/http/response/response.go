package response

import (
	"errors"
	"net/http"

	domainErrors "github.com/datpham2001/be-winsku/internal/domain/errors"
	pkgValidator "github.com/datpham2001/be-winsku/pkg/validator"
	"github.com/gin-gonic/gin"
)

var domainErrorHTTPStatusMap = map[domainErrors.ErrorCode]int{
	domainErrors.ErrCodeUserNotFound:          http.StatusNotFound,
	domainErrors.ErrCodeEmailAlreadyExists:    http.StatusConflict,
	domainErrors.ErrCodeUsernameAlreadyExists: http.StatusConflict,
	domainErrors.ErrCodeUserAlreadyExists:     http.StatusConflict,
	domainErrors.ErrCodeInvalidCredentials:    http.StatusUnauthorized,
	domainErrors.ErrCodeUserNotActive:         http.StatusUnauthorized,
	domainErrors.ErrCodeUserNotVerified:       http.StatusUnauthorized,
	domainErrors.ErrCodeInvalidEmail:          http.StatusBadRequest,
	domainErrors.ErrCodeInvalidPassword:       http.StatusBadRequest,

	domainErrors.ErrCodeOAuthInvalidCode:         http.StatusUnauthorized,
	domainErrors.ErrCodeOAuthFailedToGetUserInfo: http.StatusInternalServerError,
	domainErrors.ErrCodeOAuthEmailNotVerified:    http.StatusForbidden,
}

type Response struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string                    `json:"code"`
	Message string                    `json:"message"`
	Details []pkgValidator.FieldError `json:"details,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func ErrorWithDetails(c *gin.Context, statusCode int, code string, message string, details []pkgValidator.FieldError) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}

func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data)
}

func OK(c *gin.Context, message string, data any) {
	Success(c, http.StatusOK, message, data)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "CONFLICT", message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func ValidationError(c *gin.Context, err error) {
	fields := pkgValidator.TranslateValidationErrors(err)

	code := "VALIDATION_ERROR"
	message := "Validation failed"
	if len(fields) == 0 {
		if err != nil {
			message = err.Error()
		}
		Error(c, http.StatusBadRequest, code, message)
		return
	}

	ErrorWithDetails(c, http.StatusBadRequest, code, message, fields)
}

func UsecaseError(c *gin.Context, err error) {
	var de *domainErrors.DomainError
	if !errors.As(err, &de) {
		InternalServerError(c, "Internal server error")
		return
	}

	switch de.Type {
	case domainErrors.ClientError:
		statusCode := domainErrorHTTPStatus(de.Code)
		Error(c, statusCode, string(de.Code), de.Message)
	default:
		InternalServerError(c, "Internal server error")
	}
}

func domainErrorHTTPStatus(code domainErrors.ErrorCode) int {
	if status, ok := domainErrorHTTPStatusMap[code]; ok {
		return status
	}
	return http.StatusBadRequest
}
