package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
)

type LikeHandler struct {
	likeService service.LikeService
}

func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
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
	err = h.likeService.LikePost(userID, postID)
	if err != nil {
		if err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		if err.Error() == "すでにいいねしています" {
			return c.JSON(http.StatusConflict, errors.Conflict("すでにいいねしています"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("いいねに失敗しました"))
	}

	return c.NoContent(http.StatusCreated)
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
	err = h.likeService.UnlikePost(userID, postID)
	if err != nil {
		if err.Error() == "いいねしていません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("いいねしていません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("いいね解除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// コメントIDを取得
	commentIDStr := c.Param("id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なコメントIDです"))
	}

	// サービス層を呼び出し
	err = h.likeService.LikeComment(userID, commentID)
	if err != nil {
		if err.Error() == "コメントが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("コメントが見つかりません"))
		}
		if err.Error() == "すでにいいねしています" {
			return c.JSON(http.StatusConflict, errors.Conflict("すでにいいねしています"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("いいねに失敗しました"))
	}

	return c.NoContent(http.StatusCreated)
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// コメントIDを取得
	commentIDStr := c.Param("id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なコメントIDです"))
	}

	// サービス層を呼び出し
	err = h.likeService.UnlikeComment(userID, commentID)
	if err != nil {
		if err.Error() == "いいねしていません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("いいねしていません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("いいね解除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
}
