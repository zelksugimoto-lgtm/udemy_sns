package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo)

	t.Run("正常系_ユーザー登録成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		req := &request.RegisterRequest{
			Email:       "test@example.com",
			Username:    "testuser",
			DisplayName: "Test User",
			Password:    "password123",
		}

		result, err := authService.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Token)
		assert.NotNil(t, result.User)
		assert.Equal(t, req.Email, result.User.Email)
		assert.Equal(t, req.Username, result.User.Username)
		assert.Equal(t, req.DisplayName, result.User.DisplayName)

		// データベースに保存されたユーザーを確認
		user, err := userRepo.FindByEmail(req.Email)
		assert.NoError(t, err)
		assert.NotNil(t, user)

		// パスワードがハッシュ化されていることを確認
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		assert.NoError(t, err)
	})

	t.Run("異常系_メールアドレス重複", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// 既存ユーザーを作成
		testutil.CreateTestUser(t, db, "duplicate@example.com", "user1", "User 1", "password123")

		req := &request.RegisterRequest{
			Email:       "duplicate@example.com",
			Username:    "user2",
			DisplayName: "User 2",
			Password:    "password123",
		}

		result, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "このメールアドレスは既に使用されています", err.Error())
	})

	t.Run("異常系_ユーザー名重複", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// 既存ユーザーを作成
		testutil.CreateTestUser(t, db, "user1@example.com", "duplicateuser", "User 1", "password123")

		req := &request.RegisterRequest{
			Email:       "user2@example.com",
			Username:    "duplicateuser",
			DisplayName: "User 2",
			Password:    "password123",
		}

		result, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "このユーザー名は既に使用されています", err.Error())
	})
}

func TestAuthService_Login(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo)

	t.Run("正常系_ログイン成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		password := "password123"
		user := testutil.CreateTestUser(t, db, "login@example.com", "loginuser", "Login User", password)

		req := &request.LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		result, err := authService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Token)
		assert.NotNil(t, result.User)
		assert.Equal(t, user.Email, result.User.Email)
		assert.Equal(t, user.Username, result.User.Username)
		assert.Equal(t, user.DisplayName, result.User.DisplayName)
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		req := &request.LoginRequest{
			Email:    "notexist@example.com",
			Password: "password123",
		}

		result, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "メールアドレスまたはパスワードが正しくありません", err.Error())
	})

	t.Run("異常系_パスワードが間違っている", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "wrongpass@example.com", "wrongpassuser", "Wrong Pass User", "correctpassword")

		req := &request.LoginRequest{
			Email:    user.Email,
			Password: "wrongpassword",
		}

		result, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "メールアドレスまたはパスワードが正しくありません", err.Error())
	})
}

func TestAuthService_GetMe(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo)

	t.Run("正常系_ユーザー情報取得成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "getme@example.com", "getmeuser", "GetMe User", "password123")

		result, err := authService.GetMe(user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.DisplayName, result.DisplayName)
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// 存在しないUUIDを生成
		notExistID := testutil.CreateTestUser(t, db, "temp@example.com", "tempuser", "Temp User", "password123").ID
		db.Exec("DELETE FROM users WHERE id = ?", notExistID)

		result, err := authService.GetMe(notExistID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "ユーザーが見つかりません", err.Error())
	})
}
