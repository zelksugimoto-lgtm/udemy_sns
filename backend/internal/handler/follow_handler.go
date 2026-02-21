package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type FollowHandler struct {
	// サービスは後で実装
}

func NewFollowHandler() *FollowHandler {
	return &FollowHandler{}
}

// Follow godoc
// @Summary      フォロー
// @Description  ユーザーをフォロー
// @Tags         follows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username  path  string  true  "ユーザー名"
// @Success      201
// @Failure      400  {object}  errors.ErrorResponse "自分自身をフォローできない"
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      409  {object}  errors.ErrorResponse "既にフォロー済み"
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /users/{username}/follow [post]
func (h *FollowHandler) Follow(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(201)
}

// Unfollow godoc
// @Summary      フォロー解除
// @Description  ユーザーのフォローを解除
// @Tags         follows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username  path  string  true  "ユーザー名"
// @Success      204
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /users/{username}/follow [delete]
func (h *FollowHandler) Unfollow(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(204)
}

// GetFollowers godoc
// @Summary      フォロワー一覧取得
// @Description  ユーザーのフォロワー一覧を取得
// @Tags         follows
// @Accept       json
// @Produce      json
// @Param        username  path      string  true   "ユーザー名"
// @Param        limit     query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset    query     int     false  "オフセット（デフォルト: 0）"
// @Success      200       {object}  response.FollowListResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username}/followers [get]
func (h *FollowHandler) GetFollowers(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.FollowListResponse{})
}

// GetFollowing godoc
// @Summary      フォロー中一覧取得
// @Description  ユーザーがフォロー中のユーザー一覧を取得
// @Tags         follows
// @Accept       json
// @Produce      json
// @Param        username  path      string  true   "ユーザー名"
// @Param        limit     query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset    query     int     false  "オフセット（デフォルト: 0）"
// @Success      200       {object}  response.FollowListResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username}/following [get]
func (h *FollowHandler) GetFollowing(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.FollowListResponse{})
}
