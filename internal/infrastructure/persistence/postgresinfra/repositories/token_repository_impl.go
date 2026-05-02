package repositories

import (
	"context"
	"errors"

	"github.com/datpham2001/be-winsku/internal/domain/entities"
	"github.com/datpham2001/be-winsku/internal/infrastructure/persistence/postgresinfra/models"
	"gorm.io/gorm"
)

type tokenRepositoryImpl struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *tokenRepositoryImpl {
	return &tokenRepositoryImpl{db: db}
}

func (r *tokenRepositoryImpl) Create(ctx context.Context, token *entities.RefreshToken) error {
	model := &models.RefreshToken{}
	model.FromEntity(token)

	return r.db.WithContext(ctx).Create(model).Error
}

func (r *tokenRepositoryImpl) GetByHash(ctx context.Context, hash string) (*entities.RefreshToken, error) {
	var model models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND (revoke_at IS NULL OR revoke_at > NOW())", hash).
		First(&model).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *tokenRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.RefreshToken{}, "id = ?", id).Error
}

func (r *tokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *tokenRepositoryImpl) DeleteByHash(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).Where("token_hash = ?", hash).Delete(&models.RefreshToken{}).Error
}
