package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/pkg/errors"
)

type LikeHandler struct {
	// サービスは後で実装
}

func NewLikeHandler() *LikeHandler {
	return &LikeHandler{}
}

// LikePost godoc
// @Summary      投稿にいいね
// @Description  投稿にいいねを追加
// @Tags         likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "投稿ID"
// @Success      201
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      409  {object}  errors.ErrorResponse "既にいいね済み"
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /posts/{id}/like [post]
func (h *LikeHandler) LikePost(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(201)
}

// UnlikePost godoc
// @Summary      投稿のいいね解除
// @Description  投稿のいいねを削除
// @Tags         likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "投稿ID"
// @Success      204
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /posts/{id}/like [delete]
func (h *LikeHandler) UnlikePost(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(204)
}

// LikeComment godoc
// @Summary      コメントにいいね
// @Description  コメントにいいねを追加
// @Tags         likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "コメントID"
// @Success      201
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      409  {object}  errors.ErrorResponse "既にいいね済み"
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /comments/{id}/like [post]
func (h *LikeHandler) LikeComment(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(201)
}

// UnlikeComment godoc
// @Summary      コメントのいいね解除
// @Description  コメントのいいねを削除
// @Tags         likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "コメントID"
// @Success      204
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /comments/{id}/like [delete]
func (h *LikeHandler) UnlikeComment(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(204)
}
