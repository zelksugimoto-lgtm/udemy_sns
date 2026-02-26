package response

import (
	"time"

	"github.com/google/uuid"
)

// CommentResponse はコメントレスポンス
type CommentResponse struct {
	ID            uuid.UUID         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PostID        uuid.UUID         `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content       string            `json:"content" example:"Great post!"`
	LikesCount    int               `json:"likes_count" example:"5"`
	IsLiked       bool              `json:"is_liked" example:"false"`
	User          UserSimple        `json:"user"`
	ChildComments []CommentResponse `json:"child_comments,omitempty"`
	CreatedAt     time.Time         `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt     time.Time         `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}

// CommentListResponse はコメント一覧レスポンス
type CommentListResponse struct {
	Comments   []CommentResponse  `json:"comments"`
	Pagination PaginationResponse `json:"pagination"`
}
