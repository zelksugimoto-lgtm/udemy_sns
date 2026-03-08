package repository

import (
	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// PostMediaRepository 投稿メディアリポジトリのインターフェース
type PostMediaRepository interface {
	CreateBatch(media []model.PostMedia) error
	FindByPostID(postID uuid.UUID) ([]model.PostMedia, error)
	FindByPostIDs(postIDs []uuid.UUID) ([]model.PostMedia, error)
	DeleteByPostID(postID uuid.UUID) error
}

type postMediaRepository struct {
	db *gorm.DB
}

// NewPostMediaRepository PostMediaRepositoryのコンストラクタ
func NewPostMediaRepository(db *gorm.DB) PostMediaRepository {
	return &postMediaRepository{db: db}
}

// CreateBatch 複数のメディアを一括作成
func (r *postMediaRepository) CreateBatch(media []model.PostMedia) error {
	if len(media) == 0 {
		return nil
	}
	return r.db.Create(&media).Error
}

// FindByPostID 投稿IDでメディアを取得（display_order順）
func (r *postMediaRepository) FindByPostID(postID uuid.UUID) ([]model.PostMedia, error) {
	var media []model.PostMedia
	err := r.db.Where("post_id = ?", postID).
		Order("display_order ASC").
		Find(&media).Error
	return media, err
}

// FindByPostIDs 複数の投稿IDでメディアを一括取得
func (r *postMediaRepository) FindByPostIDs(postIDs []uuid.UUID) ([]model.PostMedia, error) {
	var media []model.PostMedia
	if len(postIDs) == 0 {
		return media, nil
	}
	err := r.db.Where("post_id IN ?", postIDs).
		Order("post_id, display_order ASC").
		Find(&media).Error
	return media, err
}

// DeleteByPostID 投稿IDでメディアを削除（論理削除）
func (r *postMediaRepository) DeleteByPostID(postID uuid.UUID) error {
	return r.db.Where("post_id = ?", postID).Delete(&model.PostMedia{}).Error
}
