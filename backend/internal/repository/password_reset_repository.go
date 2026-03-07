package repository

import (
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// PasswordResetRepository パスワードリセット申請リポジトリのインターフェース
type PasswordResetRepository interface {
	Create(request *model.PasswordResetRequest) error
	FindByID(id string) (*model.PasswordResetRequest, error)
	FindByToken(token string) (*model.PasswordResetRequest, error)
	FindByUserID(userID string, limit, offset int) ([]model.PasswordResetRequest, int64, error)
	ListPending(limit, offset int) ([]model.PasswordResetRequest, int64, error)
	List(limit, offset int) ([]model.PasswordResetRequest, int64, error)
	Update(request *model.PasswordResetRequest) error
	Delete(id string) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository パスワードリセット申請リポジトリを生成
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// Create パスワードリセット申請を作成
func (r *passwordResetRepository) Create(request *model.PasswordResetRequest) error {
	return r.db.Create(request).Error
}

// FindByID IDでパスワードリセット申請を取得
func (r *passwordResetRepository) FindByID(id string) (*model.PasswordResetRequest, error) {
	var request model.PasswordResetRequest
	if err := r.db.Preload("User").Preload("ProcessedByAdmin").Where("id = ?", id).First(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// FindByToken トークンでパスワードリセット申請を取得
func (r *passwordResetRepository) FindByToken(token string) (*model.PasswordResetRequest, error) {
	var request model.PasswordResetRequest
	if err := r.db.Preload("User").Where("reset_token = ? AND status = ?", token, "approved").First(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// FindByUserID ユーザーIDでパスワードリセット申請を取得
func (r *passwordResetRepository) FindByUserID(userID string, limit, offset int) ([]model.PasswordResetRequest, int64, error) {
	var requests []model.PasswordResetRequest
	var total int64

	if err := r.db.Model(&model.PasswordResetRequest{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("requested_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListPending 未処理のパスワードリセット申請一覧を取得
func (r *passwordResetRepository) ListPending(limit, offset int) ([]model.PasswordResetRequest, int64, error) {
	var requests []model.PasswordResetRequest
	var total int64

	if err := r.db.Model(&model.PasswordResetRequest{}).Where("status = ?", "pending").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("User").Where("status = ?", "pending").Limit(limit).Offset(offset).Order("requested_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// List 全パスワードリセット申請一覧を取得
func (r *passwordResetRepository) List(limit, offset int) ([]model.PasswordResetRequest, int64, error) {
	var requests []model.PasswordResetRequest
	var total int64

	if err := r.db.Model(&model.PasswordResetRequest{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("User").Preload("ProcessedByAdmin").Limit(limit).Offset(offset).Order("requested_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// Update パスワードリセット申請を更新
func (r *passwordResetRepository) Update(request *model.PasswordResetRequest) error {
	return r.db.Save(request).Error
}

// Delete パスワードリセット申請を削除（論理削除）
func (r *passwordResetRepository) Delete(id string) error {
	return r.db.Delete(&model.PasswordResetRequest{}, "id = ?", id).Error
}
