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

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotifications godoc
// @Summary      通知一覧取得
// @Description  ログイン中のユーザーの通知一覧を取得
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     int  false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int  false  "オフセット（デフォルト: 0）"
// @Success      200     {object}  response.NotificationListResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /notifications [get]
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
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
	result, err := h.notificationService.GetNotifications(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError("通知一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, result)
}

// MarkAsRead godoc
// @Summary      通知を既読にする
// @Description  指定した通知を既読にする
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "通知ID"
// @Success      200
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /notifications/{id}/read [patch]
func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// 通知IDを取得
	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な通知IDです"))
	}

	// サービス層を呼び出し
	err = h.notificationService.MarkAsRead(userID, notificationID)
	if err != nil {
		if err.Error() == "通知が見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound("通知が見つかりません"))
		}
		if err.Error() == "この通知を既読にする権限がありません" {
			return c.JSON(http.StatusForbidden, errors.Forbidden("この通知を既読にする権限がありません"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("通知の既読化に失敗しました"))
	}

	return c.NoContent(http.StatusOK)
}

// MarkAllAsRead godoc
// @Summary      すべての通知を既読にする
// @Description  ログイン中のユーザーのすべての通知を既読にする
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /notifications/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// サービス層を呼び出し
	err = h.notificationService.MarkAllAsRead(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError("すべての通知の既読化に失敗しました"))
	}

	return c.NoContent(http.StatusOK)
}

// GetUnreadCount godoc
// @Summary      未読通知数取得
// @Description  ログイン中のユーザーの未読通知数を取得
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]int "count: 未読通知数"
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// サービス層を呼び出し
	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.InternalError("未読通知数の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, map[string]int64{"count": count})
}
