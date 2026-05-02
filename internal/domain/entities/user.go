package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoleType string

const (
	UserRole  RoleType = "user"
	AdminRole RoleType = "admin"
)

type User struct {
	ID             string
	Email          string
	Username       string
	HashedPassword string
	FirstName      string
	LastName       string
	Bio            string
	AvatarURL      string
	IsActive       bool
	Role           RoleType
	AuthProvider   string
	ProviderID     string
	LastLoginAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUserForRegistration(
	email string, username string, password string, firstName string, lastName string,
) *User {
	return &User{
		ID:             uuid.NewString(),
		Email:          email,
		Username:       username,
		HashedPassword: password,
		FirstName:      firstName,
		LastName:       lastName,
		IsActive:       true,
		Role:           UserRole,
		AuthProvider:   "local",
	}
}

func NewUserFromOAuth(
	email string, username string, firstName string, lastName string,
	avatarURL string, authProvider string, providerID string,
) *User {
	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		AvatarURL:    avatarURL,
		IsActive:     true,
		Role:         UserRole,
		AuthProvider: authProvider,
		ProviderID:   providerID,
	}
}
