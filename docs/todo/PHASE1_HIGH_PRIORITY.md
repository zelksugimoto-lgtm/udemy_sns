# Phase 1: 高優先度機能 TODO

## 概要
基本的なSNS機能の実装（テキスト投稿のみ、メディアアップロードなし）

---

## 1. プロジェクトセットアップ

### 1.1 環境構築
- [ ] ディレクトリ構造作成
  - [ ] `backend/` ディレクトリ作成
  - [ ] `frontend/` ディレクトリ作成
- [ ] `.env.example` ファイル作成
- [ ] `.gitignore` ファイル作成
- [ ] `README.md` 作成

### 1.2 Docker環境構築
- [ ] `docker-compose.yml` 作成
  - [ ] PostgreSQL コンテナ設定
  - [ ] バックエンド (Go + Air) コンテナ設定
  - [ ] ボリューム設定（DBデータ永続化）
  - [ ] ネットワーク設定
- [ ] Docker動作確認
  - [ ] `docker compose up -d` でコンテナ起動
  - [ ] PostgreSQL接続確認
  - [ ] Air ホットリロード動作確認

---

## 2. バックエンド開発

### 2.1 プロジェクト初期化
- [ ] Go module 初期化 (`go mod init`)
- [ ] 必要なパッケージのインストール
  - [ ] Echo framework
  - [ ] GORM
  - [ ] PostgreSQL driver
  - [ ] JWT library (golang-jwt/jwt)
  - [ ] bcrypt
  - [ ] swaggo/echo-swagger
  - [ ] Air (ホットリロード)
- [ ] ディレクトリ構造作成
  - [ ] `cmd/server/`
  - [ ] `internal/handler/`
  - [ ] `internal/service/`
  - [ ] `internal/repository/`
  - [ ] `internal/model/`
  - [ ] `internal/middleware/`
  - [ ] `internal/dto/`
  - [ ] `internal/config/`
  - [ ] `pkg/auth/`
  - [ ] `pkg/validation/`
  - [ ] `pkg/errors/`
- [ ] `.air.toml` 設定ファイル作成

### 2.2 データベース設計・マイグレーション
- [ ] GORMモデル定義
  - [ ] `BaseModel` (共通フィールド)
  - [ ] `User` モデル
  - [ ] `Post` モデル
  - [ ] `Comment` モデル
  - [ ] `Like` モデル
  - [ ] `Bookmark` モデル
  - [ ] `Follow` モデル
  - [ ] `Notification` モデル
  - [ ] `Report` モデル
- [ ] マイグレーション実行
  - [ ] AutoMigrate実装
  - [ ] 初期マイグレーション実行
  - [ ] インデックス確認

### 2.3 設定・ユーティリティ
- [ ] 環境変数読み込み (`config/env.go`)
- [ ] データベース接続設定 (`config/database.go`)
- [ ] JWT生成・検証ユーティリティ (`pkg/auth/jwt.go`)
- [ ] エラーハンドリング (`pkg/errors/`)
- [ ] バリデーション (`pkg/validation/`)

### 2.4 ミドルウェア
- [ ] CORS ミドルウェア
- [ ] JWT認証ミドルウェア
- [ ] レート制限ミドルウェア
- [ ] ロギングミドルウェア
- [ ] エラーハンドリングミドルウェア

### 2.5 認証機能 (Auth)
- [ ] Repository
  - [ ] `UserRepository` 作成
    - [ ] ユーザー作成
    - [ ] メールアドレスでユーザー検索
    - [ ] ユーザー名でユーザー検索
    - [ ] ユーザーID でユーザー検索
- [ ] Service
  - [ ] `UserService` 作成
    - [ ] ユーザー登録ロジック
    - [ ] パスワードハッシュ化
    - [ ] ログインロジック
    - [ ] パスワード検証
    - [ ] JWT生成
- [ ] DTO
  - [ ] `RegisterRequest`
  - [ ] `LoginRequest`
  - [ ] `AuthResponse`
  - [ ] `UserResponse`
- [ ] Handler
  - [ ] `POST /api/v1/auth/register` - ユーザー登録
  - [ ] `POST /api/v1/auth/login` - ログイン
  - [ ] `GET /api/v1/auth/me` - 現在のユーザー情報取得
