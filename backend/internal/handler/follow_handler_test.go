package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/internal/testutil"
)

func setupFollowHandlerTest(t *testing.T) (*echo.Echo, *FollowHandler) {
	db := testutil.SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	followRepo := repository.NewFollowRepository(db)
	mockNotificationService := &MockNotificationService{}

	followService := service.NewFollowService(followRepo, userRepo, mockNotificationService)
	followHandler := NewFollowHandler(followService)

	e := echo.New()
	return e, followHandler
}

func TestFollowHandler_Follow(t *testing.T) {
	e, handler := setupFollowHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_ユーザーをフォロー", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/"+user2.Username+"/follow", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user2.Username)
		c.Set("user_id", user1.ID)

		err := handler.Follow(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("異常系_既にフォロー済み", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")

		// 先にフォローを作成
		testutil.CreateTestFollow(t, db, user1.ID, user2.ID)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/"+user2.Username+"/follow", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user2.Username)
		c.Set("user_id", user1.ID)

		err := handler.Follow(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "すでにフォローしています")
	})

	t.Run("異常系_自分自身をフォロー", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/"+user.Username+"/follow", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user.Username)
		c.Set("user_id", user.ID)

		err := handler.Follow(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "自分自身をフォローすることはできません")
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidUsername := "notexistuser"

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/"+invalidUsername+"/follow", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(invalidUsername)
		c.Set("user_id", user.ID)

		err := handler.Follow(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "ユーザーが見つかりません")
	})
}

func TestFollowHandler_Unfollow(t *testing.T) {
	e, handler := setupFollowHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_フォロー解除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")

		// 先にフォローを作成
		testutil.CreateTestFollow(t, db, user1.ID, user2.ID)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user2.Username+"/follow", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user2.Username)
		c.Set("user_id", user1.ID)

		err := handler.Unfollow(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系_フォローしていない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user2.Username+"/follow", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user2.Username)
		c.Set("user_id", user1.ID)

		err := handler.Unfollow(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "フォローしていません")
	})
}

func TestFollowHandler_GetFollowers(t *testing.T) {
	e, handler := setupFollowHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_フォロワー一覧取得", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")
		user3 := testutil.CreateTestUser(t, db, "user3@example.com", "user3", "User 3", "password123")

		// user2とuser3がuser1をフォロー
		testutil.CreateTestFollow(t, db, user2.ID, user1.ID)
		testutil.CreateTestFollow(t, db, user3.ID, user1.ID)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+user1.Username+"/followers", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user1.Username)

		err := handler.GetFollowers(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		users := response["users"].([]interface{})
		assert.GreaterOrEqual(t, len(users), 2)

		// レスポンスの中身を検証
		firstUser := users[0].(map[string]interface{})
		assert.NotEmpty(t, firstUser["username"])
		assert.NotEmpty(t, firstUser["display_name"])
	})
}

func TestFollowHandler_GetFollowing(t *testing.T) {
	e, handler := setupFollowHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_フォロー中一覧取得", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")
		user3 := testutil.CreateTestUser(t, db, "user3@example.com", "user3", "User 3", "password123")

		// user1がuser2とuser3をフォロー
		testutil.CreateTestFollow(t, db, user1.ID, user2.ID)
		testutil.CreateTestFollow(t, db, user1.ID, user3.ID)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+user1.Username+"/following", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(user1.Username)

		err := handler.GetFollowing(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		users := response["users"].([]interface{})
		assert.GreaterOrEqual(t, len(users), 2)

		// レスポンスの中身を検証
		firstUser := users[0].(map[string]interface{})
		assert.NotEmpty(t, firstUser["username"])
		assert.NotEmpty(t, firstUser["display_name"])
	})
}
