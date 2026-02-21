# 認証フロー仕様書

## 概要

JWT (JSON Web Token) を使用したステートレス認証

## 認証フロー図

### ユーザー登録フロー

```
[クライアント]
    ↓ 1. 登録フォーム送信
    POST /api/v1/auth/register
    {
      email: "user@example.com",
      password: "password123",
      username: "johndoe",
      display_name: "John Doe"
    }
    ↓
[バックエンド]
    ↓ 2. バリデーション
    - メールアドレス形式チェック
    - パスワード強度チェック (8文字以上)
    - ユーザー名形式チェック (3-50文字、英数字_のみ)
    ↓ 3. 重複チェック
    - メールアドレスの重複確認
    - ユーザー名の重複確認
    ↓ 4. パスワードハッシュ化
    - bcrypt でハッシュ化 (cost: 10)
    ↓ 5. ユーザーをDBに保存
    ↓ 6. JWT生成
    - payload: { user_id, email }
    - 有効期限: 7日間
    ↓ 7. レスポンス返却
    {
      token: "eyJhbGciOiJIUzI1NiIs...",
      user: { id, email, username, ... }
    }
    ↓
[クライアント]
    ↓ 8. トークン保存
    - localStorage に保存
    ↓ 9. ホームページへリダイレクト
```

---

### ログインフロー

```
[クライアント]
    ↓ 1. ログインフォーム送信
    POST /api/v1/auth/login
    {
      email: "user@example.com",
      password: "password123"
    }
    ↓
[バックエンド]
    ↓ 2. ユーザー検索
    - メールアドレスでユーザー検索
    - deleted_at IS NULL を確認
    ↓ 3. パスワード検証
    - bcrypt.CompareHashAndPassword
    ↓ 4. JWT生成
    - payload: { user_id, email }
    - 有効期限: 7日間
    ↓ 5. レスポンス返却
    {
      token: "eyJhbGciOiJIUzI1NiIs...",
      user: { id, email, username, ... }
    }
    ↓
[クライアント]
    ↓ 6. トークン保存
    - localStorage に保存
    ↓ 7. ホームページへリダイレクト
```

---

### 認証済みAPIリクエストフロー

```
[クライアント]
    ↓ 1. APIリクエスト
    GET /api/v1/timeline
    Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
    ↓
[バックエンド - ミドルウェア]
    ↓ 2. トークン抽出
    - Authorizationヘッダーから "Bearer " を除去
    ↓ 3. トークン検証
    - JWT署名検証
    - 有効期限確認
    ↓ 4. ユーザー情報取得
    - payload から user_id を抽出
    - DBからユーザー情報取得
    - deleted_at IS NULL を確認
    ↓ 5. コンテキストに保存
    - c.Set("user", user)
    ↓
[バックエンド - ハンドラー]
    ↓ 6. ユーザー情報取得
    - user := c.Get("user").(*model.User)
    ↓ 7. ビジネスロジック実行
    ↓ 8. レスポンス返却
```

---

## JWT構造

### Header
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Payload
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "exp": 1735689600,
  "iat": 1735084800
}
```

### Signature
```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  JWT_SECRET
)
```

---

## バックエンド実装

### JWT生成 (Go)

```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateToken(userID, email, secret string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### JWT検証 (Go)

```go
func ValidateToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        },
    )

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

### 認証ミドルウェア (Echo)

```go
func AuthMiddleware(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := ValidateToken(tokenString, secret)
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }

            // ユーザー情報取得
            var user model.User
            if err := db.First(&user, "id = ? AND deleted_at IS NULL", claims.UserID).Error; err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
            }

            c.Set("user", &user)
            return next(c)
        }
    }
}
```

---

## フロントエンド実装

### トークン保存・取得

```typescript
// utils/storage.ts
const TOKEN_KEY = 'auth_token';

export const storage = {
  setToken: (token: string) => {
    localStorage.setItem(TOKEN_KEY, token);
  },

  getToken: (): string | null => {
    return localStorage.getItem(TOKEN_KEY);
  },

  removeToken: () => {
    localStorage.removeItem(TOKEN_KEY);
  },
};
```

### Axios インターセプター

```typescript
// api/client.ts
import axios from 'axios';
import { storage } from '@/utils/storage';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
});

// リクエストインターセプター
client.interceptors.request.use((config) => {
  const token = storage.getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// レスポンスインターセプター
client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      storage.removeToken();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default client;
```

### AuthContext

```typescript
// contexts/AuthContext.tsx
import { createContext, useContext, useState, useEffect } from 'react';
import { storage } from '@/utils/storage';
import * as authAPI from '@/api/endpoints/auth';

type User = {
  id: string;
  email: string;
  username: string;
  display_name: string;
};

type AuthContextType = {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initAuth = async () => {
      const token = storage.getToken();
      if (token) {
        try {
          const { data } = await authAPI.getMe();
          setUser(data);
        } catch (error) {
          storage.removeToken();
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await authAPI.login({ email, password });
    storage.setToken(data.token);
    setUser(data.user);
  };

  const register = async (data: RegisterData) => {
    const { data: responseData } = await authAPI.register(data);
    storage.setToken(responseData.token);
    setUser(responseData.user);
  };

  const logout = () => {
    storage.removeToken();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
```

### Protected Route

```typescript
// components/common/ProtectedRoute.tsx
import { Navigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';

export const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <CircularProgress />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
};
```

---

## セキュリティ対策

### バックエンド
1. **パスワードハッシュ化**: bcrypt (cost: 10以上)
2. **JWT署名**: HS256アルゴリズム、強力なシークレットキー使用
3. **トークン有効期限**: 7日間
4. **論理削除チェック**: deleted_at IS NULL 確認
5. **HTTPS通信**: 本番環境では必須

### フロントエンド
1. **トークン保存**: localStorage（XSS対策としてhttpOnlyクッキー検討）
2. **自動ログアウト**: トークン期限切れ時
3. **入力サニタイゼーション**: XSS対策

---

## エラーハンドリング

### 登録エラー
- **409 Conflict**: メールアドレス or ユーザー名が既に存在
- **400 Bad Request**: バリデーションエラー

### ログインエラー
- **401 Unauthorized**: メールアドレス or パスワードが間違っている
- **404 Not Found**: ユーザーが存在しない

### 認証エラー
- **401 Unauthorized**: トークンが無効 or 期限切れ
- **403 Forbidden**: 権限がない

---

## 将来の拡張（Phase 4以降）

1. **リフレッシュトークン**: アクセストークン（短時間）+ リフレッシュトークン（長期間）
2. **ソーシャルログイン**: OAuth 2.0 (Google, Twitter)
3. **2要素認証 (2FA)**: TOTP (Time-based One-Time Password)
4. **パスワードリセット**: メール認証
5. **メール確認**: 登録時のメール認証