- [ ] swaggo アノテーション追加
- [ ] テスト
  - [ ] 登録API テスト
  - [ ] ログインAPI テスト

### 2.6 ユーザー機能 (Users)
- [ ] Repository
  - [ ] ユーザー情報更新
  - [ ] ユーザー名でユーザー検索
  - [ ] ユーザー検索（キーワード）
- [ ] Service
  - [ ] プロフィール取得ロジック
  - [ ] プロフィール更新ロジック
  - [ ] ユーザー検索ロジック
  - [ ] フォロワー数・フォロー中数取得
- [ ] DTO
  - [ ] `UpdateProfileRequest`
  - [ ] `UserProfileResponse`
  - [ ] `UserListResponse`
- [ ] Handler
  - [ ] `GET /api/v1/users/:username` - プロフィール取得
  - [ ] `PATCH /api/v1/users/me` - プロフィール更新
  - [ ] `GET /api/v1/users?q=keyword` - ユーザー検索
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.7 投稿機能 (Posts) - テキストのみ
- [ ] Repository
  - [ ] `PostRepository` 作成
    - [ ] 投稿作成
    - [ ] 投稿取得（ID）
    - [ ] 投稿更新
    - [ ] 投稿削除（論理削除）
    - [ ] タイムライン取得（ページネーション）
    - [ ] ユーザーの投稿一覧取得
    - [ ] いいね数・コメント数取得
- [ ] Service
  - [ ] `PostService` 作成
    - [ ] 投稿作成ロジック
    - [ ] 投稿取得ロジック
    - [ ] 投稿更新ロジック（権限チェック）
    - [ ] 投稿削除ロジック（権限チェック）
    - [ ] タイムライン取得ロジック
    - [ ] ユーザー投稿取得ロジック
- [ ] DTO
  - [ ] `CreatePostRequest`
  - [ ] `UpdatePostRequest`
  - [ ] `PostResponse`
  - [ ] `PostListResponse`
- [ ] Handler
  - [ ] `POST /api/v1/posts` - 投稿作成
  - [ ] `GET /api/v1/posts/:id` - 投稿取得
  - [ ] `PATCH /api/v1/posts/:id` - 投稿更新
  - [ ] `DELETE /api/v1/posts/:id` - 投稿削除
  - [ ] `GET /api/v1/timeline` - タイムライン取得
  - [ ] `GET /api/v1/users/:username/posts` - ユーザー投稿一覧
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.8 コメント機能 (Comments)
- [ ] Repository
  - [ ] `CommentRepository` 作成
    - [ ] コメント作成
    - [ ] コメント取得（投稿ID）
    - [ ] コメント削除（論理削除）
    - [ ] ネストコメント取得
    - [ ] いいね数取得
- [ ] Service
  - [ ] `CommentService` 作成
    - [ ] コメント作成ロジック
    - [ ] コメント取得ロジック
    - [ ] コメント削除ロジック（権限チェック）
- [ ] DTO
  - [ ] `CreateCommentRequest`
  - [ ] `CommentResponse`
  - [ ] `CommentListResponse`
- [ ] Handler
  - [ ] `POST /api/v1/posts/:id/comments` - コメント作成
  - [ ] `GET /api/v1/posts/:id/comments` - コメント取得
  - [ ] `DELETE /api/v1/comments/:id` - コメント削除
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.9 いいね機能 (Likes)
- [ ] Repository
  - [ ] `LikeRepository` 作成
    - [ ] いいね追加
    - [ ] いいね削除
    - [ ] いいね状態確認
- [ ] Service
  - [ ] `LikeService` 作成
    - [ ] 投稿いいねロジック
    - [ ] コメントいいねロジック
    - [ ] いいね解除ロジック
- [ ] DTO
  - [ ] `LikeResponse`
- [ ] Handler
  - [ ] `POST /api/v1/posts/:id/like` - 投稿いいね
  - [ ] `DELETE /api/v1/posts/:id/like` - 投稿いいね解除
  - [ ] `POST /api/v1/comments/:id/like` - コメントいいね
  - [ ] `DELETE /api/v1/comments/:id/like` - コメントいいね解除
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.10 ブックマーク機能 (Bookmarks)
- [ ] Repository
  - [ ] `BookmarkRepository` 作成
    - [ ] ブックマーク追加
    - [ ] ブックマーク削除
    - [ ] ブックマーク一覧取得
