import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import type { components } from '../api/generated/schema';
import * as authApi from '../api/endpoints/auth';
import { tokenStorage, userStorage, clearAllStorage } from '../utils/storage';
import { getErrorMessage } from '../api/client';

type User = components['schemas']['response.UserResponse'];
type RegisterRequest = components['schemas']['request.RegisterRequest'];
type LoginRequest = components['schemas']['request.LoginRequest'];

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
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
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 初期化: ローカルストレージからトークンとユーザー情報を復元
  useEffect(() => {
    const initAuth = async () => {
      const savedToken = tokenStorage.get();
      const savedUser = userStorage.get<User>();

      if (savedToken && savedUser) {
        setToken(savedToken);
        setUser(savedUser);

        // トークンが有効か確認（現在のユーザー情報を取得）
        try {
          const currentUser = await authApi.getCurrentUser();
          setUser(currentUser);
          userStorage.set(currentUser);
        } catch (err) {
          // トークンが無効な場合はクリア
          clearAllStorage();
          setToken(null);
          setUser(null);
        }
      }

      setIsLoading(false);
    };

    initAuth();
  }, []);

  // ログイン
  const login = useCallback(async (data: LoginRequest) => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await authApi.login(data);

      if (response.token && response.user) {
        setToken(response.token);
        setUser(response.user);
        tokenStorage.set(response.token);
        userStorage.set(response.user);
      }
    } catch (err) {
      const message = getErrorMessage(err);
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // ユーザー登録
  const register = useCallback(async (data: RegisterRequest) => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await authApi.register(data);

      if (response.token && response.user) {
        setToken(response.token);
        setUser(response.user);
        tokenStorage.set(response.token);
        userStorage.set(response.user);
      }
    } catch (err) {
      const message = getErrorMessage(err);
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // ログアウト
  const logout = useCallback(() => {
    clearAllStorage();
    setToken(null);
    setUser(null);
    setError(null);
  }, []);

  const value = {
    user,
    token,
    isAuthenticated: !!token && !!user,
    isLoading,
    login,
    register,
    logout,
    error,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
