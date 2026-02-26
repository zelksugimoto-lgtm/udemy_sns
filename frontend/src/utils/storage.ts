import { STORAGE_KEYS } from './constants';

// トークン管理
export const tokenStorage = {
  get: (): string | null => {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
  },
  set: (token: string): void => {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token);
  },
  remove: (): void => {
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
  },
};

// ユーザー情報管理
export const userStorage = {
  get: <T>(): T | null => {
    const user = localStorage.getItem(STORAGE_KEYS.USER);
    return user ? JSON.parse(user) : null;
  },
  set: <T>(user: T): void => {
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user));
  },
  remove: (): void => {
    localStorage.removeItem(STORAGE_KEYS.USER);
  },
};

// テーマ管理
export const themeStorage = {
  get: (): string | null => {
    return localStorage.getItem(STORAGE_KEYS.THEME);
  },
  set: (theme: string): void => {
    localStorage.setItem(STORAGE_KEYS.THEME, theme);
  },
  remove: (): void => {
    localStorage.removeItem(STORAGE_KEYS.THEME);
  },
};

// 全データクリア
export const clearAllStorage = (): void => {
  tokenStorage.remove();
  userStorage.remove();
  // テーマは保持
};
