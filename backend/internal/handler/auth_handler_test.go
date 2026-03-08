package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/internal/testutil"
)

// mockPasswordResetService はテスト用のパスワードリセットサービスのモック
type mockPasswordResetService struct{}

func (m *mockPasswordResetService) CreateRequest(email, reason string) (*model.PasswordResetRequest, error) {
	return nil, nil
}

func (m *mockPasswordResetService) ResetPassword(token, newPassword string) error {
	return nil
}

func setupAuthHandlerTest(t *testing.T) (*echo.Echo, *AuthHandler) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := service.NewAuthService(userRepo, refreshTokenRepo)
	passwordResetService := &mockPasswordResetService{}
	authHandler := NewAuthHandler(authService, passwordResetService)

	e := echo.New()
	return e, authHandler
}

func TestAuthHandler_Register(t *testing.T) {
	e, handler := setupAuthHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_有効なデータで登録成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"email": "newuser@example.com",
			"username": "newuser",
			"display_name": "New User",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// レスポンスボディの確認
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotNil(t, response["user"])

		user := response["user"].(map[string]interface{})
		assert.Equal(t, "newuser@example.com", user["email"])
		assert.Equal(t, "newuser", user["username"])
		assert.Equal(t, "New User", user["display_name"])

		// Cookieが設定されていることを確認
		cookies := rec.Result().Cookies()
		var accessTokenCookie, refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				accessTokenCookie = cookie
			}
			if cookie.Name == "refresh_token" {
				refreshTokenCookie = cookie
			}
		}

		assert.NotNil(t, accessTokenCookie, "access_tokenクッキーが設定されているべき")
		assert.NotNil(t, refreshTokenCookie, "refresh_tokenクッキーが設定されているべき")
		assert.NotEmpty(t, accessTokenCookie.Value)
		assert.NotEmpty(t, refreshTokenCookie.Value)
		assert.True(t, accessTokenCookie.HttpOnly, "access_tokenはHttpOnlyであるべき")
		assert.True(t, refreshTokenCookie.HttpOnly, "refresh_tokenはHttpOnlyであるべき")
	})

	t.Run("異常系_必須項目不足", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		testCases := []struct {
			name    string
			reqBody string
		}{
			{
				name:    "emailなし",
				reqBody: `{"username": "user1", "display_name": "User 1", "password": "password123"}`,
			},
			{
				name:    "usernameなし",
				reqBody: `{"email": "user@example.com", "display_name": "User 1", "password": "password123"}`,
			},
			{
				name:    "display_nameなし",
				reqBody: `{"email": "user@example.com", "username": "user1", "password": "password123"}`,
			},
			{
				name:    "passwordなし",
				reqBody: `{"email": "user@example.com", "username": "user1", "display_name": "User 1"}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(tc.reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				err := handler.Register(c)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, rec.Code)

				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				errorObj := response["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "は必須です")
			})
		}
	})

	t.Run("異常系_メールアドレス重複", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// 既存ユーザーを作成
		testutil.CreateTestUser(t, db, "duplicate@example.com", "user1", "User 1", "password123")

		reqBody := `{
			"email": "duplicate@example.com",
			"username": "user2",
			"display_name": "User 2",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		// ユーザー列挙攻撃対策: 409ではなく400を返す
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		// ユーザー列挙攻撃対策: 既存ユーザーか否かを明示しない
		assert.Contains(t, errorObj["message"], "登録に失敗しました")
	})

	t.Run("異常系_ユーザー名重複", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// 既存ユーザーを作成
		testutil.CreateTestUser(t, db, "user1@example.com", "duplicateuser", "User 1", "password123")

		reqBody := `{
			"email": "user2@example.com",
			"username": "duplicateuser",
			"display_name": "User 2",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		// ユーザー列挙攻撃対策: 409ではなく400を返す
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		// ユーザー列挙攻撃対策: 既存ユーザーか否かを明示しない
		assert.Contains(t, errorObj["message"], "登録に失敗しました")
	})

	t.Run("境界値_usernameが3文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"email": "newuser@example.com",
			"username": "abc",
			"display_name": "New User",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("境界値_usernameが2文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"email": "newuser@example.com",
			"username": "ab",
			"display_name": "New User",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "ユーザー名は3文字以上である必要があります")
	})

	t.Run("境界値_usernameが50文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		username := strings.Repeat("a", 50)
		reqBody := `{
			"email": "newuser@example.com",
			"username": "` + username + `",
			"display_name": "New User",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("境界値_usernameが51文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		username := strings.Repeat("a", 51)
		reqBody := `{
			"email": "newuser@example.com",
			"username": "` + username + `",
			"display_name": "New User",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "ユーザー名は50文字以内である必要があります")
	})

	t.Run("境界値_display_nameが100文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		displayName := strings.Repeat("あ", 100)
		reqBody := `{
			"email": "newuser@example.com",
			"username": "newuser",
			"display_name": "` + displayName + `",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("境界値_display_nameが101文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		displayName := strings.Repeat("あ", 101)
		reqBody := `{
			"email": "newuser@example.com",
			"username": "newuser",
			"display_name": "` + displayName + `",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "表示名は100文字以内である必要があります")
	})

	t.Run("境界値_passwordが8文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"email": "newuser@example.com",
			"username": "newuser",
			"display_name": "New User",
			"password": "12345678"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("境界値_passwordが7文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"email": "newuser@example.com",
			"username": "newuser",
			"display_name": "New User",
			"password": "1234567"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "パスワードは8文字以上である必要があります")
	})
}

func TestAuthHandler_Login(t *testing.T) {
	e, handler := setupAuthHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_正しいメール・パスワードでログイン成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		password := "password123"
		user := testutil.CreateTestUser(t, db, "login@example.com", "loginuser", "Login User", password)

		reqBody := `{
			"email": "login@example.com",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスボディの確認
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotNil(t, response["user"])

		userResponse := response["user"].(map[string]interface{})
		assert.Equal(t, user.Email, userResponse["email"])
		assert.Equal(t, user.Username, userResponse["username"])

		// Cookieが設定されていることを確認
		cookies := rec.Result().Cookies()
		var accessTokenCookie, refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				accessTokenCookie = cookie
			}
			if cookie.Name == "refresh_token" {
				refreshTokenCookie = cookie
			}
		}

		assert.NotNil(t, accessTokenCookie, "access_tokenクッキーが設定されているべき")
		assert.NotNil(t, refreshTokenCookie, "refresh_tokenクッキーが設定されているべき")
		assert.NotEmpty(t, accessTokenCookie.Value)
		assert.NotEmpty(t, refreshTokenCookie.Value)
		assert.True(t, accessTokenCookie.HttpOnly, "access_tokenはHttpOnlyであるべき")
		assert.True(t, refreshTokenCookie.HttpOnly, "refresh_tokenはHttpOnlyであるべき")
	})

	t.Run("異常系_必須項目不足", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		testCases := []struct {
			name    string
			reqBody string
		}{
			{
				name:    "emailなし",
				reqBody: `{"password": "password123"}`,
			},
			{
				name:    "passwordなし",
				reqBody: `{"email": "user@example.com"}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(tc.reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				err := handler.Login(c)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, rec.Code)

				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				errorObj := response["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "は必須です")
			})
		}
	})

	t.Run("異常系_メールアドレス不一致", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"email": "notexist@example.com",
			"password": "password123"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "メールアドレスまたはパスワードが正しくありません")
	})

	t.Run("異常系_パスワード不一致", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		testutil.CreateTestUser(t, db, "wrongpass@example.com", "wrongpassuser", "Wrong Pass User", "correctpassword")

		reqBody := `{
			"email": "wrongpass@example.com",
			"password": "wrongpassword"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "メールアドレスまたはパスワードが正しくありません")
	})
}

func TestAuthHandler_GetMe(t *testing.T) {
	e, handler := setupAuthHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_有効なJWTで自分の情報取得", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "getme@example.com", "getmeuser", "GetMe User", "password123")
		token := testutil.GenerateTestJWT(user.ID)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// AuthMiddlewareを模倣してuser_idをセット
		c.Set("user_id", user.ID)

		err := handler.GetMe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, user.Email, response["email"])
		assert.Equal(t, user.Username, response["username"])
		assert.Equal(t, user.DisplayName, response["display_name"])
	})

	t.Run("異常系_JWTなし", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetMe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "認証が必要です")
	})

	t.Run("異常系_無効なJWT", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// AuthMiddlewareを使用して実際の認証フローをテスト
		e := echo.New()
		g := e.Group("/api/v1")
		g.Use(middleware.AuthMiddleware())

		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer invalid-token")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "無効なトークンです")
	})
}
