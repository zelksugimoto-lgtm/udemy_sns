package response

import "time"

// NotificationResponse は通知レスポンス
type NotificationResponse struct {
	ID               string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	NotificationType string     `json:"notification_type" example:"like"`
	Actor            UserSimple `json:"actor"`
	TargetType       *string    `json:"target_type,omitempty" example:"Post"`
	TargetID         *string    `json:"target_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsRead           bool       `json:"is_read" example:"false"`
	CreatedAt        time.Time  `json:"created_at" example:"2025-01-01T00:00:00Z"`
}

// NotificationListResponse は通知一覧レスポンス
type NotificationListResponse struct {
	Data        []NotificationResponse `json:"data"`
	UnreadCount int                    `json:"unread_count" example:"5"`
	Pagination  PaginationResponse     `json:"pagination"`
}
