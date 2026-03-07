package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	// RequestIDHeader リクエストIDのヘッダー名
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey コンテキストキー
	RequestIDKey = "request_id"
)

// RequestID リクエストごとにユニークなIDを生成し、レスポンスヘッダーとコンテキストに設定
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// リクエストヘッダーからリクエストIDを取得（存在する場合）
			requestID := c.Request().Header.Get(RequestIDHeader)

			// リクエストIDがない場合は新規生成
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// レスポンスヘッダーに設定
			c.Response().Header().Set(RequestIDHeader, requestID)

			// コンテキストに設定（ログやエラーハンドリングで使用）
			c.Set(RequestIDKey, requestID)

			return next(c)
		}
	}
}

// GetRequestID コンテキストからリクエストIDを取得
func GetRequestID(c echo.Context) string {
	requestID, ok := c.Get(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}
