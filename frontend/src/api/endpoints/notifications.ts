import { apiClient } from '../client';
import type { components } from '../generated/schema';
import { PAGINATION } from '../../utils/constants';

type NotificationListResponse = components['schemas']['response.NotificationListResponse'];

interface PaginationParams {
  limit?: number;
  offset?: number;
}

// 通知一覧取得
export const getNotifications = async (
  params: PaginationParams = {}
): Promise<NotificationListResponse> => {
  const response = await apiClient.get<NotificationListResponse>('/notifications', {
    params: {
      limit: params.limit || PAGINATION.DEFAULT_LIMIT,
      offset: params.offset || PAGINATION.DEFAULT_OFFSET,
    },
  });
  return response.data;
};

// 通知を既読にする
export const markNotificationAsRead = async (id: string): Promise<void> => {
  await apiClient.patch(`/notifications/${id}/read`);
};

// すべての通知を既読にする
export const markAllNotificationsAsRead = async (): Promise<void> => {
  await apiClient.post('/notifications/read-all');
};

// 未読通知数取得
export const getUnreadCount = async (): Promise<number> => {
  const response = await apiClient.get<{ count: number }>('/notifications/unread-count');
  return response.data.count;
};
