package request

// CreateReportRequest は通報作成リクエスト
type CreateReportRequest struct {
	ReportableType string `json:"reportable_type" validate:"required,oneof=Post Comment User" example:"Post"`
	ReportableID   string `json:"reportable_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Reason         string `json:"reason" validate:"required,oneof=spam inappropriate harassment other" example:"spam"`
	Details        string `json:"details,omitempty" validate:"max=1000" example:"This is spam content"`
}
