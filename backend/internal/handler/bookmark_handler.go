package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type BookmarkHandler struct {
	// サービスは後で実装
}

func NewBookmarkHandler() *BookmarkHandler {
	return &BookmarkHandler{}
}

// AddBookmark godoc
// @Summary      ブックマーク追加
// @Description  投稿をブックマークに追加
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "投稿ID"
// @Success      201
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      409  {object}  errors.ErrorResponse "既にブックマーク済み"
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /posts/{id}/bookmark [post]
func (h *BookmarkHandler) AddBookmark(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(201)
}

// RemoveBookmark godoc
// @Summary      ブックマーク削除
// @Description  ブックマークを削除
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "投稿ID"
// @Success      204
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /posts/{id}/bookmark [delete]
func (h *BookmarkHandler) RemoveBookmark(c echo.Context) error {
	// TODO: 実装
	return c.NoContent(204)
}

// GetBookmarks godoc
// @Summary      ブックマーク一覧取得
// @Description  ブックマークした投稿の一覧を取得
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     int  false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int  false  "オフセット（デフォルト: 0）"
// @Success      200     {object}  response.BookmarkListResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /bookmarks [get]
func (h *BookmarkHandler) GetBookmarks(c echo.Context) error {
	// TODO: 実装
	return c.JSON(200, response.BookmarkListResponse{})
}
