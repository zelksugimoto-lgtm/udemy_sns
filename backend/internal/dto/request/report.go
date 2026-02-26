package request

import "github.com/google/uuid"

// CreateReportRequest は通報作成リクエスト
type CreateReportRequest struct {
	TargetType string    `json:"target_type" validate:"required,oneof=Post Comment User" example:"Post"`
	TargetID   uuid.UUID `json:"target_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Reason     string    `json:"reason" validate:"required,oneof=spam inappropriate_content harassment other" example:"spam"`
	Comment    *string   `json:"comment,omitempty" validate:"max=1000" example:"This is spam content"`
}
