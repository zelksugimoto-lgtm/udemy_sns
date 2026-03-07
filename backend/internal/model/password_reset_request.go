package model

import "time"

// PasswordResetRequest パスワードリセット申請モデル
type PasswordResetRequest struct {
	BaseModel
	UserID         string     `gorm:"type:uuid;not null;index" json:"user_id"`
	Reason         string     `gorm:"type:text;not null" json:"reason"`
	Status         string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // 'pending', 'approved', 'rejected'
	ResetToken     *string    `gorm:"type:varchar(255);uniqueIndex:idx_reset_token,where:deleted_at IS NULL" json:"reset_token,omitempty"`
	TokenExpiresAt *time.Time `gorm:"type:timestamp" json:"token_expires_at,omitempty"`
	RequestedAt    time.Time  `gorm:"type:timestamp;not null" json:"requested_at"`
	ProcessedAt    *time.Time `gorm:"type:timestamp" json:"processed_at,omitempty"`
	ProcessedBy    *string    `gorm:"type:uuid" json:"processed_by,omitempty"`

	// Relations
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProcessedByAdmin *Admin `gorm:"foreignKey:ProcessedBy" json:"processed_by_admin,omitempty"`
}

// TableName テーブル名を指定
func (PasswordResetRequest) TableName() string {
	return "password_reset_requests"
}

// IsTokenValid トークンが有効かどうか
func (p *PasswordResetRequest) IsTokenValid() bool {
	if p.ResetToken == nil || p.TokenExpiresAt == nil {
		return false
	}
	if p.Status != "approved" {
		return false
	}
	return time.Now().Before(*p.TokenExpiresAt)
}
