# Phase 2: 中優先度機能 Part 1 (Firebase Storage統合)

## 概要
メディア投稿機能（画像/動画）と追加のソーシャル機能の実装

---

## 1. Firebase Storage 準備

### 1.1 Firebase プロジェクトセットアップ
- [ ] Firebaseプロジェクト作成
- [ ] Firebase Storage有効化
- [ ] サービスアカウントキー取得
  - [ ] `serviceAccountKey.json` をバックエンドに配置
  - [ ] `.gitignore` に追加
- [ ] Firebase Storage セキュリティルール設定
  - [ ] 認証済みユーザーのみアップロード可能
  - [ ] ファイルサイズ制限設定
  - [ ] ディレクトリ別ルール設定

### 1.2 環境変数追加
- [ ] バックエンド `.env` に追加
  - [ ] `FIREBASE_PROJECT_ID`
  - [ ] `FIREBASE_STORAGE_BUCKET`
  - [ ] `MAX_VIDEO_SIZE_MB=20`
  - [ ] `MAX_IMAGE_SIZE_MB=5`
  - [ ] `MAX_UPLOAD_FILES=4`
- [ ] フロントエンド `.env` に追加
  - [ ] `VITE_FIREBASE_API_KEY`
  - [ ] `VITE_FIREBASE_AUTH_DOMAIN`
  - [ ] `VITE_FIREBASE_PROJECT_ID`
  - [ ] `VITE_FIREBASE_STORAGE_BUCKET`

---

## 2. バックエンド開発

### 2.1 Firebase Admin SDK 統合
- [ ] Go用Firebase Admin SDKインストール
- [ ] Firebase初期化 (`internal/config/firebase.go`)
- [ ] Storage Service作成 (`internal/service/storage_service.go`)
  - [ ] URL検証機能
  - [ ] ファイルメタデータ取得機能

### 2.2 データベーススキーマ追加
- [ ] `users` テーブル更新
  - [ ] `avatar_url` カラム追加
  - [ ] `header_url` カラム追加
- [ ] `posts` テーブル更新
  - [ ] `visibility` カラム追加 (public, followers, private)
- [ ] `post_media` テーブル作成
  - [ ] GORMモデル作成
  - [ ] マイグレーション実行
- [ ] `retweets` テーブル作成
  - [ ] GORMモデル作成
  - [ ] マイグレーション実行
- [ ] `hashtags` テーブル作成
  - [ ] GORMモデル作成
  - [ ] マイグレーション実行
- [ ] `post_hashtags` テーブル作成
  - [ ] 多対多関連設定
- [ ] `blocks` テーブル作成
  - [ ] GORMモデル作成
  - [ ] マイグレーション実行
- [ ] `reports` テーブル更新
  - [ ] `reportable_type` に `User` 追加

### 2.3 投稿機能更新（メディア対応）
- [ ] Repository更新
  - [ ] `PostMediaRepository` 作成
    - [ ] メディア追加
    - [ ] メディア取得
    - [ ] メディア削除
  - [ ] `PostRepository` 更新
    - [ ] メディア含む投稿取得
- [ ] Service更新
  - [ ] `PostService` 更新
    - [ ] メディア付き投稿作成
    - [ ] メディアURL検証
    - [ ] 公開範囲設定
- [ ] DTO更新
  - [ ] `CreatePostRequest` にメディア配列追加
  - [ ] `PostResponse` にメディア配列追加
- [ ] Handler更新
  - [ ] `POST /api/v1/posts` - メディア対応
- [ ] swaggo アノテーション更新
- [ ] テスト

### 2.4 プロフィール画像機能
- [ ] Service更新
  - [ ] `UserService` 更新
    - [ ] アイコン画像URL更新
    - [ ] ヘッダー画像URL更新
    - [ ] URL検証
- [ ] DTO更新
  - [ ] `UpdateProfileRequest` にアイコン・ヘッダー追加
  - [ ] `UserResponse` に画像URL追加
