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
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
// @Summary      投稿作成
// @Description  新しい投稿を作成
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreatePostRequest  true  "投稿内容"
// @Success      201      {object}  response.PostResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      401      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /posts [post]
func (h *PostHandler) CreatePost(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// リクエストボディをパース
	var req request.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("投稿内容は必須です"))
	}
	if req.Visibility == "" {
		req.Visibility = "public" // デフォルトは公開
	}
	if req.Visibility != "public" && req.Visibility != "followers" && req.Visibility != "private" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な公開範囲です"))
	}

	// サービス層を呼び出し
	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError("投稿の作成に失敗しました"))
	}

	return c.JSON(http.StatusCreated, post)
}

// GetPost godoc
// @Summary      投稿取得
// @Description  IDから投稿を取得
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "投稿ID"
// @Success      200 {object}  response.PostResponse
// @Failure      404 {object}  errors.ErrorResponse
// @Failure      500 {object}  errors.ErrorResponse
// @Router       /posts/{id} [get]
func (h *PostHandler) GetPost(c echo.Context) error {
	// 投稿IDを取得
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な投稿IDです"))
	}

	// オプショナルな認証情報を取得
	var currentUserID *uuid.UUID
	if userID, ok := c.Get("user_id").(uuid.UUID); ok {
		currentUserID = &userID
	}

	// サービス層を呼び出し
	post, err := h.postService.GetPost(postID, currentUserID)
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("投稿の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, post)
}

// UpdatePost godoc
// @Summary      投稿更新
// @Description  投稿内容を更新（投稿者のみ）
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                     true  "投稿ID"
// @Param        request  body      request.UpdatePostRequest  true  "更新内容"
// @Success      200      {object}  response.PostResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      401      {object}  errors.ErrorResponse
// @Failure      403      {object}  errors.ErrorResponse
// @Failure      404      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /posts/{id} [patch]
func (h *PostHandler) UpdatePost(c echo.Context) error {
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
	var req request.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if req.Visibility != nil {
		if *req.Visibility != "public" && *req.Visibility != "followers" && *req.Visibility != "private" {
			return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な公開範囲です"))
		}
	}

	// サービス層を呼び出し
	post, err := h.postService.UpdatePost(userID, postID, &req)
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		if err.Error() == "forbidden" || err.Error() == "この投稿を更新する権限がありません" {
			return c.JSON(http.StatusForbidden, errors.Forbidden("この投稿を編集する権限がありません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("投稿の更新に失敗しました"))
	}

	return c.JSON(http.StatusOK, post)
}

// DeletePost godoc
// @Summary      投稿削除
// @Description  投稿を削除（論理削除、投稿者のみ）
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "投稿ID"
// @Success      204
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /posts/{id} [delete]
func (h *PostHandler) DeletePost(c echo.Context) error {
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
	err = h.postService.DeletePost(userID, postID)
	if err != nil {
		if err.Error() == "投稿が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
		}
		if err.Error() == "この投稿を削除する権限がありません" {
			return c.JSON(http.StatusForbidden, errors.Forbidden("この投稿を削除する権限がありません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("投稿の削除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
}

// GetTimeline godoc
// @Summary      タイムライン取得
// @Description  フォロー中のユーザー + 自分の投稿を時系列で取得
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     int  false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int  false  "オフセット（デフォルト: 0）"
// @Success      200     {object}  response.PostListResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /timeline [get]
func (h *PostHandler) GetTimeline(c echo.Context) error {
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
	result, err := h.postService.GetTimeline(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError("タイムラインの取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
}

// GetUserPosts godoc
// @Summary      ユーザーの投稿一覧取得
// @Description  指定したユーザー名の投稿一覧を取得
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        username  path      string  true   "ユーザー名"
// @Param        limit     query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset    query     int     false  "オフセット（デフォルト: 0）"
// @Success      200       {object}  response.PostListResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username}/posts [get]
func (h *PostHandler) GetUserPosts(c echo.Context) error {
	// ユーザー名を取得
	username := c.Param("username")

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

	// オプショナルな認証情報を取得
	var currentUserID *uuid.UUID
	if userID, ok := c.Get("user_id").(uuid.UUID); ok {
		currentUserID = &userID
	}

	// サービス層を呼び出し
	result, err := h.postService.GetUserPosts(username, limit, offset, currentUserID)
	if err != nil {
		if err.Error() == "ユーザーが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("ユーザーが見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("投稿一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
}
