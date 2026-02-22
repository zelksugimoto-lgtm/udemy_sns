package model

import "github.com/google/uuid"

// Comment コメントモデル
type Comment struct {
	BaseModel
	PostID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"post_id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	ParentCommentID *uuid.UUID `gorm:"type:uuid;index" json:"parent_comment_id,omitempty"`
	Content         string     `gorm:"type:text;not null" json:"content"`

	// Relations
	Post          Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentComment *Comment  `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	ChildComments []Comment `gorm:"foreignKey:ParentCommentID" json:"child_comments,omitempty"`
	Likes         []Like    `gorm:"foreignKey:LikeableID;polymorphicType:LikeableType;polymorphicValue:Comment" json:"likes,omitempty"`
}

// TableName テーブル名を指定
func (Comment) TableName() string {
	return "comments"
}
