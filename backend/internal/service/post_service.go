package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// PostService 投稿サービスのインターフェース
type PostService interface {
	CreatePost(userID uuid.UUID, req *request.CreatePostRequest) (*response.PostResponse, error)
	GetPost(postID uuid.UUID, currentUserID *uuid.UUID) (*response.PostResponse, error)
	UpdatePost(userID uuid.UUID, postID uuid.UUID, req *request.UpdatePostRequest) (*response.PostResponse, error)
	DeletePost(userID uuid.UUID, postID uuid.UUID) error
	GetTimeline(userID uuid.UUID, limit, offset int) (*response.PostListResponse, error)
	GetUserPosts(username string, limit, offset int, currentUserID *uuid.UUID) (*response.PostListResponse, error)
	GetUserLikedPosts(username string, limit, offset int, currentUserID *uuid.UUID) (*response.PostListResponse, error)
}

type postService struct {
	postRepo      repository.PostRepository
	postMediaRepo repository.PostMediaRepository
	userRepo      repository.UserRepository
	followRepo    repository.FollowRepository
	likeRepo      repository.LikeRepository
	bookmarkRepo  repository.BookmarkRepository
	storageService *StorageService
}

// NewPostService 投稿サービスのコンストラクタ
func NewPostService(
	postRepo repository.PostRepository,
	postMediaRepo repository.PostMediaRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	likeRepo repository.LikeRepository,
	bookmarkRepo repository.BookmarkRepository,
	storageService *StorageService,
) PostService {
	return &postService{
		postRepo:       postRepo,
		postMediaRepo:  postMediaRepo,
		userRepo:       userRepo,
		followRepo:     followRepo,
		likeRepo:       likeRepo,
		bookmarkRepo:   bookmarkRepo,
		storageService: storageService,
	}
}

// CreatePost 投稿作成
func (s *postService) CreatePost(userID uuid.UUID, req *request.CreatePostRequest) (*response.PostResponse, error) {
	// メディアのバリデーション
	if req.Media != nil && len(*req.Media) > 0 {
		if len(*req.Media) > 4 {
			return nil, errors.New("メディアは最大4つまで添付できます")
		}

		for _, mediaReq := range *req.Media {
			// メディアタイプ検証
			if err := s.storageService.ValidateMediaType(mediaReq.MediaType); err != nil {
				return nil, err
			}

			// URLバリデーション
			if err := s.storageService.ValidateMediaURL(mediaReq.MediaURL); err != nil {
				return nil, err
			}

			// セキュリティ: メディアの所有権を検証（IDOR攻撃対策）
			if err := s.storageService.ValidateMediaOwnership(mediaReq.MediaURL, userID); err != nil {
				return nil, err
			}

			// display_orderバリデーション
			if mediaReq.DisplayOrder < 0 || mediaReq.DisplayOrder > 3 {
				return nil, errors.New("display_orderは0-3の範囲で指定してください")
			}
		}
	}

	// 投稿作成
	post := &model.Post{
		UserID:     userID,
		Content:    req.Content,
		Visibility: req.Visibility,
	}

	if err := s.postRepo.Create(post); err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("投稿作成失敗（DB）")
		return nil, err
	}

	// メディア保存
	var mediaModels []model.PostMedia
	if req.Media != nil && len(*req.Media) > 0 {
		mediaModels = make([]model.PostMedia, len(*req.Media))
		for i, mediaReq := range *req.Media {
			mediaModels[i] = model.PostMedia{
				PostID:       post.ID,
				MediaType:    mediaReq.MediaType,
				MediaURL:     mediaReq.MediaURL,
				DisplayOrder: mediaReq.DisplayOrder,
			}
		}

		if err := s.postMediaRepo.CreateBatch(mediaModels); err != nil {
			log.Error().Err(err).Str("post_id", post.ID.String()).Int("media_count", len(mediaModels)).Msg("メディア保存失敗（DB）")
			return nil, err
		}
	}

	// ユーザー情報を取得
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("ユーザー情報取得失敗")
		return nil, err
	}

	return mapPostToPostResponse(post, user, 0, 0, false, false, mediaModels), nil
}

// GetPost 投稿取得
func (s *postService) GetPost(postID uuid.UUID, currentUserID *uuid.UUID) (*response.PostResponse, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("投稿が見つかりません")
	}

	// メディア取得
	media, _ := s.postMediaRepo.FindByPostID(postID)

	// いいね数とコメント数を取得
	likesCount, _ := s.postRepo.CountLikes(postID)
	commentsCount, _ := s.postRepo.CountComments(postID)

	// いいね・ブックマーク状態を確認
	var isLiked, isBookmarked bool
	if currentUserID != nil {
		isLiked, _ = s.likeRepo.Exists(*currentUserID, "Post", postID)
		isBookmarked, _ = s.bookmarkRepo.Exists(*currentUserID, postID)
	}

	return mapPostToPostResponse(post, &post.User, int(likesCount), int(commentsCount), isLiked, isBookmarked, media), nil
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

	// メディア取得
	media, _ := s.postMediaRepo.FindByPostID(postID)

	// いいね数とコメント数を取得
	likesCount, _ := s.postRepo.CountLikes(postID)
	commentsCount, _ := s.postRepo.CountComments(postID)

	return mapPostToPostResponse(post, &post.User, int(likesCount), int(commentsCount), false, false, media), nil
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

	// 投稿IDリストを作成
	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	// 一括取得（N+1クエリ解消）
	likesCounts, _ := s.postRepo.CountLikesInBatch(postIDs)
	commentsCounts, _ := s.postRepo.CountCommentsInBatch(postIDs)
	isLikedMap, _ := s.likeRepo.ExistsInBatch(userID, "Post", postIDs)
	isBookmarkedMap, _ := s.bookmarkRepo.ExistsInBatch(userID, postIDs)

	// メディアを一括取得
	allMedia, _ := s.postMediaRepo.FindByPostIDs(postIDs)
	mediaMap := make(map[uuid.UUID][]model.PostMedia)
	for _, media := range allMedia {
		mediaMap[media.PostID] = append(mediaMap[media.PostID], media)
	}

	// PostResponseに変換
	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		likesCount := int(likesCounts[post.ID])
		commentsCount := int(commentsCounts[post.ID])
		isLiked := isLikedMap[post.ID]
		isBookmarked := isBookmarkedMap[post.ID]
		media := mediaMap[post.ID]

		postResponses[i] = *mapPostToPostResponse(&post, &post.User, likesCount, commentsCount, isLiked, isBookmarked, media)
	}

	hasMore := offset+limit < int(total)

	return &response.PostListResponse{
		Posts: postResponses,
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}

