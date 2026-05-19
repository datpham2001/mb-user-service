package models

import (
	"time"

	"github.com/datpham2001/mb-api-gateway/internal/domain/entities"
)

type RefreshToken struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string     `gorm:"user_id;not null"`
	TokenHash string     `gorm:"token_hash;not null"`
	ExpiresAt *time.Time `gorm:"expires_at;not null"`
	RevokeAt  *time.Time `gorm:"revoke_at"`
	CreatedAt time.Time  `gorm:"created_at;autoCreateTime"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (rt *RefreshToken) ToEntity() *entities.RefreshToken {
	return &entities.RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		RevokeAt:  rt.RevokeAt,
		CreatedAt: rt.CreatedAt,
	}
}

func (rt *RefreshToken) FromEntity(e *entities.RefreshToken) {
	rt.ID = e.ID
	rt.UserID = e.UserID
	rt.TokenHash = e.TokenHash
	rt.ExpiresAt = e.ExpiresAt
	rt.RevokeAt = e.RevokeAt
	rt.CreatedAt = e.CreatedAt
}
