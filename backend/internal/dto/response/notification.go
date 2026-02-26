package response

import (
	"time"

	"github.com/google/uuid"
)

// NotificationResponse は通知レスポンス
type NotificationResponse struct {
	ID         uuid.UUID   `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type       string      `json:"type" example:"like"`
	Message    string      `json:"message" example:"@johndoe liked your post"`
	IsRead     bool        `json:"is_read" example:"false"`
	Actor      *UserSimple `json:"actor,omitempty"`
	TargetType *string     `json:"target_type,omitempty" example:"Post"`
	TargetID   *uuid.UUID  `json:"target_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	PostID     *uuid.UUID  `json:"post_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt  time.Time   `json:"created_at" example:"2025-01-01T00:00:00Z"`
}

// NotificationListResponse は通知一覧レスポンス
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int                    `json:"unread_count" example:"5"`
	Pagination    PaginationResponse     `json:"pagination"`
}
