package model

import "github.com/google/uuid"

// Like いいねモデル（ポリモーフィック: Post, Comment）
type Like struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	LikeableType string    `gorm:"type:varchar(50);not null;index:idx_likes_likeable" json:"likeable_type"` // Post, Comment
	LikeableID   uuid.UUID `gorm:"type:uuid;not null;index:idx_likes_likeable" json:"likeable_id"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName テーブル名を指定
func (Like) TableName() string {
	return "likes"
}
