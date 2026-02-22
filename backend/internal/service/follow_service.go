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
	GetFollowers(username string, limit, offset int) (*response.FollowListResponse, error)
	GetFollowing(username string, limit, offset int) (*response.FollowListResponse, error)
}

type followService struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

// NewFollowService フォローサービスのコンストラクタ
func NewFollowService(
	followRepo repository.FollowRepository,
	userRepo repository.UserRepository,
) FollowService {
	return &followService{
		followRepo: followRepo,
		userRepo:   userRepo,
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

	return s.followRepo.Create(follow)
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

	return s.followRepo.Delete(followerID, user.ID)
}

// GetFollowers フォロワー一覧取得
func (s *followService) GetFollowers(username string, limit, offset int) (*response.FollowListResponse, error) {
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
		users[i] = response.UserSimple{
			ID:          follow.Follower.ID,
			Username:    follow.Follower.Username,
			DisplayName: follow.Follower.DisplayName,
			AvatarURL:   follow.Follower.AvatarURL,
		}
	}

	return &response.FollowListResponse{
		Users: users,
		Pagination: response.PaginationResponse{
			Total:  int(total),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

// GetFollowing フォロー中一覧取得
func (s *followService) GetFollowing(username string, limit, offset int) (*response.FollowListResponse, error) {
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
		users[i] = response.UserSimple{
			ID:          follow.Followed.ID,
			Username:    follow.Followed.Username,
			DisplayName: follow.Followed.DisplayName,
			AvatarURL:   follow.Followed.AvatarURL,
		}
	}

	return &response.FollowListResponse{
		Users: users,
		Pagination: response.PaginationResponse{
			Total:  int(total),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}
