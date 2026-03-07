package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/admin/session"
)

// AdminAuthMiddleware 管理者認証ミドルウェア
func AdminAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// セッションを取得
			sess, err := session.GetStore().Get(c.Request(), session.SessionName)
			if err != nil {
				log.Error().Err(err).Msg("セッション取得エラー")
				return c.Redirect(http.StatusFound, "/admin/login")
			}

			// セッションから管理者情報を取得
			adminID, _ := sess.Values[session.KeyAdminID].(string)
			isRoot, _ := sess.Values[session.KeyIsRoot].(bool)

			// 認証チェック
			if adminID == "" && !isRoot {
				return c.Redirect(http.StatusFound, "/admin/login")
			}

			// コンテキストに管理者情報を設定
			c.Set("admin_id", adminID)
			c.Set("is_root", isRoot)
			if username, ok := sess.Values[session.KeyUsername].(string); ok {
				c.Set("admin_username", username)
			}
			if role, ok := sess.Values[session.KeyRole].(string); ok {
				c.Set("admin_role", role)
			}

			return next(c)
		}
	}
}

// RootAdminMiddleware ルート管理者権限チェックミドルウェア
func RootAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isRoot, _ := c.Get("is_root").(bool)
			role, _ := c.Get("admin_role").(string)

			if !isRoot && role != "root" {
				return echo.NewHTTPError(http.StatusForbidden, "この操作にはルート管理者権限が必要です")
			}

			return next(c)
		}
	}
}