- [ ] Service
  - [ ] `BookmarkService` 作成
    - [ ] ブックマーク追加ロジック
    - [ ] ブックマーク削除ロジック
    - [ ] ブックマーク一覧取得ロジック
- [ ] DTO
  - [ ] `BookmarkResponse`
  - [ ] `BookmarkListResponse`
- [ ] Handler
  - [ ] `POST /api/v1/posts/:id/bookmark` - ブックマーク追加
  - [ ] `DELETE /api/v1/posts/:id/bookmark` - ブックマーク削除
  - [ ] `GET /api/v1/bookmarks` - ブックマーク一覧
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.11 フォロー機能 (Follows)
- [ ] Repository
  - [ ] `FollowRepository` 作成
    - [ ] フォロー追加
    - [ ] フォロー削除
    - [ ] フォロワー一覧取得
    - [ ] フォロー中一覧取得
    - [ ] フォロー状態確認
- [ ] Service
  - [ ] `FollowService` 作成
    - [ ] フォローロジック
    - [ ] フォロー解除ロジック
    - [ ] フォロワー一覧取得ロジック
    - [ ] フォロー中一覧取得ロジック
- [ ] DTO
  - [ ] `FollowResponse`
  - [ ] `FollowerListResponse`
- [ ] Handler
  - [ ] `POST /api/v1/users/:username/follow` - フォロー
  - [ ] `DELETE /api/v1/users/:username/follow` - フォロー解除
  - [ ] `GET /api/v1/users/:username/followers` - フォロワー一覧
  - [ ] `GET /api/v1/users/:username/following` - フォロー中一覧
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.12 通知機能 (Notifications)
- [ ] Repository
  - [ ] `NotificationRepository` 作成
    - [ ] 通知作成
    - [ ] 通知一覧取得
    - [ ] 未読通知数取得
    - [ ] 通知既読化
    - [ ] 全通知既読化
- [ ] Service
  - [ ] `NotificationService` 作成
    - [ ] 通知作成ロジック
    - [ ] 通知一覧取得ロジック
    - [ ] 通知既読化ロジック
- [ ] DTO
  - [ ] `NotificationResponse`
  - [ ] `NotificationListResponse`
- [ ] Handler
  - [ ] `GET /api/v1/notifications` - 通知一覧取得
  - [ ] `PATCH /api/v1/notifications/:id/read` - 通知既読化
  - [ ] `POST /api/v1/notifications/read-all` - 全通知既読化
- [ ] 通知生成の統合
  - [ ] いいね時に通知作成
  - [ ] コメント時に通知作成
  - [ ] フォロー時に通知作成
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.13 通報機能 (Reports)
- [ ] Repository
  - [ ] `ReportRepository` 作成
    - [ ] 通報作成
    - [ ] 通報一覧取得（管理者用、Phase 3で使用）
- [ ] Service
  - [ ] `ReportService` 作成
    - [ ] 通報作成ロジック
    - [ ] 重複通報チェック
- [ ] DTO
  - [ ] `CreateReportRequest`
  - [ ] `ReportResponse`
- [ ] Handler
  - [ ] `POST /api/v1/reports` - 通報作成
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.14 OpenAPI生成
- [ ] swaggo でOpenAPI定義生成
  - [ ] `swag init` 実行
  - [ ] `docs/swagger.json` 生成確認
  - [ ] `docs/swagger.yaml` 生成確認
- [ ] Swagger UI 動作確認
  - [ ] `http://localhost:8080/swagger/index.html` アクセス確認

---

## 3. フロントエンド開発

### 3.1 プロジェクト初期化
- [x] Vite + React + TypeScript プロジェクト作成
- [x] 必要なパッケージのインストール
  - [x] React Router
  - [x] Material-UI (MUI)
  - [x] Axios
  - [x] React Query (TanStack Query)
  - [x] React Hook Form
  - [x] Yup
  - [x] openapi-typescript
  - [x] date-fns
  - [x] @mui/icons-material
- [x] ディレクトリ構造作成
- [x] `.env` ファイル作成

