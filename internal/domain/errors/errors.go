package errors

import "fmt"

type ErrorType int

const (
	ClientError ErrorType = iota + 1
	ServerError
)

type ErrorCode string

const (
	ErrCodeUserNotFound          ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists     ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeEmailAlreadyExists    ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrCodeUsernameAlreadyExists ErrorCode = "USERNAME_ALREADY_EXISTS"
	ErrCodeInvalidCredentials    ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeUserNotActive         ErrorCode = "USER_NOT_ACTIVE"
	ErrCodeUserNotVerified       ErrorCode = "USER_NOT_VERIFIED"
	ErrCodeInvalidEmail          ErrorCode = "INVALID_EMAIL"
	ErrCodeInvalidPassword       ErrorCode = "INVALID_PASSWORD"
	ErrCodeInternal              ErrorCode = "INTERNAL_SERVER_ERROR"

	ErrCodeOAuthInvalidCode         ErrorCode = "OAUTH_INVALID_CODE"
	ErrCodeOAuthFailedToGetUserInfo ErrorCode = "OAUTH_FAILED_TO_GET_USER_INFO"
	ErrCodeOAuthEmailNotVerified    ErrorCode = "OAUTH_EMAIL_NOT_VERIFIED"
)

type DomainError struct {
	Type    ErrorType
	Code    ErrorCode
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) WithMessage(message string) *DomainError {
	return &DomainError{Type: e.Type, Code: e.Code, Message: message}
}

func (e *DomainError) WithMessagef(format string, args ...any) *DomainError {
	return &DomainError{Type: e.Type, Code: e.Code, Message: fmt.Sprintf(format, args...)}
}
