package repository

import (
	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// FollowRepository フォローリポジトリのインターフェース
type FollowRepository interface {
	Create(follow *model.Follow) error
	Delete(followerID, followedID uuid.UUID) error
	Exists(followerID, followedID uuid.UUID) (bool, error)
	FindFollowersByUserID(userID uuid.UUID, limit, offset int) ([]model.Follow, int64, error)
	FindFollowingByUserID(userID uuid.UUID, limit, offset int) ([]model.Follow, int64, error)
	GetFollowingIDs(userID uuid.UUID) ([]uuid.UUID, error)
}

type followRepository struct {
	db *gorm.DB
}

// NewFollowRepository フォローリポジトリのコンストラクタ
func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

// Create フォローを作成
func (r *followRepository) Create(follow *model.Follow) error {
	return r.db.Create(follow).Error
}

// Delete フォローを削除（物理削除）
func (r *followRepository) Delete(followerID, followedID uuid.UUID) error {
	return r.db.Unscoped().Where(
		"follower_id = ? AND followed_id = ?",
		followerID, followedID,
	).Delete(&model.Follow{}).Error
}

// Exists フォロー関係が存在するか確認
func (r *followRepository) Exists(followerID, followedID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where(
		"follower_id = ? AND followed_id = ?",
		followerID, followedID,
	).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindFollowersByUserID ユーザーIDでフォロワーを取得
func (r *followRepository) FindFollowersByUserID(userID uuid.UUID, limit, offset int) ([]model.Follow, int64, error) {
	var follows []model.Follow
	var total int64

	query := r.db.Model(&model.Follow{}).Where("followed_id = ?", userID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("Follower").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}

// FindFollowingByUserID ユーザーIDでフォロー中を取得
func (r *followRepository) FindFollowingByUserID(userID uuid.UUID, limit, offset int) ([]model.Follow, int64, error) {
	var follows []model.Follow
	var total int64

	query := r.db.Model(&model.Follow{}).Where("follower_id = ?", userID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("Followed").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}

// GetFollowingIDs ユーザーがフォローしているユーザーのIDリストを取得
func (r *followRepository) GetFollowingIDs(userID uuid.UUID) ([]uuid.UUID, error) {
	var followingIDs []uuid.UUID
	err := r.db.Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("followed_id", &followingIDs).Error

	if err != nil {
		return nil, err
	}

	return followingIDs, nil
}
