import { apiClient } from '../client';
import type { components } from '../generated/schema';
import { PAGINATION } from '../../utils/constants';

type UpdateProfileRequest = components['schemas']['request.UpdateProfileRequest'];
type UserProfileResponse = components['schemas']['response.UserProfileResponse'];
type UserListResponse = components['schemas']['response.UserListResponse'];
type UserResponse = components['schemas']['response.UserResponse'];
type PostListResponse = components['schemas']['response.PostListResponse'];

interface PaginationParams {
  limit?: number;
  offset?: number;
}

// ユーザープロフィール取得
export const getUserProfile = async (username: string): Promise<UserProfileResponse> => {
  const response = await apiClient.get<UserProfileResponse>(`/users/${username}`);
  return response.data;
};

// プロフィール更新
export const updateProfile = async (data: UpdateProfileRequest): Promise<UserResponse> => {
  const response = await apiClient.patch<UserResponse>('/users/me', data);
  return response.data;
};

interface SearchParams extends PaginationParams {
  query: string;
}

// ユーザー検索
export const searchUsers = async (
  params: SearchParams
): Promise<UserListResponse> => {
  const response = await apiClient.get<UserListResponse>('/users', {
    params: {
      q: params.query,
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// フォロー
export const followUser = async (username: string): Promise<void> => {
  await apiClient.post(`/users/${username}/follow`);
};

// フォロー解除
export const unfollowUser = async (username: string): Promise<void> => {
  await apiClient.delete(`/users/${username}/follow`);
};

// フォロワー一覧取得
export const getFollowers = async (
  username: string,
  params: PaginationParams = {}
): Promise<UserListResponse> => {
  const response = await apiClient.get<UserListResponse>(`/users/${username}/followers`, {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// フォロー中一覧取得
export const getFollowing = async (
  username: string,
  params: PaginationParams = {}
): Promise<UserListResponse> => {
  const response = await apiClient.get<UserListResponse>(`/users/${username}/following`, {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// ユーザーがいいねした投稿一覧取得
export const getUserLikedPosts = async (
  username: string,
  params: PaginationParams = {}
): Promise<PostListResponse> => {
  const response = await apiClient.get<PostListResponse>(`/users/${username}/likes`, {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};
