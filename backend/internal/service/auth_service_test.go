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
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := NewAuthService(userRepo, refreshTokenRepo)

	t.Run("正常系_ユーザー登録成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		req := &request.RegisterRequest{
			Email:       "test@example.com",
			Username:    "testuser",
			DisplayName: "Test User",
			Password:    "password123",
		}

		result, tokenPair, err := authService.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, tokenPair)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
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

		// リフレッシュトークンがDBに保存されていることを確認
		storedToken, err := refreshTokenRepo.FindByToken(tokenPair.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, storedToken)
		assert.Equal(t, user.ID, storedToken.UserID)
		assert.False(t, storedToken.IsRevoked)
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

		result, tokenPair, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Nil(t, tokenPair)
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

		result, tokenPair, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Nil(t, tokenPair)
		assert.Equal(t, "このユーザー名は既に使用されています", err.Error())
	})
}

func TestAuthService_Login(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := NewAuthService(userRepo, refreshTokenRepo)

	t.Run("正常系_ログイン成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		password := "password123"
		user := testutil.CreateTestUser(t, db, "login@example.com", "loginuser", "Login User", password)

		req := &request.LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		result, tokenPair, err := authService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, tokenPair)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
		assert.NotNil(t, result.User)
		assert.Equal(t, user.Email, result.User.Email)
		assert.Equal(t, user.Username, result.User.Username)
		assert.Equal(t, user.DisplayName, result.User.DisplayName)

		// リフレッシュトークンがDBに保存されていることを確認
		storedToken, err := refreshTokenRepo.FindByToken(tokenPair.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, storedToken)
		assert.Equal(t, user.ID, storedToken.UserID)
		assert.False(t, storedToken.IsRevoked)
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		req := &request.LoginRequest{
			Email:    "notexist@example.com",
			Password: "password123",
		}

		result, tokenPair, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Nil(t, tokenPair)
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

		result, tokenPair, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Nil(t, tokenPair)
		assert.Equal(t, "メールアドレスまたはパスワードが正しくありません", err.Error())
	})
}

func TestAuthService_GetMe(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := NewAuthService(userRepo, refreshTokenRepo)

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

func TestAuthService_RefreshToken(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := NewAuthService(userRepo, refreshTokenRepo)

	t.Run("正常系_トークンリフレッシュ成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成してログイン
		user := testutil.CreateTestUser(t, db, "refresh@example.com", "refreshuser", "Refresh User", "password123")
		req := &request.LoginRequest{
			Email:    user.Email,
			Password: "password123",
		}
		_, initialTokenPair, err := authService.Login(req)
		assert.NoError(t, err)

		// リフレッシュトークンを使って新しいトークンペアを取得
		newTokenPair, err := authService.RefreshToken(initialTokenPair.RefreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, newTokenPair)
		assert.NotEmpty(t, newTokenPair.AccessToken)
		assert.NotEmpty(t, newTokenPair.RefreshToken)
		assert.NotEqual(t, initialTokenPair.RefreshToken, newTokenPair.RefreshToken)

		// 古いリフレッシュトークンが無効化されていることを確認
		oldToken, err := refreshTokenRepo.FindByToken(initialTokenPair.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, oldToken)
		assert.True(t, oldToken.IsRevoked)

		// 新しいリフレッシュトークンが有効であることを確認
		newToken, err := refreshTokenRepo.FindByToken(newTokenPair.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, newToken)
		assert.False(t, newToken.IsRevoked)
	})

	t.Run("異常系_無効なリフレッシュトークン", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		newTokenPair, err := authService.RefreshToken("invalid-refresh-token")

		assert.Error(t, err)
		assert.Nil(t, newTokenPair)
		assert.Equal(t, "無効なリフレッシュトークンです", err.Error())
	})

	t.Run("異常系_失効済みリフレッシュトークン", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成してログイン
		user := testutil.CreateTestUser(t, db, "revoked@example.com", "revokeduser", "Revoked User", "password123")
		req := &request.LoginRequest{
			Email:    user.Email,
			Password: "password123",
		}
		_, tokenPair, err := authService.Login(req)
		assert.NoError(t, err)

		// トークンを失効させる
		err = authService.Logout(tokenPair.RefreshToken)
		assert.NoError(t, err)

		// 失効済みトークンでリフレッシュを試みる
		newTokenPair, err := authService.RefreshToken(tokenPair.RefreshToken)

		assert.Error(t, err)
		assert.Nil(t, newTokenPair)
		assert.Equal(t, "リフレッシュトークンの有効期限が切れているか、無効化されています", err.Error())
	})
}

func TestAuthService_Logout(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := NewAuthService(userRepo, refreshTokenRepo)

	t.Run("正常系_ログアウト成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成してログイン
		user := testutil.CreateTestUser(t, db, "logout@example.com", "logoutuser", "Logout User", "password123")
		req := &request.LoginRequest{
			Email:    user.Email,
			Password: "password123",
		}
		_, tokenPair, err := authService.Login(req)
		assert.NoError(t, err)

		// ログアウト
		err = authService.Logout(tokenPair.RefreshToken)

		assert.NoError(t, err)

		// トークンが失効していることを確認
		storedToken, err := refreshTokenRepo.FindByToken(tokenPair.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, storedToken)
		assert.True(t, storedToken.IsRevoked)
	})
}

func TestAuthService_RevokeAllTokens(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := NewAuthService(userRepo, refreshTokenRepo)

	t.Run("正常系_全デバイスログアウト成功", func(t *testing.T) {
		testutil.CleanupTestData(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "revokeall@example.com", "revokealluser", "RevokeAll User", "password123")

		// 複数のリフレッシュトークンを生成（複数デバイスからのログインをシミュレート）
		req := &request.LoginRequest{
			Email:    user.Email,
			Password: "password123",
		}
		_, tokenPair1, err := authService.Login(req)
		assert.NoError(t, err)

		_, tokenPair2, err := authService.Login(req)
		assert.NoError(t, err)

		_, tokenPair3, err := authService.Login(req)
		assert.NoError(t, err)

		// 全デバイスログアウト
		err = authService.RevokeAllTokens(user.ID)

		assert.NoError(t, err)

		// 全てのトークンが失効していることを確認
		token1, err := refreshTokenRepo.FindByToken(tokenPair1.RefreshToken)
		assert.NoError(t, err)
		assert.True(t, token1.IsRevoked)

		token2, err := refreshTokenRepo.FindByToken(tokenPair2.RefreshToken)
		assert.NoError(t, err)
		assert.True(t, token2.IsRevoked)

		token3, err := refreshTokenRepo.FindByToken(tokenPair3.RefreshToken)
		assert.NoError(t, err)
		assert.True(t, token3.IsRevoked)
	})
}
