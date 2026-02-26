package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile godoc
// @Summary      ユーザープロフィール取得
// @Description  ユーザー名からプロフィール情報を取得
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username  path      string  true  "ユーザー名"
// @Success      200       {object}  response.UserProfileResponse
// @Failure      404       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /users/{username} [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	username := c.Param("username")

	result, err := h.userService.GetProfile(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.NotFound(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateProfile godoc
// @Summary      プロフィール更新
// @Description  現在ログイン中のユーザーのプロフィールを更新
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateProfileRequest  true  "更新情報"
// @Success      200      {object}  response.UserResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      401      {object}  errors.ErrorResponse
// @Failure      500      {object}  errors.ErrorResponse
// @Router       /users/me [patch]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized(err.Error()))
	}

	var req request.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	result, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

// SearchUsers godoc
// @Summary      ユーザー検索
// @Description  キーワードでユーザーを検索
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        q       query     string  true   "検索キーワード"
// @Param        limit   query     int     false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int     false  "オフセット（デフォルト: 0）"
// @Success      200     {object}  response.UserListResponse
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) SearchUsers(c echo.Context) error {
	query := c.QueryParam("q")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	result, err := h.userService.SearchUsers(query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}
