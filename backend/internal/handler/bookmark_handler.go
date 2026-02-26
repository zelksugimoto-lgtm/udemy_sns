package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
)

type BookmarkHandler struct {
	bookmarkService service.BookmarkService
}

func NewBookmarkHandler(bookmarkService service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkService: bookmarkService,
	}
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// 投稿IDを取得
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な投稿IDです"))
	}

	// サービス層を呼び出し
	err = h.bookmarkService.AddBookmark(userID, postID)
	if err != nil {
		if err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		if err.Error() == "すでにブックマークしています" {
			return c.JSON(http.StatusConflict, errors.Conflict("すでにブックマークしています"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("ブックマークに失敗しました"))
	}

	return c.NoContent(http.StatusCreated)
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// 投稿IDを取得
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な投稿IDです"))
	}

	// サービス層を呼び出し
	err = h.bookmarkService.RemoveBookmark(userID, postID)
	if err != nil {
		if err.Error() == "ブックマークしていません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("ブックマークしていません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("ブックマーク削除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// クエリパラメータを取得
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	// デフォルト値設定
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// サービス層を呼び出し
	result, err := h.bookmarkService.GetBookmarks(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError("ブックマーク一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
}
