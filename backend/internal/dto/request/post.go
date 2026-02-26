package request

// CreatePostRequest は投稿作成リクエスト
type CreatePostRequest struct {
	Content    string `json:"content" validate:"required,min=1,max=250" example:"Hello, world!"`
	Visibility string `json:"visibility" validate:"required,oneof=public followers private" example:"public"`
}

// UpdatePostRequest は投稿更新リクエスト
type UpdatePostRequest struct {
	Content    *string `json:"content,omitempty" validate:"omitempty,min=1,max=250" example:"Updated content"`
	Visibility *string `json:"visibility,omitempty" validate:"omitempty,oneof=public followers private" example:"public"`
}
