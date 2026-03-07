package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	// AccessTokenDuration アクセストークンの有効期限（1時間）
	AccessTokenDuration = 1 * time.Hour
	// RefreshTokenDuration リフレッシュトークンの有効期限（7日間）
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// TokenPair アクセストークンとリフレッシュトークンのペア
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// AuthService 認証サービスのインターフェース
type AuthService interface {
	Register(req *request.RegisterRequest) (*response.AuthResponse, *TokenPair, error)
	Login(req *request.LoginRequest) (*response.AuthResponse, *TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	Logout(refreshToken string) error
	RevokeAllTokens(userID uuid.UUID) error
	GetMe(userID uuid.UUID) (*response.UserResponse, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

// NewAuthService 認証サービスのコンストラクタ
func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// Register ユーザー登録
func (s *authService) Register(req *request.RegisterRequest) (*response.AuthResponse, *TokenPair, error) {
	// メールアドレス重複チェック
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, errors.New("このメールアドレスは既に使用されています")
	}

	// ユーザー名重複チェック
	existingUser, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, errors.New("このユーザー名は既に使用されています")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	// ユーザー作成
	// E2Eテスト環境では自動承認、それ以外はpending
	status := "pending"
	if os.Getenv("AUTO_APPROVE_USERS") == "true" {
		status = "approved"
	}

	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		PasswordHash: string(hashedPassword),
		Bio:          "",
		Status:       status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// トークン生成
	tokenPair, err := s.generateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return &response.AuthResponse{
		User: mapUserToUserResponse(user),
	}, tokenPair, nil
}

// Login ログイン
func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, *TokenPair, error) {
	// ユーザー取得
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// パスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// トークン生成
	tokenPair, err := s.generateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return &response.AuthResponse{
		User: mapUserToUserResponse(user),
	}, tokenPair, nil
}

// RefreshToken リフレッシュトークンでアクセストークンを再発行
func (s *authService) RefreshToken(refreshToken string) (*TokenPair, error) {
	// リフレッシュトークンをDBから取得
	storedToken, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, errors.New("無効なリフレッシュトークンです")
	}

	// トークンの有効性をチェック
	if !storedToken.IsValid() {
		return nil, errors.New("リフレッシュトークンの有効期限が切れているか、無効化されています")
	}

	// 古いリフレッシュトークンを無効化（トークンローテーション）
	if err := s.refreshTokenRepo.RevokeByToken(refreshToken); err != nil {
		return nil, err
	}

	// 新しいトークンペアを生成
	tokenPair, err := s.generateTokenPair(storedToken.UserID)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// Logout ログアウト（リフレッシュトークンを無効化）
func (s *authService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.RevokeByToken(refreshToken)
}

// RevokeAllTokens 全デバイスログアウト（ユーザーの全リフレッシュトークンを無効化）
func (s *authService) RevokeAllTokens(userID uuid.UUID) error {
	return s.refreshTokenRepo.RevokeAllByUserID(userID)
}

// GetMe ログイン中のユーザー情報取得
func (s *authService) GetMe(userID uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	return mapUserToUserResponse(user), nil
}

// ValidateAccessToken アクセストークンを検証
func (s *authService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	secret := getJWTSecret()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無効な署名アルゴリズムです")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("無効なトークンです")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("トークンのクレームが不正です")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("トークンにuser_idが含まれていません")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("無効なuser_idです")
	}

	return userID, nil
}

// generateTokenPair アクセストークンとリフレッシュトークンのペアを生成
func (s *authService) generateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	// アクセストークン生成
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークン生成
	refreshTokenString, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンをDBに保存
	refreshToken := &model.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(RefreshTokenDuration),
		IsRevoked: false,
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

// generateAccessToken アクセストークンを生成（1時間有効、user_idとexpのみ）
func generateAccessToken(userID uuid.UUID) (string, error) {
	secret := getJWTSecret()

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateRandomToken ランダムなリフレッシュトークンを生成
func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// getJWTSecret JWT秘密鍵を取得
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // 開発環境用のデフォルト値
	}
	return secret
}

// mapUserToUserResponse UserをUserResponseにマッピング
func mapUserToUserResponse(user *model.User) *response.UserResponse {
	return &response.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		HeaderURL:   user.HeaderURL,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
	}
}
