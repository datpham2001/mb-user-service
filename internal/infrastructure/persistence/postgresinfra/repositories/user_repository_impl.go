package repositories

import (
	"context"
	"errors"

	"github.com/datpham2001/be-winsku/internal/domain/entities"
	"github.com/datpham2001/be-winsku/internal/infrastructure/persistence/postgresinfra/models"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	model := &models.User{}
	model.ToDB(user)

	if err := u.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (u *userRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var model models.User

	if err := u.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (u *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model models.User

	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (u *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var model models.User

	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (u *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	model := &models.User{}
	model.ToDB(user)

	result := u.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return result.Error
	}

	user.UpdatedAt = model.UpdatedAt
	return nil
}

func (u *userRepositoryImpl) Delete(ctx context.Context, id string) error {
	return u.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (u *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *userRepositoryImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
