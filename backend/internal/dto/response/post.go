package response

import (
	"time"

	"github.com/google/uuid"
)

// PostMediaResponse は投稿メディアレスポンス
type PostMediaResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	MediaType    string    `json:"media_type" example:"image"`
	MediaURL     string    `json:"media_url" example:"https://firebasestorage.googleapis.com/..."`
	ThumbnailURL string    `json:"thumbnail_url,omitempty" example:"https://firebasestorage.googleapis.com/..."`
	DisplayOrder int       `json:"display_order" example:"0"`
}

// PostResponse は投稿レスポンス
type PostResponse struct {
	ID            uuid.UUID           `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content       string              `json:"content" example:"Hello, world!"`
	Visibility    string              `json:"visibility" example:"public"`
	Media         []PostMediaResponse `json:"media,omitempty"`
	LikesCount    int                 `json:"likes_count" example:"10"`
	CommentsCount int                 `json:"comments_count" example:"5"`
	IsLiked       bool                `json:"is_liked" example:"false"`
	IsBookmarked  bool                `json:"is_bookmarked" example:"false"`
	User          UserSimple          `json:"user"`
	CreatedAt     time.Time           `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt     time.Time           `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}

// PostListResponse は投稿一覧レスポンス
type PostListResponse struct {
	Posts      []PostResponse     `json:"posts"`
	Pagination PaginationResponse `json:"pagination"`
}
