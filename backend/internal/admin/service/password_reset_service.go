package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("ユーザーが見つかりません")
	ErrInvalidToken    = errors.New("無効なトークンです")
	ErrTokenExpired    = errors.New("トークンの有効期限が切れています")
	ErrTokenUsed       = errors.New("このトークンは既に使用されています")
)

// PasswordResetService パスワードリセットサービスのインターフェース
type PasswordResetService interface {
	CreateRequest(email, reason string) (*model.PasswordResetRequest, error)
	ListRequests(status string, limit, offset int) ([]model.PasswordResetRequest, int64, error)
	GetRequestByID(id string) (*model.PasswordResetRequest, error)
	GenerateResetLink(requestID, adminID string, isRoot bool) (string, error)
	ResetPassword(token, newPassword string) error
	RejectRequest(requestID, adminID string, isRoot bool) error
}

type passwordResetService struct {
	resetRepo   repository.PasswordResetRepository
	userRepo    repository.UserRepository
	frontendURL string
}

// NewPasswordResetService パスワードリセットサービスを生成
func NewPasswordResetService(
	resetRepo repository.PasswordResetRepository,
	userRepo repository.UserRepository,
	frontendURL string,
) PasswordResetService {
	return &passwordResetService{
		resetRepo:   resetRepo,
		userRepo:    userRepo,
		frontendURL: frontendURL,
	}
}

// CreateRequest パスワードリセット申請を作成
func (s *passwordResetService) CreateRequest(email, reason string) (*model.PasswordResetRequest, error) {
	// メールアドレスでユーザーを検索
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("メールアドレスが見つかりません")
	}
	if user == nil {
		return nil, errors.New("メールアドレスが見つかりません")
	}

	// 申請を作成
	request := &model.PasswordResetRequest{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:      user.ID.String(),
		Reason:      reason,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.resetRepo.Create(request); err != nil {
		return nil, err
	}

	log.Info().
		Str("action", "password_reset_request_created").
		Str("user_id", user.ID.String()).
		Str("email", email).
		Msg("パスワードリセット申請作成")

	return request, nil
}

// ListRequests パスワードリセット申請一覧を取得
func (s *passwordResetService) ListRequests(status string, limit, offset int) ([]model.PasswordResetRequest, int64, error) {
	if status == "pending" {
		return s.resetRepo.ListPending(limit, offset)
	}
	return s.resetRepo.List(limit, offset)
}

// GetRequestByID IDでパスワードリセット申請を取得
func (s *passwordResetService) GetRequestByID(id string) (*model.PasswordResetRequest, error) {
	return s.resetRepo.FindByID(id)
}

// GenerateResetLink パスワードリセットリンクを生成
func (s *passwordResetService) GenerateResetLink(requestID, adminID string, isRoot bool) (string, error) {
	// 申請を取得
	request, err := s.resetRepo.FindByID(requestID)
	if err != nil {
		return "", err
	}
	if request == nil {
		return "", errors.New("申請が見つかりません")
	}

	// トークン生成
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// 有効期限: 24時間
	expiresAt := time.Now().Add(24 * time.Hour)

	// 申請を更新
	request.Status = "approved"
	request.ResetToken = &token
	request.TokenExpiresAt = &expiresAt
	now := time.Now()
	request.ProcessedAt = &now
	if adminID != "" {
		request.ProcessedBy = &adminID
	}
	request.UpdatedAt = time.Now()

	if err := s.resetRepo.Update(request); err != nil {
		return "", err
	}

	// 構造化ログ
	logger := log.Info().
		Str("action", "password_reset_link_generated").
		Str("request_id", requestID).
		Str("user_id", request.UserID).
		Bool("is_root", isRoot)

	if adminID != "" {
		logger = logger.Str("admin_id", adminID)
	}

	logger.Msg("パスワードリセットリンク発行")

	// リセットリンクを返す
	return s.frontendURL + "/password-reset/reset?token=" + token, nil
}

// ResetPassword トークンを使ってパスワードをリセット
func (s *passwordResetService) ResetPassword(token, newPassword string) error {
	// トークンで申請を取得
	request, err := s.resetRepo.FindByToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	// トークンの有効性をチェック
	if !request.IsTokenValid() {
		return ErrTokenExpired
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	// ユーザーのパスワードを更新
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// トークンを無効化（使用済み）
	request.Status = "used"
	request.ResetToken = nil
	request.UpdatedAt = time.Now()
	if err := s.resetRepo.Update(request); err != nil {
		return err
	}

	log.Info().
		Str("action", "password_reset").
		Str("user_id", request.UserID).
		Msg("パスワードリセット完了")

	return nil
}

// RejectRequest パスワードリセット申請を拒否
func (s *passwordResetService) RejectRequest(requestID, adminID string, isRoot bool) error {
	request, err := s.resetRepo.FindByID(requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("申請が見つかりません")
	}

	request.Status = "rejected"
	now := time.Now()
	request.ProcessedAt = &now
	if adminID != "" {
		request.ProcessedBy = &adminID
	}
	request.UpdatedAt = time.Now()

	if err := s.resetRepo.Update(request); err != nil {
		return err
	}

	// 構造化ログ
	logger := log.Info().
		Str("action", "password_reset_request_rejected").
		Str("request_id", requestID).
		Str("user_id", request.UserID).
		Bool("is_root", isRoot)

	if adminID != "" {
		logger = logger.Str("admin_id", adminID)
	}

	logger.Msg("パスワードリセット申請拒否")

	return nil
}
