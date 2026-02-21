package response

import "time"

// PostResponse は投稿レスポンス
type PostResponse struct {
	ID           string       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID       string       `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content      string       `json:"content" example:"Hello, world!"`
	LikeCount    int          `json:"like_count" example:"10"`
	CommentCount int          `json:"comment_count" example:"5"`
	IsLiked      bool         `json:"is_liked" example:"false"`
	IsBookmarked bool         `json:"is_bookmarked" example:"false"`
	User         UserSimple   `json:"user"`
	CreatedAt    time.Time    `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt    time.Time    `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}

// PostListResponse は投稿一覧レスポンス
type PostListResponse struct {
	Data       []PostResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
