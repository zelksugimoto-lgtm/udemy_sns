package middleware

import (
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// getActionDescription パスとメソッドから操作内容を日本語で返す
func getActionDescription(method, path string) string {
	// 正規表現パターン
	patterns := []struct {
		regex       *regexp.Regexp
		method      string
		description string
	}{
		// 認証系
		{regexp.MustCompile(`^/api/v1/auth/register$`), "POST", "ユーザー登録"},
		{regexp.MustCompile(`^/api/v1/auth/login$`), "POST", "ログイン"},
		{regexp.MustCompile(`^/api/v1/auth/logout$`), "POST", "ログアウト"},
		{regexp.MustCompile(`^/api/v1/auth/refresh$`), "POST", "トークン更新"},
		{regexp.MustCompile(`^/api/v1/auth/revoke-all$`), "POST", "全トークン無効化"},
		{regexp.MustCompile(`^/api/v1/auth/me$`), "GET", "ログインユーザー情報取得"},

		// 投稿系
		{regexp.MustCompile(`^/api/v1/posts$`), "POST", "投稿作成"},
		{regexp.MustCompile(`^/api/v1/posts/[^/]+$`), "GET", "投稿詳細取得"},
		{regexp.MustCompile(`^/api/v1/posts/[^/]+$`), "PATCH", "投稿編集"},
		{regexp.MustCompile(`^/api/v1/posts/[^/]+$`), "DELETE", "投稿削除"},

		// いいね系
		{regexp.MustCompile(`^/api/v1/posts/[^/]+/like$`), "POST", "投稿にいいね"},
		{regexp.MustCompile(`^/api/v1/posts/[^/]+/like$`), "DELETE", "投稿のいいね解除"},
		{regexp.MustCompile(`^/api/v1/comments/[^/]+/like$`), "POST", "コメントにいいね"},
		{regexp.MustCompile(`^/api/v1/comments/[^/]+/like$`), "DELETE", "コメントのいいね解除"},

		// コメント系
		{regexp.MustCompile(`^/api/v1/posts/[^/]+/comments$`), "POST", "コメント投稿"},
		{regexp.MustCompile(`^/api/v1/posts/[^/]+/comments$`), "GET", "コメント一覧取得"},
		{regexp.MustCompile(`^/api/v1/comments/[^/]+$`), "DELETE", "コメント削除"},

		// ブックマーク系
		{regexp.MustCompile(`^/api/v1/posts/[^/]+/bookmark$`), "POST", "ブックマーク追加"},
		{regexp.MustCompile(`^/api/v1/posts/[^/]+/bookmark$`), "DELETE", "ブックマーク削除"},
		{regexp.MustCompile(`^/api/v1/bookmarks$`), "GET", "ブックマーク一覧取得"},

		// フォロー系
		{regexp.MustCompile(`^/api/v1/users/[^/]+/follow$`), "POST", "ユーザーをフォロー"},
		{regexp.MustCompile(`^/api/v1/users/[^/]+/follow$`), "DELETE", "ユーザーのフォロー解除"},
		{regexp.MustCompile(`^/api/v1/users/[^/]+/followers$`), "GET", "フォロワー一覧取得"},
		{regexp.MustCompile(`^/api/v1/users/[^/]+/following$`), "GET", "フォロー中一覧取得"},

		// ユーザー系
		{regexp.MustCompile(`^/api/v1/users$`), "GET", "ユーザー検索"},
		{regexp.MustCompile(`^/api/v1/users/[^/]+$`), "GET", "プロフィール取得"},
		{regexp.MustCompile(`^/api/v1/users/me$`), "PATCH", "プロフィール編集"},
		{regexp.MustCompile(`^/api/v1/users/[^/]+/posts$`), "GET", "ユーザーの投稿一覧取得"},
		{regexp.MustCompile(`^/api/v1/users/[^/]+/likes$`), "GET", "ユーザーのいいね一覧取得"},

		// タイムライン
		{regexp.MustCompile(`^/api/v1/timeline$`), "GET", "タイムライン取得"},

		// 通知系
		{regexp.MustCompile(`^/api/v1/notifications$`), "GET", "通知一覧取得"},
		{regexp.MustCompile(`^/api/v1/notifications/unread-count$`), "GET", "未読通知数取得"},
		{regexp.MustCompile(`^/api/v1/notifications/[^/]+/read$`), "PATCH", "通知を既読にする"},
		{regexp.MustCompile(`^/api/v1/notifications/read-all$`), "POST", "全通知を既読にする"},

		// 通報系
		{regexp.MustCompile(`^/api/v1/reports$`), "POST", "コンテンツを通報"},

		// ヘルスチェック
		{regexp.MustCompile(`^/health$`), "GET", "ヘルスチェック"},
		{regexp.MustCompile(`^/swagger/.*$`), "GET", "API仕様書閲覧"},
	}

	// パスとメソッドでマッチング
	for _, p := range patterns {
		if p.regex.MatchString(path) && strings.EqualFold(p.method, method) {
			return p.description
		}
	}

	// マッチしない場合は空文字列
	return ""
}

// AccessLog 全APIリクエストのアクセスログを出力
func AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// リクエスト開始時刻
			start := time.Now()

			// リクエストIDを取得
			requestID := GetRequestID(c)

			// ユーザーIDを取得（認証済みの場合）
			userID, err := GetUserID(c) // 既存のGetUserID関数を使用

			// 次のハンドラーを実行
			handlerErr := next(c)

			// レスポンスタイム計算
			responseTime := time.Since(start)

			// 操作内容を取得
			method := c.Request().Method
			path := c.Request().URL.Path
			actionDescription := getActionDescription(method, path)

			// ログ出力
			logEvent := log.Info().
				Str("request_id", requestID).
				Str("method", method).
				Str("path", path).
				Int("status", c.Response().Status).
				Dur("response_time", responseTime).
				Str("remote_ip", c.RealIP())

			// 操作内容がある場合のみ追加
			if actionDescription != "" {
				logEvent = logEvent.Str("action", actionDescription)
			}

			// ユーザーIDがある場合のみ追加（エラーがない＝認証済み）
			if err == nil {
				logEvent = logEvent.Str("user_id", userID.String())
			}

			// クエリパラメータがある場合
			if c.Request().URL.RawQuery != "" {
				logEvent = logEvent.Str("query", c.Request().URL.RawQuery)
			}

			logEvent.Msg("access")

			return handlerErr
		}
	}
}
