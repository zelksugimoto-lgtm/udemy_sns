package model

import "github.com/google/uuid"

// Report 通報モデル
type Report struct {
	BaseModel
	ReporterID uuid.UUID  `gorm:"type:uuid;not null;index" json:"reporter_id"`                           // 通報者
	TargetType string     `gorm:"type:varchar(50);not null;index:idx_reports_target" json:"target_type"` // Post, Comment, User
	TargetID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_reports_target" json:"target_id"`          // 通報対象のID
	Reason     string     `gorm:"type:varchar(100);not null" json:"reason"`                               // spam, harassment, inappropriate_content, other
	Comment    *string    `gorm:"type:text" json:"comment,omitempty"`                                     // 詳細コメント
	Status     string     `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`        // pending, reviewed, resolved, rejected
	ReviewerID *uuid.UUID `gorm:"type:uuid;index" json:"reviewer_id,omitempty"`                           // レビュー担当者（管理者）

	// Relations
	Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// TableName テーブル名を指定
func (Report) TableName() string {
	return "reports"
}
