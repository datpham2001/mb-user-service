package models

import (
	"time"

	"github.com/datpham2001/mb-user-service/internal/domain/entities"
)

type User struct {
	ID             string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email          string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username       string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	HashedPassword *string    `gorm:"type:varchar(255)"` // nullable — OAuth users have no password
	FirstName      string     `gorm:"type:varchar(100)"`
	LastName       string     `gorm:"type:varchar(100)"`
	Bio            string     `gorm:"type:text"`
	AvatarURL      string     `gorm:"type:varchar(500)"`
	IsActive       bool       `gorm:"not null;default:true"`
	Role           string     `gorm:"type:varchar(20);not null;default:'user'"`
	AuthProvider   string     `gorm:"type:varchar(50);not null;default:'local'"`
	ProviderID     string     `gorm:"type:varchar(255);not null;default:''"`
	LastLoginAt    *time.Time `gorm:"column:last_login_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToEntity() *entities.User {
	var hashedPassword string
	if u.HashedPassword != nil {
		hashedPassword = *u.HashedPassword
	}

	return &entities.User{
		ID:             u.ID,
		Email:          u.Email,
		Username:       u.Username,
		HashedPassword: hashedPassword,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Bio:            u.Bio,
		AvatarURL:      u.AvatarURL,
		IsActive:       u.IsActive,
		Role:           entities.RoleType(u.Role),
		AuthProvider:   u.AuthProvider,
		ProviderID:     u.ProviderID,
		LastLoginAt:    u.LastLoginAt,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func (u *User) ToDB(entity *entities.User) {
	var hashedPassword *string
	if entity.HashedPassword != "" {
		hashedPassword = &entity.HashedPassword
	}

	u.ID = entity.ID
	u.Email = entity.Email
	u.Username = entity.Username
	u.HashedPassword = hashedPassword
	u.FirstName = entity.FirstName
	u.LastName = entity.LastName
	u.Bio = entity.Bio
	u.AvatarURL = entity.AvatarURL
	u.IsActive = entity.IsActive
	u.Role = string(entity.Role)
	u.AuthProvider = entity.AuthProvider
	u.ProviderID = entity.ProviderID
	u.LastLoginAt = entity.LastLoginAt
	u.CreatedAt = entity.CreatedAt
	u.UpdatedAt = entity.UpdatedAt
}
