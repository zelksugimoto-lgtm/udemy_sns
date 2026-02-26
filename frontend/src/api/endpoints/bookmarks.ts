import { apiClient } from '../client';
import type { components } from '../generated/schema';
import { PAGINATION } from '../../utils/constants';

type BookmarkListResponse = components['schemas']['response.BookmarkListResponse'];

interface PaginationParams {
  limit?: number;
  offset?: number;
}

// ブックマーク一覧取得
export const getBookmarks = async (
  params: PaginationParams = {}
): Promise<BookmarkListResponse> => {
  const response = await apiClient.get<BookmarkListResponse>('/bookmarks', {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};
