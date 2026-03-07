package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// PostRepository 投稿リポジトリのインターフェース
type PostRepository interface {
	Create(post *model.Post) error
	FindByID(id uuid.UUID) (*model.Post, error)
	FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Post, int64, error)
	FindLikedByUser(userID uuid.UUID, limit, offset int) ([]model.Post, int64, error)
	Update(post *model.Post) error
	Delete(id uuid.UUID) error
	GetTimeline(userID uuid.UUID, followingIDs []uuid.UUID, limit, offset int) ([]model.Post, int64, error)
	CountLikes(postID uuid.UUID) (int64, error)
	CountComments(postID uuid.UUID) (int64, error)
	CountLikesInBatch(postIDs []uuid.UUID) (map[uuid.UUID]int64, error)
	CountCommentsInBatch(postIDs []uuid.UUID) (map[uuid.UUID]int64, error)
}

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository 投稿リポジトリのコンストラクタ
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Create 投稿を作成
func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// FindByID IDで投稿を取得
func (r *postRepository) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").Where("id = ?", id).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

// FindByUserID ユーザーIDで投稿を取得
func (r *postRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("user_id = ?", userID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Update 投稿を更新
func (r *postRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete 投稿を削除（論理削除）
func (r *postRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Post{}, id).Error
}

// GetTimeline タイムラインを取得
func (r *postRepository) GetTimeline(userID uuid.UUID, followingIDs []uuid.UUID, limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	var query *gorm.DB

	// Phase 1: フォローIDが空の場合は全ユーザーの公開投稿を取得
	if len(followingIDs) == 0 {
		query = r.db.Model(&model.Post{}).Where("visibility = ?", "public")
	} else {
		// Phase 2以降: 自分とフォローしているユーザーの投稿を取得
		userIDs := append(followingIDs, userID)
		query = r.db.Model(&model.Post{}).Where("user_id IN ? AND visibility IN ?", userIDs, []string{"public", "followers"})
	}

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// CountLikes 投稿のいいね数をカウント
func (r *postRepository) CountLikes(postID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("likeable_type = ? AND likeable_id = ?", "Post", postID).
		Count(&count).Error
	return count, err
}

// CountComments 投稿のコメント数をカウント
func (r *postRepository) CountComments(postID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// FindLikedByUser ユーザーがいいねした投稿を取得
func (r *postRepository) FindLikedByUser(userID uuid.UUID, limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	// いいねテーブルとJOINして、該当ユーザーがいいねした投稿を取得
	query := r.db.Model(&model.Post{}).
		Joins("INNER JOIN likes ON likes.likeable_id = posts.id AND likes.likeable_type = ?", "Post").
		Where("likes.user_id = ?", userID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("User").
		Order("likes.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// CountLikesInBatch 複数の投稿のいいね数を一括取得
func (r *postRepository) CountLikesInBatch(postIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	type Result struct {
		LikeableID uuid.UUID `gorm:"column:likeable_id"`
		Count      int64     `gorm:"column:count"`
	}

	var results []Result
	err := r.db.Model(&model.Like{}).
		Select("likeable_id, COUNT(*) as count").
		Where("likeable_type = ? AND likeable_id IN ?", "Post", postIDs).
		Group("likeable_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// マップに変換
	counts := make(map[uuid.UUID]int64)
	for _, result := range results {
		counts[result.LikeableID] = result.Count
	}

	return counts, nil
}

// CountCommentsInBatch 複数の投稿のコメント数を一括取得
func (r *postRepository) CountCommentsInBatch(postIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	type Result struct {
		PostID uuid.UUID `gorm:"column:post_id"`
		Count  int64     `gorm:"column:count"`
	}

	var results []Result
	err := r.db.Model(&model.Comment{}).
		Select("post_id, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// マップに変換
	counts := make(map[uuid.UUID]int64)
	for _, result := range results {
		counts[result.PostID] = result.Count
	}

	return counts, nil
}
