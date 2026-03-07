package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken リフレッシュトークンモデル
type RefreshToken struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string         `gorm:"type:varchar(500);not null;unique" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IsRevoked bool           `gorm:"default:false;not null" json:"is_revoked"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// リレーション
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName テーブル名を指定
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid トークンが有効かどうかを判定
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && time.Now().Before(rt.ExpiresAt)
}
