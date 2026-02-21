package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type AuthHandler struct {
	// サービスは後で実装
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register godoc
// @Summary      ユーザー登録
// @Description  新しいユーザーを登録し、JWTトークンを返す
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
	// TODO: 実装
	return c.JSON(201, response.AuthResponse{})
}

// Login godoc
// @Summary      ログイン
// @Description  メールアドレスとパスワードでログインし、JWTトークンを返す
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
	// TODO: 実装
	return c.JSON(200, response.AuthResponse{})
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
	// TODO: 実装
	return c.JSON(200, response.UserResponse{})
}