- [ ] Handler更新
  - [ ] `PATCH /api/v1/users/me` - 画像URL対応
- [ ] swaggo アノテーション更新
- [ ] テスト

### 2.5 リツイート機能
- [ ] Repository
  - [ ] `RetweetRepository` 作成
    - [ ] リツイート追加
    - [ ] リツイート削除
    - [ ] リツイート一覧取得
    - [ ] リツイート状態確認
- [ ] Service
  - [ ] `RetweetService` 作成
    - [ ] リツイートロジック
    - [ ] リツイート解除ロジック
    - [ ] リツイート一覧取得ロジック
- [ ] DTO
  - [ ] `RetweetResponse`
  - [ ] `RetweetListResponse`
- [ ] Handler
  - [ ] `POST /api/v1/posts/:id/retweet` - リツイート
  - [ ] `DELETE /api/v1/posts/:id/retweet` - リツイート解除
  - [ ] `GET /api/v1/posts/:id/retweets` - リツイート一覧
- [ ] 通知統合
  - [ ] リツイート時に通知作成
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.6 ハッシュタグ機能
- [ ] Repository
  - [ ] `HashtagRepository` 作成
    - [ ] ハッシュタグ作成
    - [ ] ハッシュタグ検索
    - [ ] トレンドハッシュタグ取得
- [ ] Service
  - [ ] `HashtagService` 作成
    - [ ] ハッシュタグ抽出ロジック（投稿本文から）
    - [ ] ハッシュタグ付き投稿取得
    - [ ] トレンドハッシュタグ取得
- [ ] DTO
  - [ ] `HashtagResponse`
  - [ ] `TrendingHashtagResponse`
- [ ] Handler
  - [ ] `GET /api/v1/hashtags/trending` - トレンドハッシュタグ
  - [ ] `GET /api/v1/hashtags/:name/posts` - ハッシュタグ付き投稿
- [ ] 投稿作成時の統合
  - [ ] 投稿作成時にハッシュタグ自動抽出
  - [ ] `post_hashtags` テーブルに関連保存
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.7 投稿検索機能
- [ ] Repository
  - [ ] `PostRepository` 更新
    - [ ] 投稿検索クエリ実装
    - [ ] ハッシュタグ検索統合
- [ ] Service
  - [ ] `PostService` 更新
    - [ ] 検索ロジック
    - [ ] 検索結果ランキング
- [ ] DTO
  - [ ] `PostSearchRequest`
  - [ ] `PostSearchResponse`
- [ ] Handler
  - [ ] `GET /api/v1/posts/search?q=keyword` - 投稿検索
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.8 公開範囲設定
- [ ] Service更新
  - [ ] `PostService` 更新
    - [ ] 公開範囲に基づくアクセス制御
    - [ ] タイムライン取得時の公開範囲フィルタ
- [ ] DTO更新
  - [ ] `CreatePostRequest` に `visibility` 追加
  - [ ] `UpdatePostRequest` に `visibility` 追加
- [ ] Handler更新
  - [ ] 投稿取得時の権限チェック
- [ ] テスト

### 2.9 ブロック機能
- [ ] Repository
  - [ ] `BlockRepository` 作成
    - [ ] ブロック追加
    - [ ] ブロック削除
    - [ ] ブロック一覧取得
    - [ ] ブロック状態確認
- [ ] Service
  - [ ] `BlockService` 作成
    - [ ] ブロックロジック
    - [ ] ブロック解除ロジック
    - [ ] ブロック済みユーザーフィルタリング
- [ ] DTO
  - [ ] `BlockResponse`
  - [ ] `BlockListResponse`
- [ ] Handler
  - [ ] `POST /api/v1/users/:username/block` - ブロック
  - [ ] `DELETE /api/v1/users/:username/block` - ブロック解除
  - [ ] `GET /api/v1/blocks` - ブロック一覧
