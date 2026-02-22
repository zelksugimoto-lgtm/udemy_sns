package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
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
	var req request.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if req.Email == "" || req.Username == "" || req.DisplayName == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("必須項目が不足しています"))
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == "このメールアドレスは既に使用されています" || err.Error() == "このユーザー名は既に使用されています" {
			return c.JSON(http.StatusConflict, errors.Conflict(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError(err.Error()))
	}

	return c.JSON(http.StatusCreated, result)
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
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("必須項目が不足しています"))
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
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
