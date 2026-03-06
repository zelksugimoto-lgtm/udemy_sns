package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
	"github.com/yourusername/sns-app/pkg/validator"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
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

	// リクエストボディをパース
	var req request.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if err := validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
	}

	// サービス層を呼び出し
	comment, err := h.commentService.CreateComment(userID, postID, &req)
	if err != nil {
		if err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		if err.Error() == "親コメントが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("親コメントが見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("コメントの作成に失敗しました"))
	}

	return c.JSON(http.StatusCreated, comment)
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
	// 投稿IDを取得
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な投稿IDです"))
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
	result, err := h.commentService.GetComments(postID, limit, offset)
	if err != nil {
		if err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("コメント一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
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
	err = h.commentService.DeleteComment(userID, commentID)
	if err != nil {
		if err.Error() == "コメントが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("コメントが見つかりません"))
		}
		if err.Error() == "このコメントを削除する権限がありません" {
			return c.JSON(http.StatusForbidden, errors.Forbidden("このコメントを削除する権限がありません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("コメントの削除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
}