// GetUserPosts ユーザーの投稿一覧取得
func (s *postService) GetUserPosts(username string, limit, offset int, currentUserID *uuid.UUID) (*response.PostListResponse, error) {
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

	// 投稿IDリストを作成
	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	// 一括取得（N+1クエリ解消）
	likesCounts, _ := s.postRepo.CountLikesInBatch(postIDs)
	commentsCounts, _ := s.postRepo.CountCommentsInBatch(postIDs)

	// いいね・ブックマーク状態を確認（ループ外で一度だけ実行）
	var isLikedMap, isBookmarkedMap map[uuid.UUID]bool
	if currentUserID != nil {
		isLikedMap, _ = s.likeRepo.ExistsInBatch(*currentUserID, "Post", postIDs)
		isBookmarkedMap, _ = s.bookmarkRepo.ExistsInBatch(*currentUserID, postIDs)
	}

	// メディアを一括取得
	allMedia, _ := s.postMediaRepo.FindByPostIDs(postIDs)
	mediaMap := make(map[uuid.UUID][]model.PostMedia)
	for _, media := range allMedia {
		mediaMap[media.PostID] = append(mediaMap[media.PostID], media)
	}

	// PostResponseに変換
	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		likesCount := int(likesCounts[post.ID])
		commentsCount := int(commentsCounts[post.ID])
		isLiked := isLikedMap[post.ID]
		isBookmarked := isBookmarkedMap[post.ID]
		media := mediaMap[post.ID]

		postResponses[i] = *mapPostToPostResponse(&post, &post.User, likesCount, commentsCount, isLiked, isBookmarked, media)
	}

	hasMore := offset+limit < int(total)

	return &response.PostListResponse{
		Posts: postResponses,
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}

// GetUserLikedPosts ユーザーがいいねした投稿一覧取得
func (s *postService) GetUserLikedPosts(username string, limit, offset int, currentUserID *uuid.UUID) (*response.PostListResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	posts, total, err := s.postRepo.FindLikedByUser(user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	// 投稿IDリストを作成
	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	// 一括取得（N+1クエリ解消）
	likesCounts, _ := s.postRepo.CountLikesInBatch(postIDs)
	commentsCounts, _ := s.postRepo.CountCommentsInBatch(postIDs)

	// いいね・ブックマーク状態を確認（ループ外で一度だけ実行）
	var isLikedMap, isBookmarkedMap map[uuid.UUID]bool
	if currentUserID != nil {
		isLikedMap, _ = s.likeRepo.ExistsInBatch(*currentUserID, "Post", postIDs)
		isBookmarkedMap, _ = s.bookmarkRepo.ExistsInBatch(*currentUserID, postIDs)
	}

	// メディアを一括取得
	allMedia, _ := s.postMediaRepo.FindByPostIDs(postIDs)
	mediaMap := make(map[uuid.UUID][]model.PostMedia)
	for _, media := range allMedia {
		mediaMap[media.PostID] = append(mediaMap[media.PostID], media)
	}

	// PostResponseに変換
	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		likesCount := int(likesCounts[post.ID])
		commentsCount := int(commentsCounts[post.ID])
		isLiked := isLikedMap[post.ID]
		isBookmarked := isBookmarkedMap[post.ID]
		media := mediaMap[post.ID]

		postResponses[i] = *mapPostToPostResponse(&post, &post.User, likesCount, commentsCount, isLiked, isBookmarked, media)
	}

	hasMore := offset+limit < int(total)

	return &response.PostListResponse{
		Posts: postResponses,
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}

// mapPostToPostResponse PostをPostResponseにマッピング
func mapPostToPostResponse(post *model.Post, user *model.User, likesCount, commentsCount int, isLiked, isBookmarked bool, media []model.PostMedia) *response.PostResponse {
	// メディアをレスポンス形式に変換
	mediaResponses := make([]response.PostMediaResponse, len(media))
	for i, m := range media {
		mediaResponses[i] = response.PostMediaResponse{
			ID:           m.ID,
			MediaType:    m.MediaType,
			MediaURL:     m.MediaURL,
			ThumbnailURL: m.ThumbnailURL,
			DisplayOrder: m.DisplayOrder,
		}
	}

	return &response.PostResponse{
		ID:            post.ID,
		Content:       post.Content,
		Visibility:    post.Visibility,
		Media:         mediaResponses,
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
