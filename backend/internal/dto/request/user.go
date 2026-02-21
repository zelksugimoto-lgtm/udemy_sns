package request

// UpdateProfileRequest はプロフィール更新リクエスト
type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,min=1,max=100" example:"John Doe Updated"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=1000" example:"This is my bio"`
}
