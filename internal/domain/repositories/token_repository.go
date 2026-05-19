package repositories

import (
	"context"

	"github.com/datpham2001/mb-api-gateway/internal/domain/entities"
)

type ITokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*entities.RefreshToken, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteByHash(ctx context.Context, hash string) error
}
