package authinfra

import (
	"context"
	"net/url"

	"github.com/datpham2001/mb-user-service/internal/infrastructure/configinfra"
	"github.com/datpham2001/mb-user-service/pkg/httpclient"
)

const (
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

type IGoogleOAuthService interface {
	ExchangeCode(ctx context.Context, code, redirectUri string) (string, error)
	GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error)
}

type googleOAuthService struct {
	cfg    configinfra.OAuth2ProviderConfig
	client *httpclient.Client
}

func NewGoogleOAuthService(cfg configinfra.OAuth2ProviderConfig, client *httpclient.Client) IGoogleOAuthService {
	return &googleOAuthService{
		cfg:    cfg,
		client: client,
	}
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

func (s *googleOAuthService) ExchangeCode(ctx context.Context, code, redirectUri string) (string, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {s.cfg.ClientID},
		"client_secret": {s.cfg.ClientSecret},
		"redirect_uri":  {redirectUri},
		"grant_type":    {"authorization_code"},
	}

	var tokenResp googleTokenResponse
	if err := s.client.PostForm(ctx, googleTokenURL, form, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (s *googleOAuthService) GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	var userInfo GoogleUserInfo
	if err := s.client.Get(
		ctx, googleUserInfoURL, &userInfo, httpclient.WithBearerToken(accessToken),
	); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
