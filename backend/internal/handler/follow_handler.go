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

type FollowHandler struct {
	followService service.FollowService
}

func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// ユーザー名を取得
	username := c.Param("username")

	// サービス層を呼び出し
	err = h.followService.Follow(userID, username)
	if err != nil {
		if err.Error() == "ユーザーが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("ユーザーが見つかりません"))
		}
		if err.Error() == "自分自身をフォローすることはできません" {
			return c.JSON(http.StatusBadRequest, errors.BadRequest("自分自身をフォローすることはできません"))
		}
		if err.Error() == "すでにフォローしています" {
			return c.JSON(http.StatusConflict, errors.Conflict("すでにフォローしています"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("フォローに失敗しました"))
	}

	return c.NoContent(http.StatusCreated)
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
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// ユーザー名を取得
	username := c.Param("username")

	// サービス層を呼び出し
	err = h.followService.Unfollow(userID, username)
	if err != nil {
		if err.Error() == "ユーザーが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("ユーザーが見つかりません"))
		}
		if err.Error() == "フォローしていません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("フォローしていません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("フォロー解除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
}

// GetFollowers godoc
// @Summary      フォロワー一覧取得
// @Description  ユーザーのフォロワー一覧を取得
// @Tags         follows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username  path      string  true   "ユーザー名"
// @Param        limit     query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset    query     int     false  "オフセット（デフォルト: 0）"
// @Success      200       {object}  response.UserListResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username}/followers [get]
func (h *FollowHandler) GetFollowers(c echo.Context) error {
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

	// 認証ユーザーIDを取得（任意）
	var currentUserID *uuid.UUID
	userID, err := middleware.GetUserID(c)
	if err == nil {
		currentUserID = &userID
	}

	// サービス層を呼び出し
	result, err := h.followService.GetFollowers(username, limit, offset, currentUserID)
	if err != nil {
		if err.Error() == "ユーザーが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("ユーザーが見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("フォロワー一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
}

// GetFollowing godoc
// @Summary      フォロー中一覧取得
// @Description  ユーザーがフォロー中のユーザー一覧を取得
// @Tags         follows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username  path      string  true   "ユーザー名"
// @Param        limit     query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset    query     int     false  "オフセット（デフォルト: 0）"
// @Success      200       {object}  response.UserListResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username}/following [get]
func (h *FollowHandler) GetFollowing(c echo.Context) error {
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

	// 認証ユーザーIDを取得（任意）
	var currentUserID *uuid.UUID
	userID, err := middleware.GetUserID(c)
	if err == nil {
		currentUserID = &userID
	}

	// サービス層を呼び出し
	result, err := h.followService.GetFollowing(username, limit, offset, currentUserID)
	if err != nil {
		if err.Error() == "ユーザーが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("ユーザーが見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("フォロー中一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
}
