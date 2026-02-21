# フロントエンド仕様書

## 技術スタック

- **言語**: TypeScript
- **フレームワーク**: React 18+
- **UIライブラリ**: Material-UI (MUI) v5+
- **ルーティング**: React Router v6+
- **状態管理**: Context API + React Query (TanStack Query)
- **フォーム**: React Hook Form + Yup (バリデーション)
- **API通信**: Axios
- **型生成**: openapi-typescript
- **ビルドツール**: Vite
- **ホスティング**: Firebase Hosting

## プロジェクト構造

```
frontend/
├── src/
│   ├── api/                    # API通信
│   │   ├── client.ts          # Axios設定
│   │   ├── generated/         # openapi-typescript生成
│   │   └── endpoints/         # APIエンドポイント定義
│   │       ├── auth.ts
│   │       ├── posts.ts
│   │       ├── users.ts
│   │       └── ...
│   ├── components/            # UIコンポーネント
│   │   ├── common/           # 共通コンポーネント
│   │   │   ├── Button.tsx
│   │   │   ├── TextField.tsx
│   │   │   ├── Avatar.tsx
│   │   │   ├── Layout.tsx
│   │   │   └── InfiniteScroll.tsx
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   └── RegisterForm.tsx
│   │   ├── post/
│   │   │   ├── PostCard.tsx
│   │   │   ├── PostForm.tsx
│   │   │   ├── PostList.tsx
│   │   │   └── MediaUpload.tsx
│   │   ├── comment/
│   │   │   ├── CommentCard.tsx
│   │   │   ├── CommentForm.tsx
│   │   │   └── CommentThread.tsx
│   │   ├── user/
│   │   │   ├── UserCard.tsx
│   │   │   ├── UserProfile.tsx
│   │   │   └── FollowButton.tsx
│   │   ├── notification/
│   │   │   ├── NotificationList.tsx
│   │   │   └── NotificationItem.tsx
│   │   └── timeline/
│   │       └── Timeline.tsx
│   ├── pages/                # ページコンポーネント
│   │   ├── Home.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Profile.tsx
│   │   ├── PostDetail.tsx
│   │   ├── Bookmarks.tsx
│   │   ├── Notifications.tsx
│   │   └── Settings.tsx
│   ├── hooks/                # カスタムフック
│   │   ├── useAuth.ts
│   │   ├── usePosts.ts
│   │   ├── useUsers.ts
│   │   ├── useComments.ts
│   │   ├── useNotifications.ts
│   │   └── useInfiniteScroll.ts
│   ├── contexts/             # Context API
│   │   ├── AuthContext.tsx
│   │   └── ThemeContext.tsx
│   ├── utils/                # ユーティリティ
│   │   ├── validation.ts
│   │   ├── formatters.ts
│   │   ├── constants.ts
│   │   └── storage.ts
│   ├── theme/                # MUIテーマ
│   │   ├── index.ts
│   │   ├── lightTheme.ts
│   │   ├── darkTheme.ts
│   │   └── customTheme.ts
│   ├── types/                # 型定義
│   │   └── index.ts
│   ├── App.tsx
│   ├── main.tsx
│   └── vite-env.d.ts
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .env
```

## 主要ページ

### Phase 1

#### 1. ログイン (`/login`)
- メールアドレス・パスワード入力フォーム
- バリデーション（リアルタイム）
- エラー表示
- 登録ページへのリンク

#### 2. ユーザー登録 (`/register`)
- メールアドレス・パスワード・ユーザー名・表示名入力フォーム
- パスワード強度インジケーター
- ユーザー名の即時重複チェック
- バリデーション
- 登録後、自動ログイン

#### 3. ホーム / タイムライン (`/`)
- フォロー中のユーザー + 自分の投稿表示
- 無限スクロール（20件ずつ読み込み）
- 投稿フォーム（上部固定）
- リアルタイム更新（手動リフレッシュ）

#### 4. プロフィール (`/users/:username`)
- ユーザー情報表示
  - アイコン、ヘッダー画像（Phase 2）
  - 表示名、ユーザー名、自己紹介
  - フォロワー数、フォロー中数
- フォローボタン
- 投稿一覧（無限スクロール）
- タブ切り替え
  - 投稿
  - いいね（Phase 2）

#### 5. 投稿詳細 (`/posts/:id`)
- 投稿内容表示
- コメント一覧（ネスト表示）
- コメント投稿フォーム

#### 6. ブックマーク (`/bookmarks`)
- ブックマークした投稿一覧
- 無限スクロール

#### 7. 通知 (`/notifications`)
- 通知一覧
- 未読バッジ
- 通知タイプ別アイコン
- タップで該当投稿へ遷移

#### 8. 設定 (`/settings`)
- プロフィール編集（表示名、自己紹介）
- テーマ切り替え
- アカウント削除（Phase 3）

### Phase 2 追加ページ

#### 9. ハッシュタグ (`/hashtags/:name`)
- ハッシュタグ付き投稿一覧

#### 10. 検索 (`/search`)
- ユーザー検索
- 投稿検索（ハッシュタグ含む）

---

