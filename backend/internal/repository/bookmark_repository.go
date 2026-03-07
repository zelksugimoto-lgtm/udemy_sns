package repository

import (
	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// BookmarkRepository ブックマークリポジトリのインターフェース
type BookmarkRepository interface {
	Create(bookmark *model.Bookmark) error
	Delete(userID, postID uuid.UUID) error
	Exists(userID, postID uuid.UUID) (bool, error)
	ExistsInBatch(userID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error)
	FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error)
}

type bookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository ブックマークリポジトリのコンストラクタ
func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepository{db: db}
}

// Create ブックマークを作成
func (r *bookmarkRepository) Create(bookmark *model.Bookmark) error {
	return r.db.Create(bookmark).Error
}

// Delete ブックマークを削除（物理削除）
func (r *bookmarkRepository) Delete(userID, postID uuid.UUID) error {
	return r.db.Unscoped().Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Bookmark{}).Error
}

// Exists ブックマークが存在するか確認
func (r *bookmarkRepository) Exists(userID, postID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Bookmark{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsInBatch 複数のブックマークが存在するか一括確認
func (r *bookmarkRepository) ExistsInBatch(userID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	var bookmarks []model.Bookmark
	err := r.db.Model(&model.Bookmark{}).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Select("post_id").
		Find(&bookmarks).Error

	if err != nil {
		return nil, err
	}

	// マップに変換
	exists := make(map[uuid.UUID]bool)
	for _, bookmark := range bookmarks {
		exists[bookmark.PostID] = true
	}

	return exists, nil
}

// FindByUserID ユーザーIDでブックマークを取得
func (r *bookmarkRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error) {
	var bookmarks []model.Bookmark
	var total int64

	query := r.db.Model(&model.Bookmark{}).Where("user_id = ?", userID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("Post").
		Preload("Post.User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&bookmarks).Error

	if err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}
