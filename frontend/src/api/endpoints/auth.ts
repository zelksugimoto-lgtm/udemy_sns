import { apiClient } from '../client';
import type { components } from '../generated/schema';

// 型定義
type RegisterRequest = components['schemas']['request.RegisterRequest'];
type LoginRequest = components['schemas']['request.LoginRequest'];
type AuthResponse = components['schemas']['response.AuthResponse'];
type UserResponse = components['schemas']['response.UserResponse'];

// ユーザー登録
export const register = async (data: RegisterRequest): Promise<AuthResponse> => {
  const response = await apiClient.post<AuthResponse>('/auth/register', data);
  return response.data;
};

// ログイン
export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  const response = await apiClient.post<AuthResponse>('/auth/login', data);
  return response.data;
};

// 現在のユーザー情報取得
export const getCurrentUser = async (): Promise<UserResponse> => {
  const response = await apiClient.get<UserResponse>('/auth/me');
  return response.data;
};
