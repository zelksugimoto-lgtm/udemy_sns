package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type CommentHandler struct {
	// サービスは後で実装
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

// CreateComment godoc
// @Summary      コメント作成
// @Description  投稿にコメントを追加（ネスト対応）
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                        true  "投稿ID"
// @Param        request  body      request.CreateCommentRequest  true  "コメント内容"
// @Success      201      {object}  response.CommentResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      401      {object}  errors.ErrorResponse
// @Failure      404      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(c echo.Context) error {
	// TODO: 実装
	return c.JSON(201, response.CommentResponse{})
}

// GetComments godoc
// @Summary      コメント一覧取得
// @Description  投稿のコメント一覧を取得（ネスト構造）
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id      path      string  true   "投稿ID"
// @Param        limit   query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int     false  "オフセット（デフォルト: 0）"
// @Success      200     {object}  response.CommentListResponse
// @Failure      404     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /posts/{id}/comments [get]
func (h *CommentHandler) GetComments(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.CommentListResponse{})
}

// DeleteComment godoc
// @Summary      コメント削除
// @Description  コメントを削除（論理削除、コメント投稿者のみ）
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "コメントID"
// @Success      204
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(204)
}
