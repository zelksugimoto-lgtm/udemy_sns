package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// RefreshTokenRepository リフレッシュトークンリポジトリのインターフェース
type RefreshTokenRepository interface {
	Create(token *model.RefreshToken) error
	FindByToken(token string) (*model.RefreshToken, error)
	RevokeByToken(token string) error
	RevokeAllByUserID(userID uuid.UUID) error
	DeleteExpired() error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository リフレッシュトークンリポジトリのコンストラクタ
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create リフレッシュトークンを作成
func (r *refreshTokenRepository) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

// FindByToken トークンで検索
func (r *refreshTokenRepository) FindByToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

// RevokeByToken トークンを無効化
func (r *refreshTokenRepository) RevokeByToken(token string) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

// RevokeAllByUserID ユーザーの全リフレッシュトークンを無効化
func (r *refreshTokenRepository) RevokeAllByUserID(userID uuid.UUID) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

// DeleteExpired 期限切れのトークンを削除
func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Unscoped().
		Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshToken{}).Error
}
