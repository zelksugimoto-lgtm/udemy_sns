package model

import "github.com/google/uuid"

// Post 投稿モデル
type Post struct {
	BaseModel
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Visibility string    `gorm:"type:varchar(20);not null;default:'public';index" json:"visibility"` // public, followers, private

	// Relations
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments  []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Likes     []Like     `gorm:"foreignKey:LikeableID;polymorphicType:LikeableType;polymorphicValue:Post" json:"likes,omitempty"`
	Bookmarks []Bookmark `gorm:"foreignKey:PostID" json:"bookmarks,omitempty"`
}

// TableName テーブル名を指定
func (Post) TableName() string {
	return "posts"
}
