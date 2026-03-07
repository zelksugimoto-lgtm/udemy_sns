package service

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理者サービスのインターフェース
type AdminService interface {
	// 認証
	AuthenticateRoot(username, password string) (bool, error)
	AuthenticateAdmin(username, password string) (*model.Admin, error)

	// 管理者管理
	CreateAdmin(username, email, password, role string) (*model.Admin, error)
	GetAdminByID(id string) (*model.Admin, error)
	ListAdmins(limit, offset int) ([]model.Admin, int64, error)
	UpdateAdmin(admin *model.Admin) error
	DeleteAdmin(id string) error
}

type adminService struct {
	adminRepo repository.AdminRepository
}

// NewAdminService 管理者サービスを生成
func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

// AuthenticateRoot ルート管理者認証
func (s *adminService) AuthenticateRoot(username, password string) (bool, error) {
	rootUsername := os.Getenv("ADMIN_ROOT_USERNAME")
	rootPassword := os.Getenv("ADMIN_ROOT_PASSWORD")

	if rootUsername == "" || rootPassword == "" {
		log.Warn().Msg("ルート管理者の認証情報が環境変数に設定されていません")
		return false, errors.New("ルート管理者認証が設定されていません")
	}

	// ユーザー名とパスワードをチェック
	if username == rootUsername && password == rootPassword {
		log.Info().Str("username", username).Msg("ルート管理者ログイン成功")
		return true, nil
	}

	return false, nil
}

// AuthenticateAdmin 管理者認証
func (s *adminService) AuthenticateAdmin(username, password string) (*model.Admin, error) {
	// ユーザー名で管理者を取得
	admin, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// パスワードチェック
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, nil
	}

	log.Info().Str("admin_id", admin.ID.String()).Str("username", admin.Username).Msg("管理者ログイン成功")
	return admin, nil
}

// CreateAdmin 管理者を作成
func (s *adminService) CreateAdmin(username, email, password, role string) (*model.Admin, error) {
	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.adminRepo.Create(admin); err != nil {
		return nil, err
	}

	log.Info().
		Str("admin_id", admin.ID.String()).
		Str("username", admin.Username).
		Str("role", admin.Role).
		Msg("管理者作成成功")

	return admin, nil
}

// GetAdminByID IDで管理者を取得
func (s *adminService) GetAdminByID(id string) (*model.Admin, error) {
	return s.adminRepo.FindByID(id)
}

// ListAdmins 管理者一覧を取得
func (s *adminService) ListAdmins(limit, offset int) ([]model.Admin, int64, error) {
	return s.adminRepo.List(limit, offset)
}

// UpdateAdmin 管理者を更新
func (s *adminService) UpdateAdmin(admin *model.Admin) error {
	admin.UpdatedAt = time.Now()
	return s.adminRepo.Update(admin)
}

// DeleteAdmin 管理者を削除
func (s *adminService) DeleteAdmin(id string) error {
	log.Info().Str("admin_id", id).Msg("管理者削除")
	return s.adminRepo.Delete(id)
}
