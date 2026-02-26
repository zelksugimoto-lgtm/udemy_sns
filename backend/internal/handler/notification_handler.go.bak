package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/pkg/errors"
)

type NotificationHandler struct {
	// サービスは後で実装
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
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
	// TODO: 実装
	return c.JSON(200, response.NotificationListResponse{})
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
	// TODO: 実装
	return c.NoContent(200)
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
	// TODO: 実装
	return c.NoContent(200)
}
