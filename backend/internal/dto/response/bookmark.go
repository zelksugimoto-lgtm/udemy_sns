package response

// BookmarkListResponse はブックマーク一覧レスポンス
type BookmarkListResponse struct {
	Data       []PostResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
