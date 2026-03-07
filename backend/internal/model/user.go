package model

// User ユーザーモデル
type User struct {
	BaseModel
	Email        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email,where:deleted_at IS NULL" json:"email"`
	Username     string `gorm:"type:varchar(50);not null;uniqueIndex:idx_users_username,where:deleted_at IS NULL" json:"username"`
	DisplayName  string `gorm:"type:varchar(100);not null" json:"display_name"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Bio          string `gorm:"type:text" json:"bio"`
	AvatarURL    string `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	HeaderURL    string `gorm:"type:varchar(500)" json:"header_url,omitempty"`
	IsAdmin      bool   `gorm:"default:false" json:"is_admin"`
	Status       string `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // 'pending', 'approved', 'rejected'

	// Relations
	Posts         []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Likes         []Like         `gorm:"foreignKey:UserID" json:"likes,omitempty"`
	Bookmarks     []Bookmark     `gorm:"foreignKey:UserID" json:"bookmarks,omitempty"`
	Following     []Follow       `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	Followers     []Follow       `gorm:"foreignKey:FollowedID" json:"followers,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	Reports       []Report       `gorm:"foreignKey:ReporterID" json:"reports,omitempty"`
}

// IsApproved ユーザーが承認済みかどうか
func (u *User) IsApproved() bool {
	return u.Status == "approved"
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}
