package model

import "github.com/google/uuid"

// Bookmark ブックマークモデル
type Bookmark struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_bookmarks_user_post,unique" json:"user_id"`
	PostID uuid.UUID `gorm:"type:uuid;not null;index:idx_bookmarks_user_post,unique" json:"post_id"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName テーブル名を指定
func (Bookmark) TableName() string {
	return "bookmarks"
}
