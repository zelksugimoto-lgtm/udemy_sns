package repository

import (
	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// LikeRepository いいねリポジトリのインターフェース
type LikeRepository interface {
	Create(like *model.Like) error
	Delete(userID uuid.UUID, likeableType string, likeableID uuid.UUID) error
	Exists(userID uuid.UUID, likeableType string, likeableID uuid.UUID) (bool, error)
	FindByUserIDAndType(userID uuid.UUID, likeableType string, limit, offset int) ([]model.Like, int64, error)
}

type likeRepository struct {
	db *gorm.DB
}

// NewLikeRepository いいねリポジトリのコンストラクタ
func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// Create いいねを作成
func (r *likeRepository) Create(like *model.Like) error {
	return r.db.Create(like).Error
}

// Delete いいねを削除（物理削除）
func (r *likeRepository) Delete(userID uuid.UUID, likeableType string, likeableID uuid.UUID) error {
	return r.db.Unscoped().Where(
		"user_id = ? AND likeable_type = ? AND likeable_id = ?",
		userID, likeableType, likeableID,
	).Delete(&model.Like{}).Error
}

// Exists いいねが存在するか確認
func (r *likeRepository) Exists(userID uuid.UUID, likeableType string, likeableID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where(
		"user_id = ? AND likeable_type = ? AND likeable_id = ?",
		userID, likeableType, likeableID,
	).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByUserIDAndType ユーザーIDとタイプでいいねを取得
func (r *likeRepository) FindByUserIDAndType(userID uuid.UUID, likeableType string, limit, offset int) ([]model.Like, int64, error) {
	var likes []model.Like
	var total int64

	query := r.db.Model(&model.Like{}).Where("user_id = ? AND likeable_type = ?", userID, likeableType)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&likes).Error

	if err != nil {
		return nil, 0, err
	}

	return likes, total, nil
}
