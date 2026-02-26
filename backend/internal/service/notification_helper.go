package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
)

// CreateLikeNotification いいね通知を作成
func (s *notificationService) CreateLikeNotification(actorID, targetUserID uuid.UUID, targetType string, targetID uuid.UUID) error {
	// 自分自身へのアクションは通知しない
	if actorID == targetUserID {
		return nil
	}

	// 既存の通知があれば削除（いいね取り消し→再いいねのケース）
	s.DeleteNotificationByAction(actorID, targetUserID, "like", targetType, targetID)

	var message string
	var postID *uuid.UUID
	if targetType == "Post" {
		message = "があなたの投稿にいいねしました"
		postID = &targetID
	} else if targetType == "Comment" {
		message = "があなたのコメントにいいねしました"
		// コメントのいいねの場合、targetIDはコメントIDなので、postIDは設定しない（TODO: 必要に応じてコメントから投稿IDを取得）
	}

	notification := &model.Notification{
		UserID:     targetUserID,
		ActorID:    &actorID,
		Type:       "like",
		TargetType: &targetType,
		TargetID:   &targetID,
		PostID:     postID,
		Message:    message,
		IsRead:     false,
	}

	return s.notificationRepo.Create(notification)
}

// CreateCommentNotification コメント通知を作成
func (s *notificationService) CreateCommentNotification(actorID, postOwnerID uuid.UUID, postID uuid.UUID) error {
	// 自分自身へのアクションは通知しない
	if actorID == postOwnerID {
		return nil
	}

	targetType := "Post"
	message := "があなたの投稿にコメントしました"

	notification := &model.Notification{
		UserID:     postOwnerID,
		ActorID:    &actorID,
		Type:       "comment",
		TargetType: &targetType,
		TargetID:   &postID,
		PostID:     &postID,
		Message:    message,
		IsRead:     false,
	}

	return s.notificationRepo.Create(notification)
}

// CreateReplyNotification 返信通知を作成
func (s *notificationService) CreateReplyNotification(actorID, parentCommentOwnerID uuid.UUID, parentCommentID uuid.UUID, postID uuid.UUID) error {
	// 自分自身へのアクションは通知しない
	if actorID == parentCommentOwnerID {
		return nil
	}

	targetType := "Comment"
	message := "があなたのコメントに返信しました"

	notification := &model.Notification{
		UserID:     parentCommentOwnerID,
		ActorID:    &actorID,
		Type:       "reply",
		TargetType: &targetType,
		TargetID:   &parentCommentID,
		PostID:     &postID,
		Message:    message,
		IsRead:     false,
	}

	return s.notificationRepo.Create(notification)
}

// CreateFollowNotification フォロー通知を作成
func (s *notificationService) CreateFollowNotification(actorID, targetUserID uuid.UUID) error {
	// 自分自身へのアクションは通知しない
	if actorID == targetUserID {
		return nil
	}

	// 既存の通知があれば削除（フォロー解除→再フォローのケース）
	s.DeleteNotificationByAction(actorID, targetUserID, "follow", "", uuid.Nil)

	message := "があなたをフォローしました"

	notification := &model.Notification{
		UserID:  targetUserID,
		ActorID: &actorID,
		Type:    "follow",
		Message: message,
		IsRead:  false,
	}

	return s.notificationRepo.Create(notification)
}

// DeleteNotificationByAction アクション取り消し時に通知を削除
func (s *notificationService) DeleteNotificationByAction(actorID, userID uuid.UUID, notifType, targetType string, targetID uuid.UUID) error {
	return s.notificationRepo.DeleteByAction(actorID, userID, notifType, targetType, targetID)
}

// DeleteNotificationsByTarget 対象削除時に関連通知を削除
func (s *notificationService) DeleteNotificationsByTarget(targetType string, targetID uuid.UUID) error {
	return s.notificationRepo.DeleteByTarget(targetType, targetID)
}

// GetUnreadCount 未読通知数を取得
func (s *notificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.notificationRepo.CountUnread(userID)
}

// FormatNotificationMessage 通知メッセージをフォーマット（将来的にアクター名を含める場合）
func FormatNotificationMessage(actorName, baseMessage string) string {
	return fmt.Sprintf("%s%s", actorName, baseMessage)
}
