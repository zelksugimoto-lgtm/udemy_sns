package model

import "github.com/google/uuid"

// Follow フォローモデル
type Follow struct {
	BaseModel
	FollowerID uuid.UUID `gorm:"type:uuid;not null;index:idx_follows_follower_followed,unique" json:"follower_id"` // フォローする側
	FollowedID uuid.UUID `gorm:"type:uuid;not null;index:idx_follows_follower_followed,unique" json:"followed_id"` // フォローされる側

	// Relations
	Follower User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Followed User `gorm:"foreignKey:FollowedID" json:"followed,omitempty"`
}

// TableName テーブル名を指定
func (Follow) TableName() string {
	return "follows"
}
