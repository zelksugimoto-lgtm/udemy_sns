import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react';
import type { components } from '../api/generated/schema';
import * as authApi from '../api/endpoints/auth';
import { apiClient } from '../api/client';
import { userStorage, clearAllStorage } from '../utils/storage';
import { getErrorMessage } from '../api/client';

type User = components['schemas']['response.UserResponse'];
type RegisterRequest = components['schemas']['request.RegisterRequest'];
type LoginRequest = components['schemas']['request.LoginRequest'];

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (data: LoginRequest) => Promise<User>;
  register: (data: RegisterRequest) => Promise<User>;
  logout: () => Promise<void>;
  error: string | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 初期化: Cookie内のトークンが有効かAPIで確認
  useEffect(() => {
    const initAuth = async () => {
      // ローカルストレージにキャッシュされたユーザー情報があれば一時的にセット
      const savedUser = userStorage.get<User>();
      if (savedUser) {
        setUser(savedUser);
      }

      // Cookie内のトークンで現在のユーザー情報を取得
      try {
        const currentUser = await authApi.getCurrentUser();
        setUser(currentUser);
        userStorage.set(currentUser);
      } catch (err) {
        // Cookie内のトークンが無効または存在しない
        clearAllStorage();
        setUser(null);
      }

      setIsLoading(false);
    };

    initAuth();
  }, []);

  // ログイン
  const login = useCallback(async (data: LoginRequest): Promise<User> => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await authApi.login(data);

      // トークンはCookieに自動保存されるため、userのみ保存
      if (response.user) {
        setUser(response.user);
        userStorage.set(response.user);
        return response.user;
      }
      throw new Error('ユーザー情報が取得できませんでした');
    } catch (err) {
      const message = getErrorMessage(err);
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // ユーザー登録
  const register = useCallback(async (data: RegisterRequest): Promise<User> => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await authApi.register(data);

      // トークンはCookieに自動保存されるため、userのみ保存
      if (response.user) {
        setUser(response.user);
        userStorage.set(response.user);
        return response.user;
      }
      throw new Error('ユーザー情報が取得できませんでした');
    } catch (err) {
      const message = getErrorMessage(err);
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // ログアウト
  const logout = useCallback(async () => {
    try {
      // バックエンドのログアウトAPIを叩いてリフレッシュトークンを無効化
      await apiClient.post('/auth/logout');
    } catch (err) {
      console.error('Logout API error:', err);
      // エラーが発生してもローカルの状態はクリア
    } finally {
      clearAllStorage();
      setUser(null);
      setError(null);
    }
  }, []);

  const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
    error,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
