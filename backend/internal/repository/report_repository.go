package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/gorm"
)

// ReportRepository 通報リポジトリのインターフェース
type ReportRepository interface {
	Create(report *model.Report) error
	FindByID(id uuid.UUID) (*model.Report, error)
	FindAll(limit, offset int) ([]model.Report, int64, error)
	FindByStatus(status string, limit, offset int) ([]model.Report, int64, error)
	Update(report *model.Report) error
}

type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository 通報リポジトリのコンストラクタ
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// Create 通報を作成
func (r *reportRepository) Create(report *model.Report) error {
	return r.db.Create(report).Error
}

// FindByID IDで通報を取得
func (r *reportRepository) FindByID(id uuid.UUID) (*model.Report, error) {
	var report model.Report
	err := r.db.Preload("Reporter").Preload("Reviewer").Where("id = ?", id).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}

// FindAll 全通報を取得
func (r *reportRepository) FindAll(limit, offset int) ([]model.Report, int64, error) {
	var reports []model.Report
	var total int64

	query := r.db.Model(&model.Report{})

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("Reporter").
		Preload("Reviewer").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// FindByStatus ステータスで通報を取得
func (r *reportRepository) FindByStatus(status string, limit, offset int) ([]model.Report, int64, error) {
	var reports []model.Report
	var total int64

	query := r.db.Model(&model.Report{}).Where("status = ?", status)

	// 総数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	err := query.Preload("Reporter").
		Preload("Reviewer").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// Update 通報を更新
func (r *reportRepository) Update(report *model.Report) error {
	return r.db.Save(report).Error
}
