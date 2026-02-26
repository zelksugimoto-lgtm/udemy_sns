// API関連
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

// ストレージキー
export const STORAGE_KEYS = {
  TOKEN: 'auth_token',
  THEME: 'app_theme',
  USER: 'current_user',
} as const;

// ページネーション
export const PAGINATION = {
  DEFAULT_LIMIT: 20,
  DEFAULT_OFFSET: 0,
} as const;

// 文字数制限
export const CHAR_LIMITS = {
  POST_CONTENT: 250,
  COMMENT_CONTENT: 500,
  DISPLAY_NAME: 100,
  USERNAME: 50,
  BIO: 1000,
} as const;

// パスワードポリシー
export const PASSWORD_POLICY = {
  MIN_LENGTH: 8,
} as const;

// ユーザー名ポリシー
export const USERNAME_POLICY = {
  MIN_LENGTH: 3,
  MAX_LENGTH: 50,
  PATTERN: /^[a-zA-Z0-9_]+$/,
} as const;

// 通知ポーリング間隔 (ミリ秒)
export const NOTIFICATION_POLL_INTERVAL = 60000; // 1分

// テーマ
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark',
  CUSTOM1: 'custom1',
  CUSTOM2: 'custom2',
} as const;

export type ThemeType = typeof THEMES[keyof typeof THEMES];
