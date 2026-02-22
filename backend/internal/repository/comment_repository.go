package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// CommentRepository コメントリポジトリのインターフェース
type CommentRepository interface {
	Create(comment *model.Comment) error
	FindByID(id uuid.UUID) (*model.Comment, error)
	FindByPostID(postID uuid.UUID, limit, offset int) ([]model.Comment, int64, error)
	Delete(id uuid.UUID) error
	CountLikes(commentID uuid.UUID) (int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository コメントリポジトリのコンストラクタ
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create コメントを作成
func (r *commentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID IDでコメントを取得
func (r *commentRepository) FindByID(id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").Where("id = ?", id).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// FindByPostID 投稿IDでコメントを取得
func (r *commentRepository) FindByPostID(postID uuid.UUID, limit, offset int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).Where("post_id = ? AND parent_comment_id IS NULL", postID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用（親コメントのみ、子コメントもプリロード）
	err := query.Preload("User").
		Preload("ChildComments").
		Preload("ChildComments.User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// Delete コメントを削除（論理削除）
func (r *commentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// CountLikes コメントのいいね数をカウント
func (r *commentRepository) CountLikes(commentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("likeable_type = ? AND likeable_id = ?", "Comment", commentID).
		Count(&count).Error
	return count, err
}