### 3.2 型生成・API設定
- [x] openapi-typescript で型生成
  - [x] `backend/docs/openapi.yaml` から型生成
  - [x] `src/api/generated/schema.ts` に配置
- [x] Axios クライアント設定
  - [x] Base URL設定
  - [x] インターセプター設定（JWT自動付与、401ハンドリング）
- [x] APIエンドポイント定義
  - [x] `src/api/endpoints/auth.ts`
  - [x] `src/api/endpoints/users.ts`
  - [x] `src/api/endpoints/posts.ts`
  - [x] `src/api/endpoints/comments.ts`
  - [x] `src/api/endpoints/bookmarks.ts`
  - [x] `src/api/endpoints/notifications.ts`
  - [x] `src/api/endpoints/reports.ts`

### 3.3 テーマ設定
- [x] MUIテーマ作成
  - [x] `theme/lightTheme.ts`
  - [x] `theme/darkTheme.ts`
  - [ ] `theme/customTheme.ts` (2パターン) ※Phase 2で実装
- [x] ThemeContext 作成
  - [x] テーマ切り替え機能
  - [x] localStorageに永続化

### 3.4 認証機能
- [x] AuthContext 作成
  - [x] ユーザー状態管理
  - [x] login 関数
  - [x] register 関数
  - [x] logout 関数
- [x] トークン管理ユーティリティ (`utils/storage.ts`)
- [x] ProtectedRoute コンポーネント（useEffectでリダイレクト実装）
- [x] ログインページ (`pages/Login.tsx`)
  - [x] ログインフォーム
  - [x] バリデーション
  - [x] エラー表示
- [x] ユーザー登録ページ (`pages/Register.tsx`)
  - [x] 登録フォーム
  - [ ] パスワード強度インジケーター ※Phase 2で実装
  - [x] バリデーション
  - [x] エラー表示

### 3.5 レイアウト・共通コンポーネント
- [x] Layout コンポーネント
  - [x] Header
  - [x] Sidebar
  - [x] Main
  - [x] レスポンシブ対応
- [x] Header コンポーネント
  - [x] ロゴ
  - [x] 通知アイコン（未読バッジ対応）
  - [x] ユーザーメニュー
  - [x] テーマ切り替えボタン
- [x] Sidebar コンポーネント
  - [x] ナビゲーションリンク
  - [x] 投稿ボタン
- [x] InfiniteScroll コンポーネント
  - [x] 無限スクロール実装（Intersection Observer使用）
  - [x] ローディング表示
- [x] その他共通コンポーネント
  - [x] Avatar (MUIコンポーネント使用)
  - [x] Button (MUIコンポーネント使用)
  - [x] TextField (MUIコンポーネント使用)

### 3.6 ホーム / タイムライン
- [x] Timeline ページ (`pages/Home.tsx`)
  - [x] 投稿フォーム（上部固定）
  - [x] タイムライン表示
  - [x] 無限スクロール
- [x] PostForm コンポーネント
  - [x] テキストエリア（250文字制限、残り文字数表示）
  - [x] 投稿ボタン
  - [x] バリデーション
- [x] PostCard コンポーネント
  - [x] ユーザー情報表示
  - [x] 投稿日時（相対時間）
  - [x] 投稿本文
  - [x] アクションボタン（いいね、コメント、ブックマーク、共有、メニュー）
  - [x] いいね数・コメント数表示
- [x] PostList コンポーネント（Timeline内に統合）
  - [x] 投稿一覧表示
  - [x] 無限スクロール統合
- [x] カスタムフック
  - [x] React Query の useInfiniteQuery 使用
  - [x] Intersection Observer による無限スクロール

### 3.7 ユーザープロフィール
- [ ] Profile ページ (`pages/Profile.tsx`) ※Phase 2で実装
  - [ ] ユーザー情報表示
  - [ ] フォローボタン
  - [ ] 投稿一覧
  - [ ] タブ切り替え（投稿）
- [ ] UserProfile コンポーネント
  - [ ] プロフィール情報
  - [ ] フォロワー数・フォロー中数
- [ ] FollowButton コンポーネント
  - [ ] フォロー/フォロー解除
  - [ ] 状態管理

### 3.8 投稿詳細
- [x] PostDetail ページ (`pages/PostDetail.tsx`)
  - [x] 投稿内容表示
  - [x] コメント一覧
  - [x] コメント投稿フォーム
  - [x] 戻るボタン
