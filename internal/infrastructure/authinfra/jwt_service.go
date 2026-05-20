package authinfra

import (
	"fmt"
	"time"

	"github.com/datpham2001/mb-user-service/internal/infrastructure/configinfra"
	"github.com/golang-jwt/jwt/v5"
)

type IJWTService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	ValidateAccessToken(token string) (*Claims, error)
	GenerateRefreshToken(userID string, customExpiresAt *time.Time) (string, error)
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JwtService struct {
	cfg configinfra.JwtAuthConfig
}

func NewJWTService(cfg configinfra.JwtAuthConfig) *JwtService {
	return &JwtService{cfg: cfg}
}

func (s *JwtService) GenerateAccessToken(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "winsku-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}

func (s *JwtService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *JwtService) GenerateRefreshToken(userID string, customExpiresAt *time.Time) (string, error) {
	exp := time.Now().Add(s.cfg.RefreshTokenExp)
	if customExpiresAt != nil {
		exp = *customExpiresAt
	}

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "winsku-api",
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}
