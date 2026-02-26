package response

import (
	"time"

	"github.com/google/uuid"
)

// PostResponse は投稿レスポンス
type PostResponse struct {
	ID            uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content       string     `json:"content" example:"Hello, world!"`
	Visibility    string     `json:"visibility" example:"public"`
	LikesCount    int        `json:"likes_count" example:"10"`
	CommentsCount int        `json:"comments_count" example:"5"`
	IsLiked       bool       `json:"is_liked" example:"false"`
	IsBookmarked  bool       `json:"is_bookmarked" example:"false"`
	User          UserSimple `json:"user"`
	CreatedAt     time.Time  `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}

// PostListResponse は投稿一覧レスポンス
type PostListResponse struct {
	Posts      []PostResponse     `json:"posts"`
	Pagination PaginationResponse `json:"pagination"`
}
