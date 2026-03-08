package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/internal/testutil"
)

func setupPostHandlerTest(t *testing.T) (*echo.Echo, *PostHandler) {
	db := testutil.SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	postMediaRepo := repository.NewPostMediaRepository(db)
	followRepo := repository.NewFollowRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	storageService := service.NewStorageService("test-bucket")
	postService := service.NewPostService(postRepo, postMediaRepo, userRepo, followRepo, likeRepo, bookmarkRepo, storageService)
	postHandler := NewPostHandler(postService)

	e := echo.New()
	return e, postHandler
}

func TestPostHandler_CreatePost(t *testing.T) {
	e, handler := setupPostHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_有効なデータで投稿作成", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		token := testutil.GenerateTestJWT(user.ID)

		reqBody := `{
			"content": "This is a test post",
			"visibility": "public"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "This is a test post", response["content"])
		assert.Equal(t, "public", response["visibility"])
	})

	t.Run("異常系_未認証", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		reqBody := `{
			"content": "This is a test post",
			"visibility": "public"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "認証が必要です")
	})

	t.Run("異常系_contentが空", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")

		reqBody := `{
			"content": "",
			"visibility": "public"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "内容は必須です")
	})

	t.Run("境界値_contentが250文字ちょうど", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")

		// 250文字ちょうどの文字列
		content := strings.Repeat("あ", 250)
		reqBody := `{
			"content": "` + content + `",
			"visibility": "public"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("境界値_contentが251文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")

		// 251文字の文字列
		content := strings.Repeat("あ", 251)
		reqBody := `{
			"content": "` + content + `",
			"visibility": "public"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "内容は250文字以内である必要があります")
	})
}

func TestPostHandler_GetPost(t *testing.T) {
	e, handler := setupPostHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_投稿が存在する", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())

		err := handler.GetPost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "This is a test post", response["content"])
	})

	t.Run("異常系_投稿が存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		invalidID := "00000000-0000-0000-0000-000000000000"

		req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+invalidID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)

		err := handler.GetPost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "投稿が見つかりません")
	})
}

func TestPostHandler_UpdatePost(t *testing.T) {
	e, handler := setupPostHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_自分の投稿を更新", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Original content")

		reqBody := `{
			"content": "Updated content"
		}`

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/posts/"+post.ID.String(), strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.UpdatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "Updated content", response["content"])
	})

	t.Run("異常系_他人の投稿を更新", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")
		post := testutil.CreateTestPost(t, db, user1.ID, "Original content")

		reqBody := `{
			"content": "Updated content"
		}`

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/posts/"+post.ID.String(), strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user2.ID)

		err := handler.UpdatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "編集する権限がありません")
	})

	t.Run("異常系_投稿が存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidID := "00000000-0000-0000-0000-000000000000"

		reqBody := `{
			"content": "Updated content"
		}`

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/posts/"+invalidID, strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)
		c.Set("user_id", user.ID)

		err := handler.UpdatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "投稿が見つかりません")
	})
}

func TestPostHandler_DeletePost(t *testing.T) {
	e, handler := setupPostHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_自分の投稿を削除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+post.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.DeletePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系_他人の投稿を削除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")
		post := testutil.CreateTestPost(t, db, user1.ID, "This is a test post")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+post.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user2.ID)

		err := handler.DeletePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "削除する権限がありません")
	})

	t.Run("異常系_投稿が存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidID := "00000000-0000-0000-0000-000000000000"

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+invalidID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)
		c.Set("user_id", user.ID)

		err := handler.DeletePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "投稿が見つかりません")
	})
}

func TestPostHandler_GetTimeline(t *testing.T) {
	e, handler := setupPostHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_フォローしているユーザーの投稿を取得", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")

		// user1がuser2をフォロー
		testutil.CreateTestFollow(t, db, user1.ID, user2.ID)

		// user2の投稿を作成
		testutil.CreateTestPost(t, db, user2.ID, "Post by user2")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/timeline", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user1.ID)

		err := handler.GetTimeline(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		posts := response["posts"].([]interface{})
		assert.GreaterOrEqual(t, len(posts), 1)

		// レスポンスの中身を検証
		firstPost := posts[0].(map[string]interface{})
		assert.Equal(t, "Post by user2", firstPost["content"])
		assert.NotNil(t, firstPost["user"])
		postUser := firstPost["user"].(map[string]interface{})
		assert.Equal(t, user2.Username, postUser["username"])
		assert.Equal(t, user2.DisplayName, postUser["display_name"])
	})

	t.Run("正常系_自分の投稿も含まれる", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")

		// 自分の投稿を作成
		testutil.CreateTestPost(t, db, user.ID, "My own post")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/timeline", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handler.GetTimeline(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		posts := response["posts"].([]interface{})
		assert.GreaterOrEqual(t, len(posts), 1)

		// レスポンスの中身を検証
		firstPost := posts[0].(map[string]interface{})
		assert.Equal(t, "My own post", firstPost["content"])
		assert.NotNil(t, firstPost["user"])
		postUser := firstPost["user"].(map[string]interface{})
		assert.Equal(t, user.Username, postUser["username"])
	})
}