- [ ] 全機能への統合
  - [ ] タイムライン取得時にブロック済みユーザー除外
  - [ ] 投稿取得時にブロック確認
  - [ ] コメント取得時にブロック確認
- [ ] swaggo アノテーション追加
- [ ] テスト

### 2.10 通報機能拡張
- [ ] Repository更新
  - [ ] `ReportRepository` 更新
    - [ ] ユーザー通報対応
    - [ ] コメント通報対応（Phase 3へ移動可能）
- [ ] Service更新
  - [ ] `ReportService` 更新
    - [ ] ユーザー通報ロジック
    - [ ] コメント通報ロジック
- [ ] DTO更新
  - [ ] `reportable_type` に `User`, `Comment` 追加
- [ ] Handler更新
  - [ ] 通報作成時の `reportable_type` 検証
- [ ] テスト

---

## 3. フロントエンド開発

### 3.1 Firebase SDK統合
- [ ] Firebase SDKインストール
  - [ ] `firebase` パッケージ
- [ ] Firebase初期化 (`config/firebase.ts`)
  - [ ] Storage初期化
- [ ] ファイルアップロードユーティリティ作成
  - [ ] `utils/fileUpload.ts`
    - [ ] ファイルバリデーション
    - [ ] Firebase Storageアップロード
    - [ ] URL取得
    - [ ] 画像圧縮（browser-image-compression使用）

### 3.2 メディアアップロード機能
- [ ] MediaUpload コンポーネント作成
  - [ ] ファイル選択UI
  - [ ] プレビュー表示
  - [ ] アップロード進行状況表示
  - [ ] メディア削除
  - [ ] 最大4つ制限
- [ ] PostForm コンポーネント更新
  - [ ] MediaUploadコンポーネント統合
  - [ ] メディア付き投稿作成
  - [ ] 公開範囲選択UI
- [ ] PostCard コンポーネント更新
  - [ ] 画像表示（グリッドレイアウト、最大4枚）
  - [ ] 動画表示（再生コントロール）
  - [ ] 画像クリックで拡大表示
  - [ ] 動画プレーヤー

### 3.3 プロフィール画像機能
- [ ] Settings ページ更新
  - [ ] アイコン画像アップロード
  - [ ] ヘッダー画像アップロード
  - [ ] プレビュー表示
- [ ] Profile ページ更新
  - [ ] アイコン画像表示
  - [ ] ヘッダー画像表示
- [ ] Avatar コンポーネント更新
  - [ ] 画像URL対応
  - [ ] デフォルト画像表示

### 3.4 リツイート機能
- [ ] PostCard コンポーネント更新
  - [ ] リツイートボタン追加
  - [ ] リツイート数表示
  - [ ] リツイート状態表示
- [ ] Timeline更新
  - [ ] リツイート投稿表示
  - [ ] リツイートユーザー表示
- [ ] カスタムフック
  - [ ] `useRetweets` - リツイート機能
- [ ] 通知統合
  - [ ] リツイート通知表示

### 3.5 ハッシュタグ機能
- [ ] ハッシュタグリンク化
  - [ ] 投稿本文内のハッシュタグ自動リンク
  - [ ] クリックでハッシュタグページへ遷移
- [ ] ハッシュタグページ作成 (`pages/Hashtag.tsx`)
  - [ ] ハッシュタグ付き投稿一覧
  - [ ] 投稿数表示
- [ ] トレンドサイドバー作成（PC版）
  - [ ] トレンドハッシュタグ表示
  - [ ] 投稿数表示
- [ ] カスタムフック
  - [ ] `useHashtags` - ハッシュタグ機能

### 3.6 投稿検索機能
- [ ] Search ページ作成 (`pages/Search.tsx`)
  - [ ] 検索フォーム
  - [ ] タブ切り替え（投稿、ユーザー）
  - [ ] 検索結果表示
  - [ ] 無限スクロール
