package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type UserHandler struct {
	// サービスは後で実装
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetProfile godoc
// @Summary      ユーザープロフィール取得
// @Description  ユーザー名からプロフィール情報を取得
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username  path      string  true  "ユーザー名"
// @Success      200       {object}  response.UserProfileResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username} [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.UserProfileResponse{})
}

// UpdateProfile godoc
// @Summary      プロフィール更新
// @Description  現在ログイン中のユーザーのプロフィールを更新
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateProfileRequest  true  "更新情報"
// @Success      200      {object}  response.UserResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      401      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /users/me [patch]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.UserResponse{})
}

// SearchUsers godoc
// @Summary      ユーザー検索
// @Description  キーワードでユーザーを検索
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        q       query     string  true   "検索キーワード"
// @Param        limit   query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int     false  "オフセット（デフォルト: 0）"
// @Success      200     {object}  response.UserListResponse
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) SearchUsers(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.UserListResponse{})
}
