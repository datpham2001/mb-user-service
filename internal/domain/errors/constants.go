package errors

var (
	ErrUserNotFound          = &DomainError{Type: ClientError, Code: ErrCodeUserNotFound, Message: "User not found"}
	ErrUserAlreadyExists     = &DomainError{Type: ClientError, Code: ErrCodeUserAlreadyExists, Message: "User already exists"}
	ErrEmailAlreadyExists    = &DomainError{Type: ClientError, Code: ErrCodeEmailAlreadyExists, Message: "Email already exists"}
	ErrUsernameAlreadyExists = &DomainError{Type: ClientError, Code: ErrCodeUsernameAlreadyExists, Message: "Username already exists"}
	ErrInvalidCredentials    = &DomainError{Type: ClientError, Code: ErrCodeInvalidCredentials, Message: "Invalid credentials"}
	ErrUserNotActive         = &DomainError{Type: ClientError, Code: ErrCodeUserNotActive, Message: "User is not active"}
	ErrUserNotVerified       = &DomainError{Type: ClientError, Code: ErrCodeUserNotVerified, Message: "User is not verified"}
	ErrInvalidEmail          = &DomainError{Type: ClientError, Code: ErrCodeInvalidEmail, Message: "Invalid email format"}
	ErrInvalidPassword       = &DomainError{Type: ClientError, Code: ErrCodeInvalidPassword, Message: "Password is invalid"}
)

var (
	ErrInternal = &DomainError{Type: ServerError, Code: ErrCodeInternal, Message: "Internal server error"}
)

var (
	ErrOAuthInvalidCode         = &DomainError{Type: ClientError, Code: ErrCodeOAuthInvalidCode, Message: "Invalid or expired OAuth code"}
	ErrOAuthFailedToGetUserInfo = &DomainError{Type: ServerError, Code: ErrCodeOAuthFailedToGetUserInfo, Message: "Failed to get user info from OAuth provider"}
	ErrOAuthEmailNotVerified    = &DomainError{Type: ClientError, Code: ErrCodeOAuthEmailNotVerified, Message: "Google email is not verified"}
)
