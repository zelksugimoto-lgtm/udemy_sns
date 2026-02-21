# フロントエンド開発ルール

## 型安全開発

### openapi-typescript による型生成
```bash
# バックエンドのOpenAPI定義から型を生成
npx openapi-typescript http://localhost:8080/swagger/doc.json -o src/api/generated/schema.ts
```

### 生成された型の使用
```typescript
import type { components } from '@/api/generated/schema';

type User = components['schemas']['UserResponse'];
type Post = components['schemas']['PostResponse'];
```

---

## プロジェクト構造

```
frontend/src/
├── api/                    # API通信
│   ├── client.ts          # Axios設定
│   ├── generated/         # openapi-typescript生成
│   └── endpoints/         # APIエンドポイント
├── components/            # UIコンポーネント
│   ├── common/           # 共通コンポーネント
│   ├── post/
│   ├── user/
│   └── ...
├── pages/                # ページコンポーネント
├── hooks/                # カスタムフック
├── contexts/             # Context API
├── utils/                # ユーティリティ
├── theme/                # MUIテーマ
└── types/                # 型定義
```

---

## API通信

### Axios クライアント設定
```typescript
// api/client.ts
import axios from 'axios';
import { storage } from '@/utils/storage';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
});

// リクエストインターセプター（JWT自動付与）
client.interceptors.request.use((config) => {
  const token = storage.getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// レスポンスインターセプター（401エラーハンドリング）
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

### APIエンドポイント定義
```typescript
// api/endpoints/posts.ts
import client from '../client';
import type { components } from '../generated/schema';

type CreatePostRequest = components['schemas']['CreatePostRequest'];
type PostResponse = components['schemas']['PostResponse'];

export const createPost = (data: CreatePostRequest) =>
  client.post<{ data: PostResponse }>('/posts', data);

export const getPost = (id: string) =>
  client.get<{ data: PostResponse }>(`/posts/${id}`);

export const getTimeline = (params: { limit?: number; offset?: number }) =>
  client.get<{ data: PostResponse[]; pagination: any }>('/timeline', { params });
```

---

## 状態管理

### React Query (TanStack Query)
```typescript
// hooks/usePosts.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as postsAPI from '@/api/endpoints/posts';

export const usePosts = () => {
  const queryClient = useQueryClient();

  // 投稿一覧取得
  const { data, isLoading } = useQuery({
    queryKey: ['posts'],
    queryFn: () => postsAPI.getTimeline({ limit: 20, offset: 0 }),
  });

  // 投稿作成
  const createMutation = useMutation({
    mutationFn: postsAPI.createPost,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
    },
  });

  return {
    posts: data?.data.data ?? [],
    isLoading,
    createPost: createMutation.mutate,
  };
};
```

### 無限スクロール
```typescript
import { useInfiniteQuery } from '@tanstack/react-query';

export const useTimelinePosts = () => {
  return useInfiniteQuery({
    queryKey: ['timeline'],
    queryFn: ({ pageParam = 0 }) =>
      postsAPI.getTimeline({ limit: 20, offset: pageParam }),
    getNextPageParam: (lastPage, pages) =>
      lastPage.data.pagination.has_more ? pages.length * 20 : undefined,
  });
};
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
          setUser(data.data);
        } catch {
          storage.removeToken();
        }
      }
      setIsLoading(false);
    };
    initAuth();
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await authAPI.login({ email, password });
    storage.setToken(data.data.token);
    setUser(data.data.user);
  };

  const logout = () => {
    storage.removeToken();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};
```

---

## コンポーネント設計

### Atomic Designの考え方
- **Atoms**: Button, TextField, Avatar
- **Molecules**: PostCard, CommentCard
- **Organisms**: PostList, Timeline
- **Templates**: Layout
- **Pages**: Home, Profile

### コンポーネント例
```typescript
// components/post/PostCard.tsx
import { Card, CardContent, Typography, IconButton } from '@mui/material';
import { Favorite, Comment, Bookmark } from '@mui/icons-material';
import type { components } from '@/api/generated/schema';

type Post = components['schemas']['PostResponse'];

type PostCardProps = {
  post: Post;
  onLike: (id: string) => void;
  onComment: (id: string) => void;
};

export const PostCard = ({ post, onLike, onComment }: PostCardProps) => {
  return (
    <Card>
      <CardContent>
        <Typography variant="h6">{post.user.display_name}</Typography>
        <Typography variant="body2">@{post.user.username}</Typography>
        <Typography variant="body1" sx={{ mt: 2 }}>
          {post.content}
        </Typography>

        <Box sx={{ mt: 2, display: 'flex', gap: 2 }}>
          <IconButton onClick={() => onLike(post.id)}>
            <Favorite color={post.is_liked ? 'error' : 'default'} />
            <Typography variant="caption">{post.like_count}</Typography>
          </IconButton>

          <IconButton onClick={() => onComment(post.id)}>
            <Comment />
            <Typography variant="caption">{post.comment_count}</Typography>
          </IconButton>
        </Box>
      </CardContent>
    </Card>
  );
};
```

---

## レスポンシブデザイン

### MUI Breakpoints
```typescript
import { useTheme, useMediaQuery } from '@mui/material';

