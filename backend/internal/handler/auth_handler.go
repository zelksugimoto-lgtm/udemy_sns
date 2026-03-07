package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
	"github.com/yourusername/sns-app/pkg/validator"
)

const (
	// AccessTokenCookie アクセストークンのCookie名
	AccessTokenCookie = "access_token"
	// RefreshTokenCookie リフレッシュトークンのCookie名
	RefreshTokenCookie = "refresh_token"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      ユーザー登録
// @Description  新しいユーザーを登録し、HttpOnly CookieでJWTトークンを返す
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.RegisterRequest  true  "登録情報"
// @Success      201      {object}  response.AuthResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      409      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req request.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if err := validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
	}

	result, tokenPair, err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == "このメールアドレスは既に使用されています" || err.Error() == "このユーザー名は既に使用されています" {
			return c.JSON(http.StatusConflict, errors.Conflict(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError(err.Error()))
	}

	// Cookieを設定
	setAuthCookies(c, tokenPair.AccessToken, tokenPair.RefreshToken)

	return c.JSON(http.StatusCreated, result)
}

// Login godoc
// @Summary      ログイン
// @Description  メールアドレスとパスワードでログインし、HttpOnly CookieでJWTトークンを返す
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.LoginRequest  true  "ログイン情報"
// @Success      200      {object}  response.AuthResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      401      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if err := validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
	}

	result, tokenPair, err := h.authService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized(err.Error()))
	}

	// Cookieを設定
	setAuthCookies(c, tokenPair.AccessToken, tokenPair.RefreshToken)

	return c.JSON(http.StatusOK, result)
}

// RefreshToken godoc
// @Summary      トークンリフレッシュ
// @Description  リフレッシュトークンを使用して新しいアクセストークンを発行
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// CookieからリフレッシュトークンをRefreshTokenCookie取得
	cookie, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("リフレッシュトークンが見つかりません"))
	}

	// トークンリフレッシュ
	tokenPair, err := h.authService.RefreshToken(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized(err.Error()))
	}

	// 新しいCookieを設定
	setAuthCookies(c, tokenPair.AccessToken, tokenPair.RefreshToken)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "トークンをリフレッシュしました",
	})
}

// Logout godoc
// @Summary      ログアウト
// @Description  リフレッシュトークンを無効化し、Cookieを削除
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	// CookieからリフレッシュトークンをRefreshTokenCookie取得
	cookie, err := c.Cookie(RefreshTokenCookie)
	if err == nil && cookie.Value != "" {
		// リフレッシュトークンを無効化
		_ = h.authService.Logout(cookie.Value)
	}

	// Cookieを削除
	clearAuthCookies(c)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ログアウトしました",
	})
}

// RevokeAllTokens godoc
// @Summary      全デバイスログアウト
// @Description  ユーザーの全リフレッシュトークンを無効化
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /auth/revoke-all [post]
func (h *AuthHandler) RevokeAllTokens(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized(err.Error()))
	}

	if err := h.authService.RevokeAllTokens(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError(err.Error()))
	}

	// 現在のデバイスのCookieも削除
	clearAuthCookies(c)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "全てのデバイスからログアウトしました",
	})
}

// GetMe godoc
// @Summary      現在のユーザー情報取得
// @Description  JWTトークンから現在ログイン中のユーザー情報を取得
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.UserResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized(err.Error()))
	}

	result, err := h.authService.GetMe(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

// setAuthCookies アクセストークンとリフレッシュトークンをCookieに設定
func setAuthCookies(c echo.Context, accessToken, refreshToken string) {
	// 環境判定（本番環境かどうか）
	isProduction := os.Getenv("ENV") == "production"

	// SameSite設定: 本番環境（HTTPS）ではNone、開発環境（HTTP）ではLax
	// SameSite=NoneはSecure=trueが必須のため、開発環境ではLaxを使用
	sameSite := http.SameSiteLaxMode
	if isProduction {
		sameSite = http.SameSiteNoneMode
	}

	// アクセストークンCookie（1時間）
	accessCookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(service.AccessTokenDuration),
		HttpOnly: true,
		Secure:   isProduction, // 本番環境のみSecure
		SameSite: sameSite,
	}
	c.SetCookie(accessCookie)

	// リフレッシュトークンCookie（7日間）
	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(service.RefreshTokenDuration),
		HttpOnly: true,
		Secure:   isProduction, // 本番環境のみSecure
		SameSite: sameSite,
	}
	c.SetCookie(refreshCookie)
}

// clearAuthCookies 認証Cookieを削除
func clearAuthCookies(c echo.Context) {
	// アクセストークンCookieを削除
	accessCookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	c.SetCookie(accessCookie)

	// リフレッシュトークンCookieを削除
	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	c.SetCookie(refreshCookie)
}
