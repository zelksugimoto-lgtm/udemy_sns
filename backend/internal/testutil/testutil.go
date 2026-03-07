package testutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/config"
	"github.com/yourusername/sns-app/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TestDB テスト用データベース接続を保持
var TestDB *gorm.DB

// SetupTestDB テスト用データベースをセットアップ
func SetupTestDB(t *testing.T) *gorm.DB {
	if TestDB == nil {
		db, err := config.InitDB()
		if err != nil {
			t.Fatalf("テスト用DB接続失敗: %v", err)
		}

		// マイグレーション実行
		if err := config.AutoMigrate(db); err != nil {
			t.Fatalf("テスト用DBマイグレーション失敗: %v", err)
		}

		TestDB = db
		log.Println("テスト用DBセットアップ完了")
	}

	return TestDB
}

// CleanupTestData テスト用データをクリーンアップ
func CleanupTestData(t *testing.T, db *gorm.DB) {
	// 外部キー制約を考慮して逆順に削除
	tables := []string{
		"notifications", "reports", "bookmarks", "likes",
		"comments", "follows", "posts", "refresh_tokens", "users",
	}

	// TRUNCATEを使用して高速かつ確実にデータを削除
	// RESTART IDENTITYでシーケンスもリセット、CASCADEで外部キー制約も処理
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			// TRUNCATEが失敗した場合はDELETEにフォールバック
			if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
				t.Logf("警告: %s テーブルのクリーンアップ失敗: %v", table, err)
			}
		}
	}
}

// CreateTestUser テスト用ユーザーを作成
func CreateTestUser(t *testing.T, db *gorm.DB, email, username, displayName, password string) *model.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("パスワードハッシュ化失敗: %v", err)
	}

	user := &model.User{
		Email:        email,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: string(hashedPassword),
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("テスト用ユーザー作成失敗: %v", err)
	}

	return user
}

// CreateTestPost テスト用投稿を作成
func CreateTestPost(t *testing.T, db *gorm.DB, userID uuid.UUID, content string) *model.Post {
	post := &model.Post{
		UserID:  userID,
		Content: content,
	}

	if err := db.Create(post).Error; err != nil {
		t.Fatalf("テスト用投稿作成失敗: %v", err)
	}

	return post
}

// CreateTestComment テスト用コメントを作成
func CreateTestComment(t *testing.T, db *gorm.DB, userID, postID uuid.UUID, content string) *model.Comment {
	comment := &model.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	if err := db.Create(comment).Error; err != nil {
		t.Fatalf("テスト用コメント作成失敗: %v", err)
	}

	return comment
}

// CreateTestFollow テスト用フォローを作成
func CreateTestFollow(t *testing.T, db *gorm.DB, followerID, followeeID uuid.UUID) *model.Follow {
	follow := &model.Follow{
		FollowerID: followerID,
		FollowedID: followeeID,
	}

	if err := db.Create(follow).Error; err != nil {
		t.Fatalf("テスト用フォロー作成失敗: %v", err)
	}

	return follow
}

// GenerateTestJWT テスト用JWTトークンを生成
func GenerateTestJWT(userID uuid.UUID) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret-key"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(fmt.Sprintf("JWT生成失敗: %v", err))
	}

	return tokenString
}

// MakeRequest HTTPリクエストを作成してテスト実行
func MakeRequest(t *testing.T, e *echo.Echo, method, path, body, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}

// AssertJSONResponse JSONレスポンスを検証
func AssertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(t, expectedStatus, rec.Code, "ステータスコードが一致しません")

	if expectedBody != nil {
		var actual interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &actual)
		assert.NoError(t, err, "JSONパースエラー")

		expectedJSON, _ := json.Marshal(expectedBody)
		var expected interface{}
		json.Unmarshal(expectedJSON, &expected)

		assert.Equal(t, expected, actual, "レスポンスボディが一致しません")
	}
}

// AssertErrorResponse エラーレスポンスを検証（日本語メッセージ）
func AssertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	assert.Equal(t, expectedStatus, rec.Code, "ステータスコードが一致しません")

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "JSONパースエラー")

	message, ok := response["message"].(string)
	assert.True(t, ok, "messageフィールドが存在しません")
	assert.Contains(t, message, expectedMessage, "エラーメッセージが期待した内容を含んでいません")
}
