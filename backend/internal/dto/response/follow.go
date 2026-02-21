package response

// FollowListResponse はフォロー/フォロワー一覧レスポンス
type FollowListResponse struct {
	Data       []UserSimple       `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
