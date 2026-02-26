package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/dto/response"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// UserService ユーザーサービスのインターフェース
type UserService interface {
	GetProfile(username string, currentUserID *uuid.UUID) (*response.UserProfileResponse, error)
	UpdateProfile(userID uuid.UUID, req *request.UpdateProfileRequest) (*response.UserResponse, error)
	SearchUsers(query string, limit, offset int, currentUserID *uuid.UUID) (*response.UserListResponse, error)
}

type userService struct {
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
	postRepo   repository.PostRepository
}

// NewUserService ユーザーサービスのコンストラクタ
func NewUserService(
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	postRepo repository.PostRepository,
) UserService {
	return &userService{
		userRepo:   userRepo,
		followRepo: followRepo,
		postRepo:   postRepo,
	}
}

// GetProfile プロフィール取得
func (s *userService) GetProfile(username string, currentUserID *uuid.UUID) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	// フォロワー数取得
	_, followersCount, err := s.followRepo.FindFollowersByUserID(user.ID, 1, 0)
	if err != nil {
		return nil, err
	}

	// フォロー中数取得
	_, followingCount, err := s.followRepo.FindFollowingByUserID(user.ID, 1, 0)
	if err != nil {
		return nil, err
	}

	// 投稿数取得
	_, postsCount, err := s.postRepo.FindByUserID(user.ID, 1, 0)
	if err != nil {
		return nil, err
	}

	// フォロー状態確認
	var isFollowing bool
	if currentUserID != nil && *currentUserID != user.ID {
		isFollowing, err = s.followRepo.Exists(*currentUserID, user.ID)
		if err != nil {
			return nil, err
		}
	}

	return &response.UserProfileResponse{
		ID:             user.ID,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Bio:            user.Bio,
		AvatarURL:      user.AvatarURL,
		HeaderURL:      user.HeaderURL,
		FollowersCount: int(followersCount),
		FollowingCount: int(followingCount),
		PostsCount:     int(postsCount),
		IsFollowing:    &isFollowing,
		CreatedAt:      user.CreatedAt,
	}, nil
}

// UpdateProfile プロフィール更新
func (s *userService) UpdateProfile(userID uuid.UUID, req *request.UpdateProfileRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	// 更新
	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}
	if req.HeaderURL != nil {
		user.HeaderURL = *req.HeaderURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		HeaderURL:   user.HeaderURL,
		CreatedAt:   user.CreatedAt,
	}, nil
}

// SearchUsers ユーザー検索
func (s *userService) SearchUsers(query string, limit, offset int, currentUserID *uuid.UUID) (*response.UserListResponse, error) {
	users, total, err := s.userRepo.Search(query, limit, offset)
	if err != nil {
		return nil, err
	}

	userSimples := make([]response.UserSimple, len(users))
	for i, user := range users {
		userSimples[i] = s.mapUserToUserSimple(&user, currentUserID)
	}

	hasMore := offset+limit < int(total)

	return &response.UserListResponse{
		Users: userSimples,
		Pagination: response.PaginationResponse{
			Total:   int(total),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}, nil
}

// mapUserToUserSimple UserをUserSimpleにマッピング
func (s *userService) mapUserToUserSimple(user *model.User, currentUserID *uuid.UUID) response.UserSimple {
	var isFollowing bool
	if currentUserID != nil && *currentUserID != user.ID {
		isFollowing, _ = s.followRepo.Exists(*currentUserID, user.ID)
	}

	return response.UserSimple{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		IsFollowing: &isFollowing,
	}
}
