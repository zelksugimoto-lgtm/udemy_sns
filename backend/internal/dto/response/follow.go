package response

// FollowListResponse はフォロー/フォロワー一覧レスポンス
type FollowListResponse struct {
	Users      []UserSimple       `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}
