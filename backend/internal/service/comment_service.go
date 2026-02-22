package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// CommentService コメントサービスのインターフェース
type CommentService interface {
	CreateComment(userID uuid.UUID, postID uuid.UUID, req *request.CreateCommentRequest) (*response.CommentResponse, error)
	GetComments(postID uuid.UUID, limit, offset int) (*response.CommentListResponse, error)
	DeleteComment(userID uuid.UUID, commentID uuid.UUID) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	likeRepo    repository.LikeRepository
}

// NewCommentService コメントサービスのコンストラクタ
func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	likeRepo repository.LikeRepository,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		likeRepo:    likeRepo,
	}
}

// CreateComment コメント作成
func (s *commentService) CreateComment(userID uuid.UUID, postID uuid.UUID, req *request.CreateCommentRequest) (*response.CommentResponse, error) {
	// 投稿が存在するか確認
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("投稿が見つかりません")
	}

	comment := &model.Comment{
		PostID:          postID,
		UserID:          userID,
		Content:         req.Content,
		ParentCommentID: req.ParentCommentID,
	}

	// 親コメントが指定されている場合、存在確認
	if req.ParentCommentID != nil {
		parentComment, err := s.commentRepo.FindByID(*req.ParentCommentID)
		if err != nil {
			return nil, err
		}
		if parentComment == nil {
			return nil, errors.New("親コメントが見つかりません")
		}
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// ユーザー情報を取得
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return mapCommentToCommentResponse(comment, user, 0, false, nil), nil
}

// GetComments コメント一覧取得
func (s *commentService) GetComments(postID uuid.UUID, limit, offset int) (*response.CommentListResponse, error) {
	// 投稿が存在するか確認
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("投稿が見つかりません")
	}

	comments, total, err := s.commentRepo.FindByPostID(postID, limit, offset)
	if err != nil {
		return nil, err
	}

	// CommentResponseに変換
	commentResponses := make([]response.CommentResponse, len(comments))
	for i, comment := range comments {
		likesCount, _ := s.commentRepo.CountLikes(comment.ID)

		// 子コメントを変換
		var childComments []response.CommentResponse
		if len(comment.ChildComments) > 0 {
			childComments = make([]response.CommentResponse, len(comment.ChildComments))
			for j, child := range comment.ChildComments {
				childLikesCount, _ := s.commentRepo.CountLikes(child.ID)
				childComments[j] = *mapCommentToCommentResponse(&child, &child.User, int(childLikesCount), false, nil)
			}
		}

		commentResponses[i] = *mapCommentToCommentResponse(&comment, &comment.User, int(likesCount), false, childComments)
	}

	return &response.CommentListResponse{
		Comments: commentResponses,
		Pagination: response.PaginationResponse{
			Total:  int(total),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

// DeleteComment コメント削除
func (s *commentService) DeleteComment(userID uuid.UUID, commentID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("コメントが見つかりません")
	}

	// 権限チェック
	if comment.UserID != userID {
		return errors.New("このコメントを削除する権限がありません")
	}

	return s.commentRepo.Delete(commentID)
}

// mapCommentToCommentResponse CommentをCommentResponseにマッピング
func mapCommentToCommentResponse(
	comment *model.Comment,
	user *model.User,
	likesCount int,
	isLiked bool,
	childComments []response.CommentResponse,
) *response.CommentResponse {
	return &response.CommentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		Content:    comment.Content,
		LikesCount: likesCount,
		IsLiked:    isLiked,
		User: response.UserSimple{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			AvatarURL:   user.AvatarURL,
		},
		ChildComments: childComments,
		CreatedAt:     comment.CreatedAt,
		UpdatedAt:     comment.UpdatedAt,
	}
}
