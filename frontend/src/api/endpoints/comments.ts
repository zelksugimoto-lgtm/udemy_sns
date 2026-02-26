import { apiClient } from '../client';
import type { components } from '../generated/schema';
import { PAGINATION } from '../../utils/constants';

type CreateCommentRequest = components['schemas']['request.CreateCommentRequest'];
type CommentResponse = components['schemas']['response.CommentResponse'];
type CommentListResponse = components['schemas']['response.CommentListResponse'];

interface PaginationParams {
  limit?: number;
  offset?: number;
}

// コメント作成
export const createComment = async (
  postId: string,
  data: CreateCommentRequest
): Promise<CommentResponse> => {
  const response = await apiClient.post<CommentResponse>(`/posts/${postId}/comments`, data);
  return response.data;
};

// コメント一覧取得
export const getComments = async (
  postId: string,
  params: PaginationParams = {}
): Promise<CommentListResponse> => {
  const response = await apiClient.get<CommentListResponse>(`/posts/${postId}/comments`, {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// コメント削除
export const deleteComment = async (id: string): Promise<void> => {
  await apiClient.delete(`/comments/${id}`);
};

// コメントにいいね
export const likeComment = async (id: string): Promise<void> => {
  await apiClient.post(`/comments/${id}/like`);
};

// コメントのいいね解除
export const unlikeComment = async (id: string): Promise<void> => {
  await apiClient.delete(`/comments/${id}/like`);
};
