package model

import "github.com/google/uuid"

// Notification 通知モデル
type Notification struct {
	BaseModel
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`                                        // 通知を受け取るユーザー
	ActorID    *uuid.UUID `gorm:"type:uuid;index" json:"actor_id,omitempty"`                                      // 通知のアクションを起こしたユーザー
	Type       string     `gorm:"type:varchar(50);not null;index" json:"type"`                                    // like, comment, follow, mention
	TargetType *string    `gorm:"type:varchar(50);index:idx_notifications_target" json:"target_type,omitempty"`  // Post, Comment
	TargetID   *uuid.UUID `gorm:"type:uuid;index:idx_notifications_target" json:"target_id,omitempty"`           // 対象のID
	Message    string     `gorm:"type:text;not null" json:"message"`                                              // 通知メッセージ
	IsRead     bool       `gorm:"not null;default:false;index" json:"is_read"`                                    // 既読フラグ

	// Relations
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Actor *User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

// TableName テーブル名を指定
func (Notification) TableName() string {
	return "notifications"
}
