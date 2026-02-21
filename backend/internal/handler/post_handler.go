package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type PostHandler struct {
	// サービスは後で実装
}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
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
	// TODO: 実装
	return c.JSON(201, response.PostResponse{})
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
	// TODO: 実装
	return c.JSON(200, response.PostResponse{})
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
	// TODO: 実装
	return c.JSON(200, response.PostResponse{})
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
	// TODO: 実装
	return c.NoContent(204)
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
	// TODO: 実装
	return c.JSON(200, response.PostListResponse{})
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
	// TODO: 実装
	return c.JSON(200, response.PostListResponse{})
}
