package model

import "github.com/google/uuid"

// PostMedia 投稿メディアモデル（画像・動画）
type PostMedia struct {
	BaseModel
	PostID       uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
	MediaType    string    `gorm:"type:varchar(20);not null" json:"media_type"`    // image, video
	MediaURL     string    `gorm:"type:varchar(1000);not null" json:"media_url"`   // Firebase Storage URL（署名付きURL対応）
	ThumbnailURL string    `gorm:"type:varchar(1000)" json:"thumbnail_url"`        // サムネイルURL（Phase 4）
	DisplayOrder int       `gorm:"not null;default:0" json:"display_order"`        // 表示順序（0-3）

	// Relations
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName テーブル名を指定
func (PostMedia) TableName() string {
	return "post_media"
}
