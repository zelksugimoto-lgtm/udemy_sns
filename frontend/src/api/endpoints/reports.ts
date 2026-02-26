import { apiClient } from '../client';
import type { components } from '../generated/schema';

type CreateReportRequest = components['schemas']['request.CreateReportRequest'];

// 通報作成
export const createReport = async (data: CreateReportRequest): Promise<void> => {
  await apiClient.post('/reports', data);
};
