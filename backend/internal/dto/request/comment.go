package request

// CreateCommentRequest はコメント作成リクエスト
type CreateCommentRequest struct {
	Content         string  `json:"content" validate:"required,min=1,max=500" example:"Great post!"`
	ParentCommentID *string `json:"parent_comment_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
}
