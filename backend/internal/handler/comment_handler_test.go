package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/internal/testutil"
)

// MockNotificationService モック通知サービス
type MockNotificationService struct{}

func (m *MockNotificationService) CreateFollowNotification(actorID, targetUserID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) CreateLikeNotification(actorID, targetUserID uuid.UUID, targetType string, targetID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) CreateCommentNotification(actorID, postOwnerID, postID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) CreateReplyNotification(actorID, parentCommentOwnerID, parentCommentID, postID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) DeleteNotificationByAction(actorID, receiverID uuid.UUID, actionType, targetType string, targetID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) GetNotifications(userID uuid.UUID, limit, offset int) (*response.NotificationListResponse, error) {
	return nil, nil
}
func (m *MockNotificationService) MarkAsRead(userID, notificationID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) CreateNotification(notification *model.Notification) error {
	return nil
}
func (m *MockNotificationService) DeleteNotificationsByTarget(targetType string, targetID uuid.UUID) error {
	return nil
}
func (m *MockNotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return 0, nil
}

func setupCommentHandlerTest(t *testing.T) (*echo.Echo, *CommentHandler) {
	db := testutil.SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	mockNotificationService := &MockNotificationService{}

	commentService := service.NewCommentService(commentRepo, postRepo, userRepo, likeRepo, mockNotificationService)
	commentHandler := NewCommentHandler(commentService)

	e := echo.New()
	return e, commentHandler
}

func TestCommentHandler_CreateComment(t *testing.T) {
	e, handler := setupCommentHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_有効なデータでコメント作成", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		reqBody := `{
			"content": "This is a test comment"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/comments", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.CreateComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "This is a test comment", response["content"])
	})

	t.Run("異常系_未認証", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		reqBody := `{
			"content": "This is a test comment"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/comments", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())

		err := handler.CreateComment(c)

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
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		reqBody := `{
			"content": ""
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/comments", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.CreateComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "内容は必須です")
	})

	t.Run("境界値_contentが500文字ちょうど", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		// 500文字ちょうどの文字列
		content := strings.Repeat("あ", 500)
		reqBody := `{
			"content": "` + content + `"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/comments", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.CreateComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("境界値_contentが501文字", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")

		// 501文字の文字列
		content := strings.Repeat("あ", 501)
		reqBody := `{
			"content": "` + content + `"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+post.ID.String()+"/comments", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())
		c.Set("user_id", user.ID)

		err := handler.CreateComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "内容は500文字以内である必要があります")
	})

	t.Run("異常系_投稿が存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidID := "00000000-0000-0000-0000-000000000000"

		reqBody := `{
			"content": "This is a test comment"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+invalidID+"/comments", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)
		c.Set("user_id", user.ID)

		err := handler.CreateComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "投稿が見つかりません")
	})
}

func TestCommentHandler_GetComments(t *testing.T) {
	e, handler := setupCommentHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_投稿のコメント一覧取得", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		testutil.CreateTestComment(t, db, user.ID, post.ID, "Comment 1")
		testutil.CreateTestComment(t, db, user.ID, post.ID, "Comment 2")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.ID.String()+"/comments", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())

		err := handler.GetComments(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		comments := response["comments"].([]interface{})
		assert.GreaterOrEqual(t, len(comments), 2)

		// レスポンスの中身を検証
		firstComment := comments[0].(map[string]interface{})
		assert.NotEmpty(t, firstComment["content"])
		assert.NotNil(t, firstComment["user"])
		commentUser := firstComment["user"].(map[string]interface{})
		assert.Equal(t, user.Username, commentUser["username"])
	})

	t.Run("正常系_返信も含まれる", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		parentComment := testutil.CreateTestComment(t, db, user.ID, post.ID, "Parent comment")

		// 返信コメントを作成
		replyComment := testutil.CreateTestComment(t, db, user.ID, post.ID, "Reply comment")
		replyComment.ParentCommentID = &parentComment.ID
		db.Save(replyComment)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.ID.String()+"/comments", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(post.ID.String())

		err := handler.GetComments(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		comments := response["comments"].([]interface{})
		assert.GreaterOrEqual(t, len(comments), 1)
	})
}

func TestCommentHandler_DeleteComment(t *testing.T) {
	e, handler := setupCommentHandlerTest(t)
	db := testutil.SetupTestDB(t)

	t.Run("正常系_自分のコメントを削除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "This is a test post")
		comment := testutil.CreateTestComment(t, db, user.ID, post.ID, "This is a test comment")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+comment.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(comment.ID.String())
		c.Set("user_id", user.ID)

		err := handler.DeleteComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系_他人のコメントを削除", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "User 1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "User 2", "password123")
		post := testutil.CreateTestPost(t, db, user1.ID, "This is a test post")
		comment := testutil.CreateTestComment(t, db, user1.ID, post.ID, "This is a test comment")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+comment.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(comment.ID.String())
		c.Set("user_id", user2.ID)

		err := handler.DeleteComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "削除する権限がありません")
	})

	t.Run("異常系_コメントが存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "Test User", "password123")
		invalidID := "00000000-0000-0000-0000-000000000000"

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+invalidID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(invalidID)
		c.Set("user_id", user.ID)

		err := handler.DeleteComment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		errorObj := response["error"].(map[string]interface{})
		assert.Contains(t, errorObj["message"], "コメントが見つかりません")
	})
}
