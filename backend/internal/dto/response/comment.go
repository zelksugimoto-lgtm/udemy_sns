package response

import "time"

// CommentResponse はコメントレスポンス
type CommentResponse struct {
	ID              string            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PostID          string            `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID          string            `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParentCommentID *string           `json:"parent_comment_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content         string            `json:"content" example:"Great post!"`
	LikeCount       int               `json:"like_count" example:"5"`
	IsLiked         bool              `json:"is_liked" example:"false"`
	User            UserSimple        `json:"user"`
	ChildComments   []CommentResponse `json:"child_comments,omitempty"`
	CreatedAt       time.Time         `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt       time.Time         `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}

// CommentListResponse はコメント一覧レスポンス
type CommentListResponse struct {
	Data       []CommentResponse  `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
