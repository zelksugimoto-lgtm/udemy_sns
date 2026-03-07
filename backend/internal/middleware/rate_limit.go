package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// bucket レートリミット用のバケット
type bucket struct {
	count     int
	resetTime time.Time
}

// RateLimiter レートリミッター
type RateLimiter struct {
	limit   int                  // 制限回数
	window  time.Duration        // 時間ウィンドウ
	buckets map[string]*bucket   // キーごとのバケット
	mu      sync.RWMutex         // スレッドセーフ用のmutex
}

// NewRateLimiter レートリミッターのコンストラクタ
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		limit:   limit,
		window:  window,
		buckets: make(map[string]*bucket),
	}

	// 定期的に古いバケットをクリーンアップ（メモリリーク防止）
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return limiter
}

// Allow レート制限をチェックし、許可するかどうかを返す
// 戻り値: (許可するか, 制限値, 残り回数, リセット時刻)
func (rl *RateLimiter) Allow(key string) (bool, int, int, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	// バケットが存在しないか、リセット時刻を過ぎている場合は新しいバケットを作成
	if !exists || now.After(b.resetTime) {
		b = &bucket{
			count:     0,
			resetTime: now.Add(rl.window),
		}
		rl.buckets[key] = b
	}

	// 制限回数以内かチェック
	if b.count < rl.limit {
		b.count++
		remaining := rl.limit - b.count
		return true, rl.limit, remaining, b.resetTime
	}

	// 制限超過
	return false, rl.limit, 0, b.resetTime
}

// cleanup 古いバケットをクリーンアップ
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, b := range rl.buckets {
		if now.After(b.resetTime) {
			delete(rl.buckets, key)
		}
	}
}

// getClientIP クライアントのIPアドレスを取得
// Cloud Run環境を考慮し、X-Forwarded-Forヘッダーから取得
func getClientIP(c echo.Context) string {
	// X-Forwarded-For ヘッダーをチェック（Cloud Run, リバースプロキシ対応）
	xff := c.Request().Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For: client, proxy1, proxy2
		// 最初のIPアドレス（クライアントの実IP）を取得
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// X-Real-IP ヘッダーをチェック
	xri := c.Request().Header.Get("X-Real-IP")
	if xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
		}
	}

	// フォールバック: RemoteAddr から取得
	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return c.Request().RemoteAddr
	}
	return ip
}

// RateLimitMiddleware 汎用レートリミットミドルウェア
func RateLimitMiddleware(limiter *RateLimiter, keyFunc func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// キーを取得
			key := keyFunc(c)

			// レート制限チェック
			allowed, limit, remaining, resetTime := limiter.Allow(key)

			// レスポンスヘッダーを設定
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			if !allowed {
				// レート制限超過
				retryAfter := int(time.Until(resetTime).Seconds())
				if retryAfter < 0 {
					retryAfter = 0
				}
				c.Response().Header().Set("Retry-After", strconv.Itoa(retryAfter))

				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "RATE_LIMIT_EXCEEDED",
						"message": "レート制限を超過しました。しばらく待ってから再試行してください。",
					},
				})
			}

			return next(c)
		}
	}
}

// AuthRateLimitMiddleware 認証系API用のレートリミットミドルウェア
// IPアドレス単位で5回/分の制限
func AuthRateLimitMiddleware() echo.MiddlewareFunc {
	limiter := NewRateLimiter(5, 1*time.Minute)
	return RateLimitMiddleware(limiter, func(c echo.Context) string {
		return getClientIP(c)
	})
}

// GeneralRateLimitMiddleware 一般API用のレートリミットミドルウェア
// 認証済みユーザーはユーザーID単位、未認証はIPアドレス単位で60回/分の制限
func GeneralRateLimitMiddleware() echo.MiddlewareFunc {
	limiter := NewRateLimiter(60, 1*time.Minute)
	return RateLimitMiddleware(limiter, func(c echo.Context) string {
		// 認証済みユーザーの場合はユーザーIDを使用
		userID := c.Get("userID")
		if userID != nil {
			return "user:" + userID.(string)
		}
		// 未認証の場合はIPアドレスを使用
		return "ip:" + getClientIP(c)
	})
}
