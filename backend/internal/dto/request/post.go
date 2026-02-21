package request

// CreatePostRequest は投稿作成リクエスト
type CreatePostRequest struct {
	Content string `json:"content" validate:"required,min=1,max=250" example:"Hello, world!"`
}

// UpdatePostRequest は投稿更新リクエスト
type UpdatePostRequest struct {
	Content string `json:"content" validate:"required,min=1,max=250" example:"Updated content"`
}
