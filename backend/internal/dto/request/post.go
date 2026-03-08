package request

// PostMediaRequest は投稿メディアリクエスト
type PostMediaRequest struct {
	MediaType    string `json:"media_type" validate:"required,oneof=image video" example:"image"`
	MediaURL     string `json:"media_url" validate:"required,url" example:"https://firebasestorage.googleapis.com/..."`
	DisplayOrder int    `json:"display_order" validate:"gte=0,lte=3" example:"0"`
}

// CreatePostRequest は投稿作成リクエスト
type CreatePostRequest struct {
	Content    string              `json:"content" validate:"required,min=1,max=250" example:"Hello, world!"`
	Visibility string              `json:"visibility" validate:"required,oneof=public followers private" example:"public"`
	Media      *[]PostMediaRequest `json:"media,omitempty" validate:"omitempty,max=4,dive"`
}

// UpdatePostRequest は投稿更新リクエスト
type UpdatePostRequest struct {
	Content    *string `json:"content,omitempty" validate:"omitempty,min=1,max=250" example:"Updated content"`
	Visibility *string `json:"visibility,omitempty" validate:"omitempty,oneof=public followers private" example:"public"`
}
