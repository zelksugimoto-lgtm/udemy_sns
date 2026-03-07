package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// UserManagementService ユーザー管理サービスのインターフェース
type UserManagementService interface {
	ListUsers(status string, limit, offset int) ([]model.User, int64, error)
	GetUserByID(id string) (*model.User, error)
	ApproveUser(userID, adminID string, isRoot bool) error
	RejectUser(userID, adminID string, isRoot bool) error
	GetUserStats() (map[string]int64, error)
}

type userManagementService struct {
	userRepo repository.UserRepository
}

// NewUserManagementService ユーザー管理サービスを生成
func NewUserManagementService(userRepo repository.UserRepository) UserManagementService {
	return &userManagementService{
		userRepo: userRepo,
	}
}

// ListUsers ユーザー一覧を取得
func (s *userManagementService) ListUsers(status string, limit, offset int) ([]model.User, int64, error) {
	return s.userRepo.ListByStatus(status, limit, offset)
}

// GetUserByID IDでユーザーを取得
func (s *userManagementService) GetUserByID(id string) (*model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.userRepo.FindByID(userID)
}

// ApproveUser ユーザーを承認
func (s *userManagementService) ApproveUser(userID, adminID string, isRoot bool) error {
	// ユーザーを取得
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// ステータスを承認に変更
	user.Status = "approved"
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 構造化ログ
	logger := log.Info().
		Str("action", "user_approved").
		Str("target_user_id", userID).
		Str("target_username", user.Username).
		Bool("is_root", isRoot)

	if adminID != "" {
		logger = logger.Str("admin_id", adminID)
	}

	logger.Msg("ユーザー承認")

	return nil
}

// RejectUser ユーザーを拒否
func (s *userManagementService) RejectUser(userID, adminID string, isRoot bool) error {
	// ユーザーを取得
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// ステータスを拒否に変更
	user.Status = "rejected"
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 構造化ログ
	logger := log.Info().
		Str("action", "user_rejected").
		Str("target_user_id", userID).
		Str("target_username", user.Username).
		Bool("is_root", isRoot)

	if adminID != "" {
		logger = logger.Str("admin_id", adminID)
	}

	logger.Msg("ユーザー拒否")

	return nil
}

// GetUserStats ユーザー統計を取得
func (s *userManagementService) GetUserStats() (map[string]int64, error) {
	return s.userRepo.CountByStatus()
}
