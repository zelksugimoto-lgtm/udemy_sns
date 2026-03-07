package model

// Admin 管理者モデル
type Admin struct {
	BaseModel
	Username     string `gorm:"type:varchar(50);not null;uniqueIndex:idx_admins_username,where:deleted_at IS NULL" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Email        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_admins_email,where:deleted_at IS NULL" json:"email"`
	Role         string `gorm:"type:varchar(20);not null;default:'admin'" json:"role"` // 'root' or 'admin'
}

// TableName テーブル名を指定
func (Admin) TableName() string {
	return "admins"
}

// IsRoot ルート管理者かどうか
func (a *Admin) IsRoot() bool {
	return a.Role == "root"
}
