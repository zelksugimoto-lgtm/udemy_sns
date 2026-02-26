package response

import (
	"time"

	"github.com/google/uuid"
)

// AuthResponse は認証レスポンス（登録・ログイン共通）
type AuthResponse struct {
	Token string        `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  *UserResponse `json:"user"`
}

// UserResponse はユーザー情報レスポンス
type UserResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email       string    `json:"email" example:"user@example.com"`
	Username    string    `json:"username" example:"johndoe"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	Bio         string    `json:"bio" example:"Hello, I'm John!"`
	AvatarURL   string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	HeaderURL   string    `json:"header_url,omitempty" example:"https://example.com/header.jpg"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-01T00:00:00Z"`
}

// UserProfileResponse はプロフィールレスポンス
type UserProfileResponse struct {
	ID             uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username       string    `json:"username" example:"johndoe"`
	DisplayName    string    `json:"display_name" example:"John Doe"`
	Bio            string    `json:"bio" example:"Hello, I'm John!"`
	AvatarURL      string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	HeaderURL      string    `json:"header_url,omitempty" example:"https://example.com/header.jpg"`
	FollowersCount int       `json:"followers_count" example:"100"`
	FollowingCount int       `json:"following_count" example:"50"`
	PostsCount     int       `json:"posts_count" example:"25"`
	IsFollowing    *bool     `json:"is_following" example:"false"`
	CreatedAt      time.Time `json:"created_at" example:"2025-01-01T00:00:00Z"`
}

// UserListResponse はユーザー一覧レスポンス
type UserListResponse struct {
	Users      []UserSimple       `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}

// UserSimple は簡易ユーザー情報
type UserSimple struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username    string    `json:"username" example:"johndoe"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	Bio         string    `json:"bio,omitempty" example:"Hello, I'm John!"`
	AvatarURL   string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	IsFollowing *bool     `json:"is_following" example:"false"`
}

// PaginationResponse はページネーション情報
type PaginationResponse struct {
	Total   int  `json:"total" example:"100"`
	Limit   int  `json:"limit" example:"20"`
	Offset  int  `json:"offset" example:"0"`
	HasMore bool `json:"has_more" example:"true"`
}
