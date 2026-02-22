package service

import (
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

// AuthService 認証サービスのインターフェース
type AuthService interface {
	Register(req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
	GetMe(userID uuid.UUID) (*response.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService 認証サービスのコンストラクタ
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Register ユーザー登録
func (s *authService) Register(req *request.RegisterRequest) (*response.AuthResponse, error) {
	// メールアドレス重複チェック
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("このメールアドレスは既に使用されています")
	}

	// ユーザー名重複チェック
	existingUser, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("このユーザー名は既に使用されています")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ユーザー作成
	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		PasswordHash: string(hashedPassword),
		Bio:          "",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// JWT生成
	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User:  mapUserToUserResponse(user),
	}, nil
}

// Login ログイン
func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	// ユーザー取得
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// パスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// JWT生成
	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User:  mapUserToUserResponse(user),
	}, nil
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

// generateJWT JWT生成
func generateJWT(userID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // 開発環境用のデフォルト値
	}

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7日間有効
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// mapUserToUserResponse UserをUserResponseにマッピング
func mapUserToUserResponse(user *model.User) response.UserResponse {
	return response.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		HeaderURL:   user.HeaderURL,
		CreatedAt:   user.CreatedAt,
	}
}
