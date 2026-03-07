package session

import (
	"encoding/gob"
	"os"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "admin_session"
	// セッションキー
	KeyAdminID   = "admin_id"
	KeyUsername  = "username"
	KeyRole      = "role"
	KeyIsRoot    = "is_root"
)

var store *sessions.CookieStore

// Init セッションストアを初期化
func Init() {
	// AdminSessionデータ型を登録（gob encoding用）
	gob.Register(map[string]interface{}{})

	// セッションシークレットキーを環境変数から取得
	secret := os.Getenv("ADMIN_SESSION_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production" // デフォルト値（本番環境では必ず変更）
	}

	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/admin",
		MaxAge:   30 * 60, // 30分
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production", // 本番環境ではtrue
		SameSite: 1,      // Lax
	}
}

// GetStore セッションストアを取得
func GetStore() *sessions.CookieStore {
	if store == nil {
		Init()
	}
	return store
}

// AdminSession 管理者セッションデータ
type AdminSession struct {
	AdminID  string
	Username string
	Role     string
	IsRoot   bool
}

// IsAuthenticated 認証済みかどうか
func (s *AdminSession) IsAuthenticated() bool {
	return s.AdminID != "" || s.IsRoot
}