- [x] CommentCard コンポーネント
  - [x] コメント表示
  - [x] いいねボタン
  - [x] 返信ボタン
  - [x] ネスト表示（最大3レベル）
- [x] CommentForm コンポーネント
  - [x] コメント投稿フォーム
  - [x] 返信フォーム
- [x] カスタムフック
  - [x] React Query の useMutation 使用

### 3.9 ブックマーク
- [x] Bookmarks ページ (`pages/Bookmarks.tsx`)
  - [x] ブックマーク一覧表示
  - [x] 通常のリスト表示（無限スクロールはPhase 2で実装）

### 3.10 通知
- [x] Notifications ページ (`pages/Notifications.tsx`)
  - [x] 通知一覧表示
  - [x] 未読バッジ対応
  - [x] 既読化（自動）
  - [x] 通知タイプ別アイコン
  - [x] 通知内容表示
  - [x] タップで該当投稿へ遷移
- [ ] 1分間隔ポーリング実装 ※Phase 2で実装

### 3.11 設定
- [x] Settings ページ (`pages/Settings.tsx`)
  - [x] プロフィール編集フォーム
  - [x] テーマ切り替え
- [x] プロフィール編集機能
  - [x] 表示名、自己紹介編集
  - [x] バリデーション

### 3.12 ユーザー検索
- [ ] ユーザー検索機能 ※Phase 2で実装
  - [ ] 検索フォーム
  - [ ] 検索結果表示
  - [ ] ページネーション
- [ ] UserCard コンポーネント
  - [ ] ユーザー情報表示
  - [ ] フォローボタン

### 3.13 その他機能
- [x] 通報機能UI（メニューに表示、詳細はPhase 2）
  - [ ] 通報ダイアログ
  - [ ] 通報理由選択
  - [ ] 詳細入力
- [x] エラーハンドリング
  - [x] エラーメッセージ表示
  - [x] 401エラー時ログイン画面へリダイレクト
- [x] ローディング表示
  - [x] ページ読み込み中
  - [x] API通信中

---

## 4. 統合テスト

### 4.1 バックエンドテスト
- [ ] ユーザー登録 → ログイン
- [ ] 投稿作成 → 取得 → 更新 → 削除
- [ ] コメント作成 → 取得 → 削除
- [ ] いいね機能
- [ ] ブックマーク機能
- [ ] フォロー機能
- [ ] 通知機能
- [ ] 通報機能

### 4.2 フロントエンド・E2Eテスト
- [ ] ユーザー登録フロー
- [ ] ログインフロー
- [ ] 投稿作成フロー
- [ ] コメント投稿フロー
- [ ] いいね・ブックマーク
- [ ] フォロー・フォロー解除
- [ ] 通知確認
- [ ] プロフィール編集

### 4.3 統合動作確認
- [ ] Docker環境で全機能動作確認
- [ ] レスポンシブデザイン確認（PC、タブレット、スマホ）
- [ ] パフォーマンス確認
- [ ] セキュリティチェック

---

## 5. ドキュメント・コードレビュー

- [ ] README.md 更新
  - [ ] セットアップ手順
  - [ ] 環境変数説明
  - [ ] 起動方法
- [ ] API仕様書確認（Swagger UI）
- [ ] コードレビュー
  - [ ] コーディング規約準拠確認
  - [ ] セキュリティチェック
  - [ ] パフォーマンス確認
- [ ] TODO更新
  - [ ] 完了したタスクをチェック
  - [ ] 次のPhaseへの準備

---

## 完了条件

- [ ] すべてのAPIエンドポイントが正常動作
- [ ] フロントエンドの全ページが正常表示
- [ ] Docker環境で問題なく起動
- [ ] レスポンシブデザインが正しく動作
- [ ] 通知のポーリングが1分間隔で動作
- [ ] 無限スクロールが正常動作
- [ ] Swagger UIでAPI仕様が確認できる
- [ ] テストが全てパス
- [ ] Phase 2 への準備完了

---

**注意事項:**
- コードを変更するたびに、このTODOリストを確認・更新すること
- 完了したタスクは必ずチェックを入れること
- 問題が発生した場合は、TODOに詳細を記録すること