## コンポーネント設計

### 共通コンポーネント

#### Layout
```tsx
<Layout>
  <Header />
  <Sidebar />
  <Main>{children}</Main>
</Layout>
```
- レスポンシブ対応（PC: 3カラム、スマホ: 1カラム + ドロワー）

#### Header
- ロゴ
- 検索バー（Phase 2）
- 通知アイコン（未読バッジ付き）
- ユーザーメニュー

#### Sidebar
- ナビゲーションリンク
  - ホーム
  - 通知
  - ブックマーク
  - プロフィール
  - 設定
- 投稿ボタン

#### InfiniteScroll
```tsx
<InfiniteScroll
  loadMore={loadMore}
  hasMore={hasMore}
  loader={<CircularProgress />}
>
  {items.map(item => <Item key={item.id} {...item} />)}
</InfiniteScroll>
```

---

### 投稿コンポーネント

#### PostCard
- ユーザー情報（アイコン、名前、ユーザー名）
- 投稿日時（相対時間: "3分前"）
- 投稿本文
- メディア（画像/動画）（Phase 2）
- アクションボタン
  - いいね（いいね数表示）
  - コメント（コメント数表示）
  - リツイート（Phase 2）
  - ブックマーク
  - 共有
  - メニュー（編集、削除、通報）

#### PostForm
- テキストエリア（250文字制限、残り文字数表示）
- メディアアップロードボタン（Phase 2）
- 公開範囲選択（Phase 2）
- 投稿ボタン

---

### コメントコンポーネント

#### CommentCard
- PostCardと同様のUI
- 親コメントへの返信アイコン
- ネスト表示（インデント）

#### CommentThread
- コメントをネスト構造で表示
- 再帰的レンダリング

---

## 状態管理

### AuthContext
```tsx
type AuthContextType = {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
};
```

### React Query
- API通信のキャッシング
- 自動リフレッシュ
- 楽観的更新

**例: 投稿一覧取得**
```tsx
const { data, fetchNextPage, hasNextPage } = useInfiniteQuery({
  queryKey: ['timeline'],
  queryFn: ({ pageParam = 0 }) => api.getTimeline({ offset: pageParam }),
  getNextPageParam: (lastPage, pages) =>
    lastPage.has_more ? pages.length * 20 : undefined,
});
```

---

## テーマ切り替え

### テーマ種類
1. **Light Theme** (デフォルト)
2. **Dark Theme**
3. **Custom Theme 1** (青基調)
4. **Custom Theme 2** (緑基調)

### ThemeContext
```tsx
type ThemeContextType = {
  currentTheme: 'light' | 'dark' | 'custom1' | 'custom2';
  setTheme: (theme: string) => void;
};
```

### テーマ設定の永続化
localStorageに保存

---

## バリデーション

### ユーザー登録
- **email**: 必須、メール形式
- **password**: 必須、8文字以上、英数字含む
- **username**: 必須、3-50文字、英数字とアンダースコアのみ
- **display_name**: 必須、1-100文字

### 投稿
- **content**: 必須、1-250文字

### コメント
- **content**: 必須、1-500文字

---

## レスポンシブデザイン

### ブレークポイント（MUI）
- **xs**: 0px-599px (スマホ)
- **sm**: 600px-959px (タブレット)
- **md**: 960px-1279px (小型PC)
- **lg**: 1280px-1919px (PC)
- **xl**: 1920px以上 (大型PC)

### レイアウト
#### PC (md以上)
```
┌─────────────────────────────────────────┐
│ Header                                  │
├──────┬──────────────────────┬───────────┤
│ Side │ Timeline (投稿一覧)  │ Trending  │
│ bar  │                      │ (Phase 2) │
└──────┴──────────────────────┴───────────┘
```

#### スマホ (xs-sm)
```
┌─────────────────────┐
│ Header              │
├─────────────────────┤
│ Timeline            │
│                     │
│                     │
└─────────────────────┘
Bottom Navigation
```

---

## 通知機能（ポーリング）

### 実装方法
```tsx
useEffect(() => {
  const interval = setInterval(() => {
    queryClient.invalidateQueries(['notifications']);
  }, 60000); // 1分間隔

  return () => clearInterval(interval);
}, []);
```

### 未読バッジ
ヘッダーの通知アイコンに未読件数を表示

---

## エラーハンドリング

### APIエラー
- 401: ログイン画面へリダイレクト
- 403: 権限エラーメッセージ表示
- 404: 404ページ表示
- 500: エラーメッセージ表示

### フォームエラー
- リアルタイムバリデーション
- エラーメッセージをフィールド下に表示

---

## アクセシビリティ

- セマンティックHTML使用
- aria-label設定
- キーボードナビゲーション対応
- スクリーンリーダー対応

---

## パフォーマンス最適化

- コード分割（React.lazy）
- 画像の遅延読み込み
- 仮想スクロール（react-window）
- React.memo によるメモ化
- useMemo, useCallback 活用

---

## 環境変数

```bash
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

### プレビュー
```bash
npm run preview
```

### デプロイ (Firebase Hosting)
```bash
firebase deploy --only hosting
```
