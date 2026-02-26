package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// FollowService フォローサービスのインターフェース
type FollowService interface {
	Follow(followerID uuid.UUID, username string) error
	Unfollow(followerID uuid.UUID, username string) error
	GetFollowers(username string, limit, offset int, currentUserID *uuid.UUID) (*response.UserListResponse, error)
	GetFollowing(username string, limit, offset int, currentUserID *uuid.UUID) (*response.UserListResponse, error)
}

type followService struct {
	followRepo          repository.FollowRepository
	userRepo            repository.UserRepository
	notificationService NotificationService
}

// NewFollowService フォローサービスのコンストラクタ
func NewFollowService(
	followRepo repository.FollowRepository,
	userRepo repository.UserRepository,
	notificationService NotificationService,
) FollowService {
	return &followService{
		followRepo:          followRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

// Follow フォロー
func (s *followService) Follow(followerID uuid.UUID, username string) error {
	// フォロー対象のユーザーを取得
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("ユーザーが見つかりません")
	}

	// 自分自身をフォローしようとしていないか確認
	if followerID == user.ID {
		return errors.New("自分自身をフォローすることはできません")
	}

	// すでにフォローしているか確認
	exists, err := s.followRepo.Exists(followerID, user.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("すでにフォローしています")
	}

	follow := &model.Follow{
		FollowerID: followerID,
		FollowedID: user.ID,
	}

	// フォローを作成
	if err := s.followRepo.Create(follow); err != nil {
		return err
	}

	// 通知を作成（自分自身のフォローは除外される）
	s.notificationService.CreateFollowNotification(followerID, user.ID)

	return nil
}

// Unfollow フォロー解除
func (s *followService) Unfollow(followerID uuid.UUID, username string) error {
	// フォロー対象のユーザーを取得
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("ユーザーが見つかりません")
	}

	// フォローしているか確認
	exists, err := s.followRepo.Exists(followerID, user.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("フォローしていません")
	}

	// フォロー解除
	if err := s.followRepo.Delete(followerID, user.ID); err != nil {
		return err
	}

	// 通知を削除
	s.notificationService.DeleteNotificationByAction(followerID, user.ID, "follow", "", uuid.Nil)

	return nil
}

// GetFollowers フォロワー一覧取得
func (s *followService) GetFollowers(username string, limit, offset int, currentUserID *uuid.UUID) (*response.UserListResponse, error) {
	// ユーザーを取得
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	follows, total, err := s.followRepo.FindFollowersByUserID(user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	// UserSimpleに変換
	users := make([]response.UserSimple, len(follows))
	for i, follow := range follows {
		var isFollowing bool
		if currentUserID != nil && *currentUserID != follow.Follower.ID {
			isFollowing, _ = s.followRepo.Exists(*currentUserID, follow.Follower.ID)
		}

		users[i] = response.UserSimple{
			ID:          follow.Follower.ID,
			Username:    follow.Follower.Username,
			DisplayName: follow.Follower.DisplayName,
			Bio:         follow.Follower.Bio,
			AvatarURL:   follow.Follower.AvatarURL,
			IsFollowing: &isFollowing,
		}
	}

	hasMore := offset+limit < int(total)

	return &response.UserListResponse{
		Users: users,
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}

// GetFollowing フォロー中一覧取得
func (s *followService) GetFollowing(username string, limit, offset int, currentUserID *uuid.UUID) (*response.UserListResponse, error) {
	// ユーザーを取得
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	follows, total, err := s.followRepo.FindFollowingByUserID(user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	// UserSimpleに変換
	users := make([]response.UserSimple, len(follows))
	for i, follow := range follows {
		var isFollowing bool
		if currentUserID != nil && *currentUserID != follow.Followed.ID {
			isFollowing, _ = s.followRepo.Exists(*currentUserID, follow.Followed.ID)
		}

		users[i] = response.UserSimple{
			ID:          follow.Followed.ID,
			Username:    follow.Followed.Username,
			DisplayName: follow.Followed.DisplayName,
			Bio:         follow.Followed.Bio,
			AvatarURL:   follow.Followed.AvatarURL,
			IsFollowing: &isFollowing,
		}
	}

	hasMore := offset+limit < int(total)

	return &response.UserListResponse{
		Users: users,
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}
