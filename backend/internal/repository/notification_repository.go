package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// NotificationRepository 通知リポジトリのインターフェース
type NotificationRepository interface {
	Create(notification *model.Notification) error
	FindByID(id uuid.UUID) (*model.Notification, error)
	FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Notification, int64, error)
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	CountUnread(userID uuid.UUID) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 通知リポジトリのコンストラクタ
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// Create 通知を作成
func (r *notificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// FindByID IDで通知を取得
func (r *notificationRepository) FindByID(id uuid.UUID) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.Preload("Actor").Where("id = ?", id).First(&notification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

// FindByUserID ユーザーIDで通知を取得
func (r *notificationRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("Actor").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkAsRead 通知を既読にする
func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllAsRead 全通知を既読にする
func (r *notificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

// CountUnread 未読通知数をカウント
func (r *notificationRepository) CountUnread(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
