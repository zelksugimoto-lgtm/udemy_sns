package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		window         time.Duration
		requests       int
		expectedAllows int
	}{
		{
			name:           "正常系_制限内のリクエスト",
			limit:          5,
			window:         1 * time.Minute,
			requests:       3,
			expectedAllows: 3,
		},
		{
			name:           "正常系_制限ちょうどのリクエスト",
			limit:          5,
			window:         1 * time.Minute,
			requests:       5,
			expectedAllows: 5,
		},
		{
			name:           "異常系_制限超過",
			limit:          5,
			window:         1 * time.Minute,
			requests:       10,
			expectedAllows: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewRateLimiter(tt.limit, tt.window)
			key := "test-key"
			allows := 0

			for i := 0; i < tt.requests; i++ {
				allowed, _, _, _ := limiter.Allow(key)
				if allowed {
					allows++
				}
			}

			assert.Equal(t, tt.expectedAllows, allows, "許可された回数が期待値と一致すること")
		})
	}
}

func TestRateLimiter_ResetAfterWindow(t *testing.T) {
	t.Run("正常系_ウィンドウ経過後にリセット", func(t *testing.T) {
		limiter := NewRateLimiter(2, 100*time.Millisecond)
		key := "test-key"

		// 2回リクエスト（制限に達する）
		allowed1, _, _, _ := limiter.Allow(key)
		allowed2, _, _, _ := limiter.Allow(key)
		allowed3, _, _, _ := limiter.Allow(key)

		assert.True(t, allowed1, "1回目は許可されること")
		assert.True(t, allowed2, "2回目は許可されること")
		assert.False(t, allowed3, "3回目は拒否されること")

		// ウィンドウが経過するまで待機
		time.Sleep(150 * time.Millisecond)

		// ウィンドウ経過後は再度許可される
		allowed4, _, _, _ := limiter.Allow(key)
		assert.True(t, allowed4, "ウィンドウ経過後は再度許可されること")
	})
}

func TestRateLimiter_ThreadSafety(t *testing.T) {
	t.Run("正常系_並行アクセスでも正しく動作", func(t *testing.T) {
		limiter := NewRateLimiter(10, 1*time.Minute)
		key := "test-key"

		done := make(chan bool, 20)
		allows := make(chan bool, 20)

		// 20個のゴルーチンから同時にリクエスト
		for i := 0; i < 20; i++ {
			go func() {
				allowed, _, _, _ := limiter.Allow(key)
				allows <- allowed
				done <- true
			}()
		}

		// 全ゴルーチンの完了を待機
		for i := 0; i < 20; i++ {
			<-done
		}
		close(allows)

		// 許可された回数をカウント
		count := 0
		for allowed := range allows {
			if allowed {
				count++
			}
		}

		assert.Equal(t, 10, count, "制限値（10）だけ許可されること")
	})
}

func TestAuthRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		requests           int
		expectedSuccesses  int
		expectedStatusCode int
	}{
		{
			name:               "正常系_制限内のリクエスト（5回）",
			requests:           5,
			expectedSuccesses:  5,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "異常系_制限超過（6回目）",
			requests:           6,
			expectedSuccesses:  5,
			expectedStatusCode: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			middleware := AuthRateLimitMiddleware()

			handler := func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
			}

			successes := 0
			var lastStatusCode int
			var lastResponse *httptest.ResponseRecorder

			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("X-Forwarded-For", "192.168.1.1")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				h := middleware(handler)
				err := h(c)

				if err == nil && rec.Code == http.StatusOK {
					successes++
				}

				lastStatusCode = rec.Code
				lastResponse = rec
			}

			assert.Equal(t, tt.expectedSuccesses, successes, "成功した回数が期待値と一致すること")
			assert.Equal(t, tt.expectedStatusCode, lastStatusCode, "最後のリクエストのステータスコードが期待値と一致すること")

			// レート制限超過時のヘッダーチェック
			if lastStatusCode == http.StatusTooManyRequests {
				assert.NotEmpty(t, lastResponse.Header().Get("X-RateLimit-Limit"), "X-RateLimit-Limitヘッダーが設定されていること")
				assert.Equal(t, "0", lastResponse.Header().Get("X-RateLimit-Remaining"), "X-RateLimit-Remainingが0であること")
				assert.NotEmpty(t, lastResponse.Header().Get("X-RateLimit-Reset"), "X-RateLimit-Resetヘッダーが設定されていること")
				assert.NotEmpty(t, lastResponse.Header().Get("Retry-After"), "Retry-Afterヘッダーが設定されていること")
			}
		})
	}
}

