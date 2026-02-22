package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// BookmarkService ブックマークサービスのインターフェース
type BookmarkService interface {
	AddBookmark(userID uuid.UUID, postID uuid.UUID) error
	RemoveBookmark(userID uuid.UUID, postID uuid.UUID) error
	GetBookmarks(userID uuid.UUID, limit, offset int) (*response.BookmarkListResponse, error)
}

type bookmarkService struct {
	bookmarkRepo repository.BookmarkRepository
	postRepo     repository.PostRepository
}

// NewBookmarkService ブックマークサービスのコンストラクタ
func NewBookmarkService(
	bookmarkRepo repository.BookmarkRepository,
	postRepo repository.PostRepository,
) BookmarkService {
	return &bookmarkService{
		bookmarkRepo: bookmarkRepo,
		postRepo:     postRepo,
	}
}

// AddBookmark ブックマーク追加
func (s *bookmarkService) AddBookmark(userID uuid.UUID, postID uuid.UUID) error {
	// 投稿が存在するか確認
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("投稿が見つかりません")
	}

	// すでにブックマークしているか確認
	exists, err := s.bookmarkRepo.Exists(userID, postID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("すでにブックマークしています")
	}

	bookmark := &model.Bookmark{
		UserID: userID,
		PostID: postID,
	}

	return s.bookmarkRepo.Create(bookmark)
}

// RemoveBookmark ブックマーク削除
func (s *bookmarkService) RemoveBookmark(userID uuid.UUID, postID uuid.UUID) error {
	// ブックマークしているか確認
	exists, err := s.bookmarkRepo.Exists(userID, postID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("ブックマークしていません")
	}

	return s.bookmarkRepo.Delete(userID, postID)
}

// GetBookmarks ブックマーク一覧取得
func (s *bookmarkService) GetBookmarks(userID uuid.UUID, limit, offset int) (*response.BookmarkListResponse, error) {
	bookmarks, total, err := s.bookmarkRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// PostResponseに変換
	posts := make([]response.PostResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		likesCount, _ := s.postRepo.CountLikes(bookmark.PostID)
		commentsCount, _ := s.postRepo.CountComments(bookmark.PostID)

		posts[i] = *mapPostToPostResponse(&bookmark.Post, &bookmark.Post.User, int(likesCount), int(commentsCount), false)
	}

	return &response.BookmarkListResponse{
		Posts: posts,
		Pagination: response.PaginationResponse{
			Total:  int(total),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}
