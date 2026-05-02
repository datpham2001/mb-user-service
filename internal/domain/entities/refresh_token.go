package entities

import "time"

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt *time.Time
	RevokeAt  *time.Time
	CreatedAt time.Time
}
