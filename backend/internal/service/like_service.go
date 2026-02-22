package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// LikeService いいねサービスのインターフェース
type LikeService interface {
	LikePost(userID uuid.UUID, postID uuid.UUID) error
	UnlikePost(userID uuid.UUID, postID uuid.UUID) error
	LikeComment(userID uuid.UUID, commentID uuid.UUID) error
	UnlikeComment(userID uuid.UUID, commentID uuid.UUID) error
}

type likeService struct {
	likeRepo    repository.LikeRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
}

// NewLikeService いいねサービスのコンストラクタ
func NewLikeService(
	likeRepo repository.LikeRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
) LikeService {
	return &likeService{
		likeRepo:    likeRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

// LikePost 投稿にいいね
func (s *likeService) LikePost(userID uuid.UUID, postID uuid.UUID) error {
	// 投稿が存在するか確認
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("投稿が見つかりません")
	}

	// すでにいいねしているか確認
	exists, err := s.likeRepo.Exists(userID, "Post", postID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("すでにいいねしています")
	}

	like := &model.Like{
		UserID:       userID,
		LikeableType: "Post",
		LikeableID:   postID,
	}

	return s.likeRepo.Create(like)
}

// UnlikePost 投稿のいいねを解除
func (s *likeService) UnlikePost(userID uuid.UUID, postID uuid.UUID) error {
	// いいねしているか確認
	exists, err := s.likeRepo.Exists(userID, "Post", postID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("いいねしていません")
	}

	return s.likeRepo.Delete(userID, "Post", postID)
}

// LikeComment コメントにいいね
func (s *likeService) LikeComment(userID uuid.UUID, commentID uuid.UUID) error {
	// コメントが存在するか確認
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("コメントが見つかりません")
	}

	// すでにいいねしているか確認
	exists, err := s.likeRepo.Exists(userID, "Comment", commentID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("すでにいいねしています")
	}

	like := &model.Like{
		UserID:       userID,
		LikeableType: "Comment",
		LikeableID:   commentID,
	}

	return s.likeRepo.Create(like)
}

// UnlikeComment コメントのいいねを解除
func (s *likeService) UnlikeComment(userID uuid.UUID, commentID uuid.UUID) error {
	// いいねしているか確認
	exists, err := s.likeRepo.Exists(userID, "Comment", commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("いいねしていません")
	}

	return s.likeRepo.Delete(userID, "Comment", commentID)
}