const Component = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const isTablet = useMediaQuery(theme.breakpoints.between('sm', 'md'));
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));

  return (
    <Box sx={{
      display: 'flex',
      flexDirection: { xs: 'column', md: 'row' },
      gap: { xs: 1, md: 2 },
    }}>
      {/* コンテンツ */}
    </Box>
  );
};
```

### レイアウト例
```typescript
// components/common/Layout.tsx
export const Layout = ({ children }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  return (
    <Box sx={{ display: 'flex' }}>
      <Header />
      {!isMobile && <Sidebar />}
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        {children}
      </Box>
      {isMobile && <BottomNavigation />}
    </Box>
  );
};
```

---

## テーマ管理

### カスタムテーマ作成
```typescript
// theme/lightTheme.ts
import { createTheme } from '@mui/material/styles';

export const lightTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
  },
});
```

### ThemeProvider
```typescript
// App.tsx
import { ThemeProvider } from '@mui/material/styles';
import { lightTheme } from './theme/lightTheme';

function App() {
  return (
    <ThemeProvider theme={lightTheme}>
      {/* アプリ */}
    </ThemeProvider>
  );
}
```

---

## フォーム処理

### React Hook Form + Yup
```typescript
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

const schema = yup.object({
  email: yup.string().email('Invalid email').required('Email is required'),
  password: yup.string().min(8, 'Password must be at least 8 characters').required(),
});

type FormData = yup.InferType<typeof schema>;

export const LoginForm = () => {
  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: yupResolver(schema),
  });

  const onSubmit = (data: FormData) => {
    // ログイン処理
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <TextField
        {...register('email')}
        label="Email"
        error={!!errors.email}
        helperText={errors.email?.message}
      />
      <TextField
        {...register('password')}
        type="password"
        label="Password"
        error={!!errors.password}
        helperText={errors.password?.message}
      />
      <Button type="submit">Login</Button>
    </form>
  );
};
```

---

## パフォーマンス最適化

### コード分割
```typescript
import { lazy, Suspense } from 'react';

const Home = lazy(() => import('./pages/Home'));
const Profile = lazy(() => import('./pages/Profile'));

function App() {
  return (
    <Suspense fallback={<CircularProgress />}>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/profile/:username" element={<Profile />} />
      </Routes>
    </Suspense>
  );
}
```

### メモ化
```typescript
import { memo, useMemo, useCallback } from 'react';

export const PostCard = memo(({ post, onLike }) => {
  const formattedDate = useMemo(
    () => formatDate(post.created_at),
    [post.created_at]
  );

  const handleLike = useCallback(() => {
    onLike(post.id);
  }, [post.id, onLike]);

  return <Card>{/* ... */}</Card>;
});
```

---

## エラーハンドリング

### エラーバウンダリ
```typescript
import { Component, ErrorInfo, ReactNode } from 'react';

type Props = { children: ReactNode };
type State = { hasError: boolean };

class ErrorBoundary extends Component<Props, State> {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <h1>Something went wrong.</h1>;
    }
    return this.props.children;
  }
}
```

---

## テスト

### Vitest + React Testing Library
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { PostCard } from './PostCard';

test('renders post content', () => {
  const post = { id: '1', content: 'Test post', user: { display_name: 'User' } };
  render(<PostCard post={post} onLike={jest.fn()} />);

  expect(screen.getByText('Test post')).toBeInTheDocument();
});

test('calls onLike when like button clicked', () => {
  const onLike = jest.fn();
  const post = { id: '1', content: 'Test', user: {} };
  render(<PostCard post={post} onLike={onLike} />);

  fireEvent.click(screen.getByRole('button', { name: /like/i }));
  expect(onLike).toHaveBeenCalledWith('1');
});
```

---

## 環境変数

```bash
# .env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_FIREBASE_API_KEY=your-api-key
VITE_FIREBASE_AUTH_DOMAIN=your-auth-domain
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-bucket
```

---

## ビルド・デプロイ

### 開発
```bash
npm run dev
```

### ビルド
```bash
npm run build
```

### Firebase Hosting デプロイ
```bash
firebase deploy --only hosting
```

---

## 参照
- React: https://react.dev/
- MUI: https://mui.com/
- React Query: https://tanstack.com/query/
- React Hook Form: https://react-hook-form.com/
- Vite: https://vitejs.dev/
