import { apiClient } from '../client';
import type { components } from '../generated/schema';
import { PAGINATION } from '../../utils/constants';

// 型定義
type CreatePostRequest = components['schemas']['request.CreatePostRequest'];
type UpdatePostRequest = components['schemas']['request.UpdatePostRequest'];
type PostResponse = components['schemas']['response.PostResponse'];
type PostListResponse = components['schemas']['response.PostListResponse'];

// ページネーションパラメータ
interface PaginationParams {
  limit?: number;
  offset?: number;
}

// 投稿作成
export const createPost = async (data: CreatePostRequest): Promise<PostResponse> => {
  const response = await apiClient.post<PostResponse>('/posts', data);
  return response.data;
};

// 投稿取得
export const getPost = async (id: string): Promise<PostResponse> => {
  const response = await apiClient.get<PostResponse>(`/posts/${id}`);
  return response.data;
};

// 投稿更新
export const updatePost = async (id: string, data: UpdatePostRequest): Promise<PostResponse> => {
  const response = await apiClient.patch<PostResponse>(`/posts/${id}`, data);
  return response.data;
};

// 投稿削除
export const deletePost = async (id: string): Promise<void> => {
  await apiClient.delete(`/posts/${id}`);
};

// タイムライン取得
export const getTimeline = async (params: PaginationParams = {}): Promise<PostListResponse> => {
  const response = await apiClient.get<PostListResponse>('/timeline', {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// ユーザーの投稿一覧取得
export const getUserPosts = async (
  username: string,
  params: PaginationParams = {}
): Promise<PostListResponse> => {
  const response = await apiClient.get<PostListResponse>(`/users/${username}/posts`, {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// 投稿にいいね
export const likePost = async (id: string): Promise<void> => {
  await apiClient.post(`/posts/${id}/like`);
};

// 投稿のいいね解除
export const unlikePost = async (id: string): Promise<void> => {
  await apiClient.delete(`/posts/${id}/like`);
};

// ブックマーク追加
export const bookmarkPost = async (id: string): Promise<void> => {
  await apiClient.post(`/posts/${id}/bookmark`);
};

// ブックマーク削除
export const unbookmarkPost = async (id: string): Promise<void> => {
  await apiClient.delete(`/posts/${id}/bookmark`);
};
