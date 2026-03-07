package errors

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail はエラーの詳細情報
type ErrorDetail struct {
	Code      string            `json:"code" example:"BAD_REQUEST"`
	Message   string            `json:"message" example:"Invalid request parameters"`
	RequestID string            `json:"request_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Details   []ValidationError `json:"details,omitempty"`
}

// ValidationError はバリデーションエラーの詳細
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Invalid email format"`
}

// BadRequest は400エラーを返す
func BadRequest(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	}
}

// Unauthorized は401エラーを返す
func Unauthorized(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	}
}

// Forbidden は403エラーを返す
func Forbidden(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "FORBIDDEN",
			Message: message,
		},
	}
}

// NotFound は404エラーを返す
func NotFound(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "NOT_FOUND",
			Message: message,
		},
	}
}

// Conflict は409エラーを返す
func Conflict(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "CONFLICT",
			Message: message,
		},
	}
}

// InternalError は500エラーを返す
func InternalError(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
		},
	}
}

// ValidationErrors はバリデーションエラーを返す
func ValidationErrors(errs []ValidationError) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: errs,
		},
	}
}

// WithRequestID エラーレスポンスにリクエストIDを追加
func (e *ErrorResponse) WithRequestID(requestID string) *ErrorResponse {
	e.Error.RequestID = requestID
	return e
}

// InternalErrorWithRequestID は500エラーをリクエストIDと共に返す
// ユーザーには詳細を見せず、リクエストIDのみ提示
func InternalErrorWithRequestID(requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:      "INTERNAL_SERVER_ERROR",
			Message:   "サーバー内部エラーが発生しました。お手数ですが、リクエストID（" + requestID + "）をお伝えください。",
			RequestID: requestID,
		},
	}
}
