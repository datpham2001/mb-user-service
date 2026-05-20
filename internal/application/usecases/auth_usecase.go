package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
	"unicode"

	"github.com/datpham2001/mb-user-service/internal/application/dto"
	"github.com/datpham2001/mb-user-service/internal/domain/entities"
	domainErrors "github.com/datpham2001/mb-user-service/internal/domain/errors"
	"github.com/datpham2001/mb-user-service/internal/domain/repositories"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/authinfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/configinfra"
	"github.com/datpham2001/mb-user-service/internal/infrastructure/loginfra"
	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST = 12
)

type IAuthUsecase interface {
	RegisterAccount(ctx context.Context, req *dto.RegisterAccountRequest) (*dto.RegisterAccountResponse, error)
	LoginAccount(ctx context.Context, req *dto.LoginAccountRequest) (*dto.LoginAccountResponse, error)
	GoogleOAuthLogin(ctx context.Context, req *dto.GoogleOAuthRequest) (*dto.LoginAccountResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginAccountResponse, error)
	Logout(ctx context.Context, req *dto.LogoutRequest) error
}

type authUsecase struct {
	userRepo    repositories.IUserRepository
	tokenRepo   repositories.ITokenRepository
	jwtSvc      authinfra.IJWTService
	googleOAuth authinfra.IGoogleOAuthService
	cfg         configinfra.JwtAuthConfig
	logger      *loginfra.Logger
}

func NewAuthUsecase(
	userRepo repositories.IUserRepository,
	tokenRepo repositories.ITokenRepository,
	jwtSvc authinfra.IJWTService,
	googleOAuth authinfra.IGoogleOAuthService,
	cfg configinfra.JwtAuthConfig,
	logger *loginfra.Logger,
) *authUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		jwtSvc:      jwtSvc,
		googleOAuth: googleOAuth,
		cfg:         cfg,
		logger:      logger,
	}
}

func (s *authUsecase) RegisterAccount(ctx context.Context, req *dto.RegisterAccountRequest) (*dto.RegisterAccountResponse, error) {
	emailExists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("failed to check email existence: %v", err)
		return nil, domainErrors.ErrInternal
	}
	if emailExists {
		return nil, domainErrors.ErrEmailAlreadyExists
	}

	usernameExists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, domainErrors.ErrInternal
	}
	if usernameExists {
		return nil, domainErrors.ErrUsernameAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), BCRYPT_COST)
	if err != nil {
		return nil, domainErrors.ErrInternal
	}

	user := entities.NewUserForRegistration(
		req.Email,
		req.Username,
		string(hashedPassword),
		req.FirstName,
		req.LastName,
	)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, domainErrors.ErrInternal
	}

	return dto.NewRegisterAccountResponse(user), nil
}

func (s *authUsecase) LoginAccount(ctx context.Context, req *dto.LoginAccountRequest) (*dto.LoginAccountResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domainErrors.ErrInternal
	}
	if user == nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domainErrors.ErrUserNotActive
	}

	return s.generateTokenPair(ctx, user, nil, req.RememberMe)
}

func (s *authUsecase) GoogleOAuthLogin(ctx context.Context, req *dto.GoogleOAuthRequest) (*dto.LoginAccountResponse, error) {
	accessToken, err := s.googleOAuth.ExchangeCode(ctx, req.Code, req.RedirectUri)
	if err != nil {
		s.logger.Errorf("google oauth exchange failed: %v", err)
		return nil, domainErrors.ErrOAuthInvalidCode
	}

	googleUser, err := s.googleOAuth.GetUserInfo(ctx, accessToken)
	if err != nil {
		s.logger.Errorf("google oauth get user info failed: %v", err)
		return nil, domainErrors.ErrOAuthFailedToGetUserInfo
	}

	if !googleUser.EmailVerified {
		return nil, domainErrors.ErrOAuthEmailNotVerified
	}

	user, err := s.userRepo.GetByEmail(ctx, googleUser.Email)
	if err != nil {
		s.logger.Errorf("failed to fetch user by email: %v", err)
		return nil, domainErrors.ErrInternal
	}

	if user == nil {
		username := s.buildOAuthUsername(googleUser.Email, googleUser.Sub)
		user = entities.NewUserFromOAuth(
			googleUser.Email,
			username,
			googleUser.GivenName,
			googleUser.FamilyName,
			googleUser.Picture,
			"google",
			googleUser.Sub,
		)
		if err := s.userRepo.Create(ctx, user); err != nil {
			s.logger.Errorf("failed to create oauth user: %v", err)
			return nil, domainErrors.ErrInternal
		}
	} else {
		now := time.Now()
		user.LastLoginAt = &now
		if err := s.userRepo.Update(ctx, user); err != nil {
			s.logger.Errorf("failed to update last_login_at: %v", err)
			return nil, domainErrors.ErrInternal
		}
	}

	if !user.IsActive {
		return nil, domainErrors.ErrUserNotActive
	}

	return s.generateTokenPair(ctx, user, nil, req.RememberMe)
}

func (s *authUsecase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginAccountResponse, error) {
	hash := s.hashToken(req.RefreshToken)
	tokenEntity, err := s.tokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil, domainErrors.ErrInternal
	}
	if tokenEntity == nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	_ = s.tokenRepo.DeleteByHash(ctx, hash)

	user, err := s.userRepo.GetByID(ctx, tokenEntity.UserID)
	if err != nil {
		return nil, domainErrors.ErrInternal
	}
	if user == nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	return s.generateTokenPair(ctx, user, tokenEntity.ExpiresAt, false)
}

func (s *authUsecase) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	hash := s.hashToken(req.RefreshToken)
	return s.tokenRepo.DeleteByHash(ctx, hash)
}

func (s *authUsecase) generateTokenPair(
	ctx context.Context,
	user *entities.User,
	customExpiresAt *time.Time,
	rememberMe bool,
) (*dto.LoginAccountResponse, error) {
	accessToken, err := s.jwtSvc.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, domainErrors.ErrInternal
	}

	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID, customExpiresAt)
	if err != nil {
		return nil, domainErrors.ErrInternal
	}

	hash := s.hashToken(refreshToken)

	var expiresAt time.Time
	switch {
	case customExpiresAt != nil:
		expiresAt = *customExpiresAt
	case rememberMe:
		expiresAt = time.Now().Add(s.cfg.RememberMeTokenExp)
	default:
		expiresAt = time.Now().Add(s.cfg.RefreshTokenExp)
	}

	err = s.tokenRepo.Create(ctx, &entities.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: &expiresAt,
	})
	if err != nil {
		return nil, domainErrors.ErrInternal
	}

	return &dto.LoginAccountResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (s *authUsecase) buildOAuthUsername(email, sub string) string {
	prefix := strings.Split(email, "@")[0]

	prefix = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, prefix)

	if len(prefix) > 20 {
		prefix = prefix[:20]
	}

	suffix := sub
	if len(suffix) > 6 {
		suffix = suffix[:6]
	}

	return prefix + "_" + suffix
}

func (s *authUsecase) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
