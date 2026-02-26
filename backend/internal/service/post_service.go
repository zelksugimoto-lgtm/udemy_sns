package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// PostService 投稿サービスのインターフェース
type PostService interface {
	CreatePost(userID uuid.UUID, req *request.CreatePostRequest) (*response.PostResponse, error)
	GetPost(postID uuid.UUID) (*response.PostResponse, error)
	UpdatePost(userID uuid.UUID, postID uuid.UUID, req *request.UpdatePostRequest) (*response.PostResponse, error)
	DeletePost(userID uuid.UUID, postID uuid.UUID) error
	GetTimeline(userID uuid.UUID, limit, offset int) (*response.PostListResponse, error)
	GetUserPosts(username string, limit, offset int) (*response.PostListResponse, error)
}

type postService struct {
	postRepo     repository.PostRepository
	userRepo     repository.UserRepository
	followRepo   repository.FollowRepository
	likeRepo     repository.LikeRepository
	bookmarkRepo repository.BookmarkRepository
}

// NewPostService 投稿サービスのコンストラクタ
func NewPostService(
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	likeRepo repository.LikeRepository,
	bookmarkRepo repository.BookmarkRepository,
) PostService {
	return &postService{
		postRepo:     postRepo,
		userRepo:     userRepo,
		followRepo:   followRepo,
		likeRepo:     likeRepo,
		bookmarkRepo: bookmarkRepo,
	}
}

// CreatePost 投稿作成
func (s *postService) CreatePost(userID uuid.UUID, req *request.CreatePostRequest) (*response.PostResponse, error) {
	post := &model.Post{
		UserID:     userID,
		Content:    req.Content,
		Visibility: req.Visibility,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// ユーザー情報を取得
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return mapPostToPostResponse(post, user, 0, 0, false, false), nil
}

// GetPost 投稿取得
func (s *postService) GetPost(postID uuid.UUID) (*response.PostResponse, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("投稿が見つかりません")
	}

	// いいね数とコメント数を取得
	likesCount, _ := s.postRepo.CountLikes(postID)
	commentsCount, _ := s.postRepo.CountComments(postID)

	return mapPostToPostResponse(post, &post.User, int(likesCount), int(commentsCount), false, false), nil
}

// UpdatePost 投稿更新
func (s *postService) UpdatePost(userID uuid.UUID, postID uuid.UUID, req *request.UpdatePostRequest) (*response.PostResponse, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("投稿が見つかりません")
	}

	// 権限チェック
	if post.UserID != userID {
		return nil, errors.New("この投稿を更新する権限がありません")
	}

	// 更新
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Visibility != nil {
		post.Visibility = *req.Visibility
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	// いいね数とコメント数を取得
	likesCount, _ := s.postRepo.CountLikes(postID)
	commentsCount, _ := s.postRepo.CountComments(postID)

	return mapPostToPostResponse(post, &post.User, int(likesCount), int(commentsCount), false, false), nil
}

// DeletePost 投稿削除
func (s *postService) DeletePost(userID uuid.UUID, postID uuid.UUID) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("投稿が見つかりません")
	}

	// 権限チェック
	if post.UserID != userID {
		return errors.New("この投稿を削除する権限がありません")
	}

	return s.postRepo.Delete(postID)
}

// GetTimeline タイムライン取得
// Phase 1: 全ユーザーの投稿を表示（フォロー機能は Phase 2 で実装）
func (s *postService) GetTimeline(userID uuid.UUID, limit, offset int) (*response.PostListResponse, error) {
	// Phase 1: 全ユーザーの投稿を取得（公開投稿のみ）
	// フォロー機能はPhase 2で有効化
	var followingIDs []uuid.UUID // 空のスライス（全ユーザーを対象）

	// タイムライン取得
	posts, total, err := s.postRepo.GetTimeline(userID, followingIDs, limit, offset)
	if err != nil {
		return nil, err
	}

	// PostResponseに変換
	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		likesCount, _ := s.postRepo.CountLikes(post.ID)
		commentsCount, _ := s.postRepo.CountComments(post.ID)
		isLiked, _ := s.likeRepo.Exists(userID, "Post", post.ID)
		isBookmarked, _ := s.bookmarkRepo.Exists(userID, post.ID)

		postResponses[i] = *mapPostToPostResponse(&post, &post.User, int(likesCount), int(commentsCount), isLiked, isBookmarked)
	}

	return &response.PostListResponse{
		Posts: postResponses,
		Pagination: response.PaginationResponse{
			Total:  int(total),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

// GetUserPosts ユーザーの投稿一覧取得
func (s *postService) GetUserPosts(username string, limit, offset int) (*response.PostListResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	posts, total, err := s.postRepo.FindByUserID(user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	// PostResponseに変換
	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		likesCount, _ := s.postRepo.CountLikes(post.ID)
		commentsCount, _ := s.postRepo.CountComments(post.ID)

		postResponses[i] = *mapPostToPostResponse(&post, &post.User, int(likesCount), int(commentsCount), false, false)
	}

	return &response.PostListResponse{
		Posts: postResponses,
		Pagination: response.PaginationResponse{
			Total:  int(total),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

// mapPostToPostResponse PostをPostResponseにマッピング
func mapPostToPostResponse(post *model.Post, user *model.User, likesCount, commentsCount int, isLiked, isBookmarked bool) *response.PostResponse {
	return &response.PostResponse{
		ID:            post.ID,
		Content:       post.Content,
		Visibility:    post.Visibility,
		LikesCount:    likesCount,
		CommentsCount: commentsCount,
		IsLiked:       isLiked,
		IsBookmarked:  isBookmarked,
		User: response.UserSimple{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			AvatarURL:   user.AvatarURL,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
