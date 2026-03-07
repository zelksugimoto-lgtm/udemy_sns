package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
	Search(query string, limit, offset int) ([]model.User, int64, error)
	ListByStatus(status string, limit, offset int) ([]model.User, int64, error)
	CountByStatus() (map[string]int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository ユーザーリポジトリのコンストラクタ
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create ユーザーを作成
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID IDでユーザーを取得
func (r *userRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail メールアドレスでユーザーを取得
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername ユーザー名でユーザーを取得
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update ユーザーを更新
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete ユーザーを削除（論理削除）
func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, id).Error
}

// Search ユーザーを検索
func (r *userRepository) Search(query string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	searchQuery := r.db.Model(&model.User{})

	if query != "" {
		// SQLワイルドカード文字をエスケープ
		escapedQuery := escapeWildcards(query)
		searchPattern := "%" + escapedQuery + "%"
		searchQuery = searchQuery.Where(
			"username ILIKE ? OR display_name ILIKE ?",
			searchPattern, searchPattern,
		)
	}

	// 総数取得
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := searchQuery.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByStatus ステータスでユーザー一覧を取得
func (r *userRepository) ListByStatus(status string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CountByStatus ステータスごとのユーザー数を取得
func (r *userRepository) CountByStatus() (map[string]int64, error) {
	counts := make(map[string]int64)
	statuses := []string{"pending", "approved", "rejected"}

	for _, status := range statuses {
		var count int64
		if err := r.db.Model(&model.User{}).Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		counts[status] = count
	}

	return counts, nil
}

// escapeWildcards SQLワイルドカード文字（%, _）をエスケープ
func escapeWildcards(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\") // バックスラッシュを最初にエスケープ
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
