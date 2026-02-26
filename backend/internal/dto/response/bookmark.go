package response

// BookmarkListResponse はブックマーク一覧レスポンス
type BookmarkListResponse struct {
	Posts      []PostResponse     `json:"posts"`
	Pagination PaginationResponse `json:"pagination"`
}
