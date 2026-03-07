package repository

import (
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// AdminRepository 管理者リポジトリのインターフェース
type AdminRepository interface {
	Create(admin *model.Admin) error
	FindByID(id string) (*model.Admin, error)
	FindByUsername(username string) (*model.Admin, error)
	FindByEmail(email string) (*model.Admin, error)
	List(limit, offset int) ([]model.Admin, int64, error)
	Update(admin *model.Admin) error
	Delete(id string) error
}

type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 管理者リポジトリを生成
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// Create 管理者を作成
func (r *adminRepository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

// FindByID IDで管理者を取得
func (r *adminRepository) FindByID(id string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByUsername ユーザー名で管理者を取得
func (r *adminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByEmail メールアドレスで管理者を取得
func (r *adminRepository) FindByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// List 管理者一覧を取得
func (r *adminRepository) List(limit, offset int) ([]model.Admin, int64, error) {
	var admins []model.Admin
	var total int64

	if err := r.db.Model(&model.Admin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// Update 管理者を更新
func (r *adminRepository) Update(admin *model.Admin) error {
	return r.db.Save(admin).Error
}

// Delete 管理者を削除（論理削除）
func (r *adminRepository) Delete(id string) error {
	return r.db.Delete(&model.Admin{}, "id = ?", id).Error
}