- [ ] Header コンポーネント更新
  - [ ] 検索バー追加
  - [ ] 検索候補表示（オートコンプリート）
- [ ] カスタムフック
  - [ ] `useSearch` - 検索機能

### 3.7 公開範囲設定
- [ ] PostForm コンポーネント更新
  - [ ] 公開範囲選択ドロップダウン
  - [ ] アイコン表示（public, followers, private）
- [ ] PostCard コンポーネント更新
  - [ ] 公開範囲アイコン表示

### 3.8 ブロック機能
- [ ] ブロック機能UI
  - [ ] ブロックボタン
  - [ ] ブロック解除ボタン
  - [ ] 確認ダイアログ
- [ ] ブロック一覧ページ (`pages/Blocks.tsx`)
  - [ ] ブロック済みユーザー一覧
  - [ ] ブロック解除機能
- [ ] カスタムフック
  - [ ] `useBlocks` - ブロック機能
- [ ] タイムライン・検索結果からブロック済みユーザー除外

### 3.9 通報機能拡張
- [ ] 通報ダイアログ更新
  - [ ] ユーザー通報対応
  - [ ] コメント通報対応
- [ ] 通報理由選択更新
  - [ ] 理由選択肢追加

### 3.10 UI/UX改善
- [ ] 画像表示の最適化
  - [ ] 遅延読み込み（Lazy Loading）
  - [ ] プレースホルダー表示
- [ ] 動画再生の最適化
  - [ ] 自動再生オプション
  - [ ] ミュート機能
- [ ] レスポンシブ対応確認
  - [ ] メディア表示
  - [ ] 検索UI
  - [ ] トレンドサイドバー（スマホでは非表示）

---

## 4. 統合テスト

### 4.1 バックエンドテスト
- [ ] メディア付き投稿作成 → 取得
- [ ] プロフィール画像更新
- [ ] リツイート機能
- [ ] ハッシュタグ機能
- [ ] 投稿検索機能
- [ ] 公開範囲設定
- [ ] ブロック機能
- [ ] 通報機能（拡張）

### 4.2 フロントエンド・E2Eテスト
- [ ] 画像アップロード → 投稿作成
- [ ] 動画アップロード → 投稿作成
- [ ] プロフィール画像変更
- [ ] リツイート → タイムライン確認
- [ ] ハッシュタグクリック → 投稿一覧表示
- [ ] 投稿検索
- [ ] ブロック → タイムライン除外確認
- [ ] 公開範囲設定 → 閲覧制限確認

### 4.3 Firebase Storage確認
- [ ] ファイルアップロード動作確認
- [ ] セキュリティルール確認
- [ ] ファイルサイズ制限確認
- [ ] アクセス権限確認

---

## 5. ドキュメント・コードレビュー

- [ ] README.md 更新
  - [ ] Firebase設定手順追加
  - [ ] 環境変数追加説明
- [ ] API仕様書確認（Swagger UI）
- [ ] コードレビュー
- [ ] TODO更新
  - [ ] 完了したタスクをチェック
  - [ ] Phase 3への準備

---

## 完了条件

- [ ] Firebase Storageが正常動作
- [ ] 画像・動画アップロードが正常動作
- [ ] プロフィール画像が表示される
- [ ] リツイート機能が正常動作
- [ ] ハッシュタグ機能が正常動作
- [ ] 投稿検索が正常動作
- [ ] 公開範囲設定が正しく機能
- [ ] ブロック機能が正常動作
- [ ] メディア表示が最適化されている
- [ ] レスポンシブデザイン確認
- [ ] テストが全てパス
- [ ] Phase 3 への準備完了

---

**注意事項:**
- コードを変更するたびに、このTODOリストを確認・更新すること
- 完了したタスクは必ずチェックを入れること
- Firebase Storageのコスト管理に注意
