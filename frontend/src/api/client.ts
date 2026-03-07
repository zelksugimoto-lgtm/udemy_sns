import axios from 'axios';
import type { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { API_BASE_URL } from '../utils/constants';
import { clearAllStorage } from '../utils/storage';

// Axiosインスタンス作成（Cookie-based認証対応）
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
  withCredentials: true, // Cookieを自動送信
});

// リフレッシュ処理中フラグ（多重リフレッシュ防止）
let isRefreshing = false;
let refreshSubscribers: Array<(token?: string) => void> = [];

// リフレッシュ完了時に待機中のリクエストを再実行
const onRefreshed = () => {
  refreshSubscribers.forEach((callback) => callback());
  refreshSubscribers = [];
};

// リフレッシュ失敗時
const onRefreshError = () => {
  refreshSubscribers = [];
};

// リクエストを待機列に追加
const addRefreshSubscriber = (callback: () => void) => {
  refreshSubscribers.push(callback);
};

// レスポンスインターセプター（エラーハンドリング + 自動リフレッシュ）
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // 401エラー: 認証エラー → トークンリフレッシュを試みる
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      // リフレッシュAPIへのリクエスト自体が401の場合はログアウト
      if (originalRequest.url?.includes('/auth/refresh')) {
        clearAllStorage();
        if (window.location.pathname !== '/login' && window.location.pathname !== '/register') {
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }

      // 既にリフレッシュ処理中の場合は待機
      if (isRefreshing) {
        return new Promise((resolve) => {
          addRefreshSubscriber(() => {
            resolve(apiClient(originalRequest));
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // リフレッシュトークンでアクセストークンを再発行
        await axios.post(
          `${API_BASE_URL}/auth/refresh`,
          {},
          {
            withCredentials: true, // Cookie送信
          }
        );

        // リフレッシュ成功: 元のリクエストを再実行
        isRefreshing = false;
        onRefreshed();
        return apiClient(originalRequest);
      } catch (refreshError) {
        // リフレッシュ失敗: ログアウトしてログインページへ
        isRefreshing = false;
        onRefreshError();
        clearAllStorage();
        if (window.location.pathname !== '/login' && window.location.pathname !== '/register') {
          window.location.href = '/login';
        }
        return Promise.reject(refreshError);
      }
    }

    // 403エラー: 権限エラー
    if (error.response?.status === 403) {
      console.error('Permission denied');
    }

    // 500エラー: サーバーエラー
    if (error.response?.status === 500) {
      console.error('Server error');
    }

    return Promise.reject(error);
  }
);

// エラーメッセージ抽出ヘルパー
export const getErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { error?: { message?: string } };
    return data?.error?.message || error.message || 'An error occurred';
  }
  if (error instanceof Error) {
    return error.message;
  }
  return 'An unknown error occurred';
};
