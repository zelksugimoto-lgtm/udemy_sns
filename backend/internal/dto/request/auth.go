package request

// RegisterRequest はユーザー登録リクエスト
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string `json:"password" validate:"required,min=8" example:"password123"`
	Username    string `json:"username" validate:"required,min=3,max=50,alphanum" example:"johndoe"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" example:"John Doe"`
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}