func TestGeneralRateLimitMiddleware(t *testing.T) {
	t.Run("正常系_60回まで許可", func(t *testing.T) {
		e := echo.New()
		middleware := GeneralRateLimitMiddleware()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}

		successes := 0

		for i := 0; i < 60; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", "192.168.1.1")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := middleware(handler)
			err := h(c)

			if err == nil && rec.Code == http.StatusOK {
				successes++
			}
		}

		assert.Equal(t, 60, successes, "60回まで許可されること")
	})

	t.Run("異常系_61回目は拒否", func(t *testing.T) {
		e := echo.New()
		middleware := GeneralRateLimitMiddleware()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}

		var lastStatusCode int

		for i := 0; i < 61; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", "192.168.1.1")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := middleware(handler)
			_ = h(c)

			lastStatusCode = rec.Code
		}

		assert.Equal(t, http.StatusTooManyRequests, lastStatusCode, "61回目は429が返ること")
	})
}

func TestRateLimitMiddleware_ResponseHeaders(t *testing.T) {
	t.Run("正常系_レスポンスヘッダーが正しく設定される", func(t *testing.T) {
		e := echo.New()
		middleware := AuthRateLimitMiddleware()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := middleware(handler)
		_ = h(c)

		// ヘッダーの存在チェック
		assert.Equal(t, "5", rec.Header().Get("X-RateLimit-Limit"), "X-RateLimit-Limitが5であること")
		assert.Equal(t, "4", rec.Header().Get("X-RateLimit-Remaining"), "X-RateLimit-Remainingが4であること")
		assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"), "X-RateLimit-Resetが設定されていること")

		// X-RateLimit-Reset が Unix タイムスタンプであることを確認
		resetTime, err := strconv.ParseInt(rec.Header().Get("X-RateLimit-Reset"), 10, 64)
		assert.NoError(t, err, "X-RateLimit-ResetがUnixタイムスタンプであること")
		assert.Greater(t, resetTime, time.Now().Unix(), "X-RateLimit-Resetが未来の時刻であること")
	})
}

func TestRateLimitMiddleware_UserIDvsIP(t *testing.T) {
	t.Run("正常系_認証済みユーザーはユーザーID単位でレート制限", func(t *testing.T) {
		e := echo.New()
		middleware := GeneralRateLimitMiddleware()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}

		// ユーザーID "user1" で60回リクエスト
		for i := 0; i < 60; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", "192.168.1.1")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userID", "user1")

			h := middleware(handler)
			_ = h(c)
		}

		// ユーザーID "user2" からのリクエストは別カウント（許可される）
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1") // 同じIP
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userID", "user2")

		h := middleware(handler)
		_ = h(c)

		assert.Equal(t, http.StatusOK, rec.Code, "異なるユーザーIDからのリクエストは許可されること")
	})

	t.Run("正常系_未認証ユーザーはIP単位でレート制限", func(t *testing.T) {
		e := echo.New()
		middleware := GeneralRateLimitMiddleware()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}

		// IP "192.168.1.100" で60回リクエスト
		for i := 0; i < 60; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", "192.168.1.100")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// userIDを設定しない（未認証）

			h := middleware(handler)
			_ = h(c)
		}

		// 異なるIPからのリクエストは別カウント（許可される）
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.200")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := middleware(handler)
		_ = h(c)

		assert.Equal(t, http.StatusOK, rec.Code, "異なるIPからのリクエストは許可されること")
	})
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		forwardedFor   string
		realIP         string
		remoteAddr     string
		expectedIP     string
	}{
		{
			name:         "正常系_X-Forwarded-Forから取得",
			forwardedFor: "203.0.113.1, 198.51.100.1",
			expectedIP:   "203.0.113.1",
		},
		{
			name:       "正常系_X-Real-IPから取得",
			realIP:     "203.0.113.2",
			expectedIP: "203.0.113.2",
		},
		{
			name:       "正常系_RemoteAddrから取得",
			remoteAddr: "203.0.113.3:12345",
			expectedIP: "203.0.113.3",
		},
		{
			name:         "正常系_X-Forwarded-Forが優先",
			forwardedFor: "203.0.113.1",
			realIP:       "203.0.113.2",
			remoteAddr:   "203.0.113.3:12345",
			expectedIP:   "203.0.113.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.forwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.forwardedFor)
			}
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}
			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ip := getClientIP(c)
			assert.Equal(t, tt.expectedIP, ip, "取得したIPアドレスが期待値と一致すること")
		})
	}
}
