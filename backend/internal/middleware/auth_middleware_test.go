package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// JWTシークレットを設定
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("正常系_有効なJWTで認証成功", func(t *testing.T) {
		e := echo.New()
		userID := uuid.New()

		// 有効なJWTトークンを生成
		claims := jwt.MapClaims{
			"user_id": userID.String(),
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		// ミドルウェアを通過するハンドラー
		handler := func(c echo.Context) error {
			contextUserID, err := GetUserID(c)
			assert.NoError(t, err)
			assert.Equal(t, userID, contextUserID)
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := AuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系_Authorizationヘッダーなし", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := AuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "認証が必要です")
	})

	t.Run("異常系_Bearer形式でない", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "InvalidFormat token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := AuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "無効な認証形式です")
	})

	t.Run("異常系_無効なトークン", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer invalid.token.here")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := AuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "無効なトークンです")
	})

	t.Run("異常系_有効期限切れトークン", func(t *testing.T) {
		e := echo.New()
		userID := uuid.New()

		// 有効期限切れのJWTトークンを生成
		claims := jwt.MapClaims{
			"user_id": userID.String(),
			"exp":     time.Now().Add(-time.Hour).Unix(), // 1時間前に期限切れ
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := AuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "無効なトークンです")
	})
}

func TestOptionalAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("正常系_有効なJWTで認証成功", func(t *testing.T) {
		e := echo.New()
		userID := uuid.New()

		// 有効なJWTトークンを生成
		claims := jwt.MapClaims{
			"user_id": userID.String(),
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		handler := func(c echo.Context) error {
			contextUserID, err := GetUserID(c)
			assert.NoError(t, err)
			assert.Equal(t, userID, contextUserID)
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := OptionalAuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("正常系_Authorizationヘッダーなしでも次へ進む", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			_, err := GetUserID(c)
			assert.Error(t, err) // user_idが設定されていないのでエラー
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := OptionalAuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("正常系_無効なトークンでも次へ進む", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			_, err := GetUserID(c)
			assert.Error(t, err) // user_idが設定されていないのでエラー
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer invalid.token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := OptionalAuthMiddleware()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("正常系_user_idが設定されている", func(t *testing.T) {
		e := echo.New()
		userID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID)

		result, err := GetUserID(c)

		assert.NoError(t, err)
		assert.Equal(t, userID, result)
	})

	t.Run("異常系_user_idが設定されていない", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result, err := GetUserID(c)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, result)
		assert.Equal(t, "認証が必要です", err.Error())
	})
}
