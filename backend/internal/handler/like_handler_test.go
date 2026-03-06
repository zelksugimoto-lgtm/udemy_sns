package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/internal/testutil"
	"gorm.io/gorm"
)

func setupLikeHandlerTest(t *testing.T) (*echo.Echo, *LikeHandler, *gorm.DB) {
	db := testutil.SetupTestDB(t)

	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	mockNotificationService := &MockNotificationService{}

	likeService := service.NewLikeService(likeRepo, postRepo, commentRepo, mockNotificationService)
	likeHandler := NewLikeHandler(likeService)

	e := echo.New()
	return e, likeHandler, db
}

func TestLikeHandler_LikePost(t *testing.T) {
	e, handler, db := setupLikeHandlerTest(t)

	t.Run("正常系_投稿にいいね", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.LikePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("異常系_既にいいね済み", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		// 先にいいねを作成
		like := &model.Like{
			UserID:       user.ID,
			LikeableType: "Post",
			LikeableID:   post.ID,
		}
		db.Create(like)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.LikePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "すでにいいねしています")
	})

	t.Run("異常系_投稿が存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidID := "00000000-0000-0000-0000-000000000000"

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+invalidID+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)
		c.Set("user_id", user.ID)

		err := handler.LikePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "投稿が見つかりません")
	})
}

func TestLikeHandler_UnlikePost(t *testing.T) {
	e, handler, db := setupLikeHandlerTest(t)

	t.Run("正常系_いいね解除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		// 先にいいねを作成
		like := &model.Like{
			UserID:       user.ID,
			LikeableType: "Post",
			LikeableID:   post.ID,
		}
		db.Create(like)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+post.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.UnlikePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系_いいねしていない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+post.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.UnlikePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "いいねしていません")
	})
}

func TestLikeHandler_LikeComment(t *testing.T) {
	e, handler, db := setupLikeHandlerTest(t)

	t.Run("正常系_コメントにいいね", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		comment := testutil.CreateTestComment(t, db, user.ID, post.ID, "This is a test comment")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments/"+comment.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(comment.ID.String())
		c.Set("user_id", user.ID)

		err := handler.LikeComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("異常系_既にいいね済み", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		comment := testutil.CreateTestComment(t, db, user.ID, post.ID, "This is a test comment")

		// 先にいいねを作成
		like := &model.Like{
			UserID:       user.ID,
			LikeableType: "Comment",
			LikeableID:   comment.ID,
		}
		db.Create(like)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments/"+comment.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(comment.ID.String())
		c.Set("user_id", user.ID)

		err := handler.LikeComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "すでにいいねしています")
	})

	t.Run("異常系_コメントが存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidID := "00000000-0000-0000-0000-000000000000"

		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments/"+invalidID+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)
		c.Set("user_id", user.ID)

		err := handler.LikeComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "コメントが見つかりません")
	})
}

func TestLikeHandler_UnlikeComment(t *testing.T) {
	e, handler, db := setupLikeHandlerTest(t)

	t.Run("正常系_コメントいいね解除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		comment := testutil.CreateTestComment(t, db, user.ID, post.ID, "This is a test comment")

		// 先にいいねを作成
		like := &model.Like{
			UserID:       user.ID,
			LikeableType: "Comment",
			LikeableID:   comment.ID,
		}
		db.Create(like)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+comment.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(comment.ID.String())
		c.Set("user_id", user.ID)

		err := handler.UnlikeComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系_いいねしていない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		comment := testutil.CreateTestComment(t, db, user.ID, post.ID, "This is a test comment")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+comment.ID.String()+"/like", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(comment.ID.String())
		c.Set("user_id", user.ID)

		err := handler.UnlikeComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "いいねしていません")
	})
}
