package request

import "github.com/datpham2001/mb-user-service/internal/application/dto"

type RegisterAccountReq struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (r *RegisterAccountReq) ToDTO() *dto.RegisterAccountRequest {
	if r == nil {
		return nil
	}

	return &dto.RegisterAccountRequest{
		Email:    r.Email,
		Username: r.Username,
		Password: r.Password,
	}
}

type LoginAccountReq struct {
	Email      string `json:"email"    binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8,max=72"`
	RememberMe bool   `json:"remember_me"`
}

func (r *LoginAccountReq) ToDTO() *dto.LoginAccountRequest {
	if r == nil {
		return nil
	}

	return &dto.LoginAccountRequest{
		Email:      r.Email,
		Password:   r.Password,
		RememberMe: r.RememberMe,
	}
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *RefreshTokenReq) ToDTO() *dto.RefreshTokenRequest {
	if r == nil {
		return nil
	}
	return &dto.RefreshTokenRequest{
		RefreshToken: r.RefreshToken,
	}
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *LogoutReq) ToDTO() *dto.LogoutRequest {
	if r == nil {
		return nil
	}
	return &dto.LogoutRequest{
		RefreshToken: r.RefreshToken,
	}
}

type GoogleCallbackReq struct {
	Code        string `json:"code" binding:"required"`
	RedirectUri string `json:"redirect_uri" binding:"required,url"`
	RememberMe  bool   `json:"remember_me"`
}

func (r *GoogleCallbackReq) ToDTO() *dto.GoogleOAuthRequest {
	if r == nil {
		return nil
	}

	return &dto.GoogleOAuthRequest{
		Code:        r.Code,
		RedirectUri: r.RedirectUri,
		RememberMe:  r.RememberMe,
	}
}
