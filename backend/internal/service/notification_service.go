package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// NotificationService 通知サービスのインターフェース
type NotificationService interface {
	GetNotifications(userID uuid.UUID, limit, offset int) (*response.NotificationListResponse, error)
	MarkAsRead(userID uuid.UUID, notificationID uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	CreateNotification(notification *model.Notification) error

	// 通知生成ヘルパー
	CreateLikeNotification(actorID, targetUserID uuid.UUID, targetType string, targetID uuid.UUID) error
	CreateCommentNotification(actorID, postOwnerID uuid.UUID, postID uuid.UUID) error
	CreateReplyNotification(actorID, parentCommentOwnerID uuid.UUID, parentCommentID uuid.UUID, postID uuid.UUID) error
	CreateFollowNotification(actorID, targetUserID uuid.UUID) error

	// 通知削除ヘルパー
	DeleteNotificationByAction(actorID, userID uuid.UUID, notifType, targetType string, targetID uuid.UUID) error
	DeleteNotificationsByTarget(targetType string, targetID uuid.UUID) error

	// 未読数取得
	GetUnreadCount(userID uuid.UUID) (int64, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

// NewNotificationService 通知サービスのコンストラクタ
func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

// GetNotifications 通知一覧取得
func (s *notificationService) GetNotifications(userID uuid.UUID, limit, offset int) (*response.NotificationListResponse, error) {
	notifications, total, err := s.notificationRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// 未読数取得
	unreadCount, err := s.notificationRepo.CountUnread(userID)
	if err != nil {
		return nil, err
	}

	// NotificationResponseに変換
	notificationResponses := make([]response.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		var actor *response.UserSimple
		if notification.Actor != nil {
			actor = &response.UserSimple{
				ID:          notification.Actor.ID,
				Username:    notification.Actor.Username,
				DisplayName: notification.Actor.DisplayName,
				AvatarURL:   notification.Actor.AvatarURL,
			}
		}

		notificationResponses[i] = response.NotificationResponse{
			ID:         notification.ID,
			Type:       notification.Type,
			Message:    notification.Message,
			IsRead:     notification.IsRead,
			Actor:      actor,
			TargetType: notification.TargetType,
			TargetID:   notification.TargetID,
			PostID:     notification.PostID,
			CreatedAt:  notification.CreatedAt,
		}
	}

	hasMore := offset+limit < int(total)

	return &response.NotificationListResponse{
		Notifications: notificationResponses,
		UnreadCount:   int(unreadCount),
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}

// MarkAsRead 通知を既読にする
func (s *notificationService) MarkAsRead(userID uuid.UUID, notificationID uuid.UUID) error {
	notification, err := s.notificationRepo.FindByID(notificationID)
	if err != nil {
		return err
	}
	if notification == nil {
		return errors.New("通知が見つかりません")
	}

	// 権限チェック
	if notification.UserID != userID {
		return errors.New("この通知を既読にする権限がありません")
	}

	return s.notificationRepo.MarkAsRead(notificationID)
}

// MarkAllAsRead 全通知を既読にする
func (s *notificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

// CreateNotification 通知を作成（内部用）
func (s *notificationService) CreateNotification(notification *model.Notification) error {
	return s.notificationRepo.Create(notification)
}
