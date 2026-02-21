# Phase 3: 中優先度機能 Part 2 (管理・分析機能)

## 概要
管理者機能とアナリティクス機能の実装

---

## 1. バックエンド開発

### 1.1 データベーススキーマ追加
- [ ] `admin_actions` テーブル作成
  - [ ] GORMモデル作成
  - [ ] マイグレーション実行
- [ ] `post_analytics` テーブル作成
  - [ ] GORMモデル作成
  - [ ] マイグレーション実行

### 1.2 管理者機能

#### 管理者権限チェック
- [ ] Middleware作成
  - [ ] `AdminMiddleware` - is_admin フラグ確認
  - [ ] 管理者専用エンドポイントに適用

#### 通報管理
- [ ] Repository更新
  - [ ] `ReportRepository` 更新
    - [ ] 通報一覧取得（ステータス別）
    - [ ] 通報詳細取得
    - [ ] 通報ステータス更新
- [ ] Service更新
  - [ ] `ReportService` 更新
    - [ ] 通報一覧取得ロジック
    - [ ] 通報対応ロジック
    - [ ] 通報却下ロジック
- [ ] DTO作成
  - [ ] `ReportListResponse` (管理者用)
  - [ ] `UpdateReportStatusRequest`
- [ ] Handler作成
  - [ ] `GET /api/v1/admin/reports` - 通報一覧（管理者専用）
  - [ ] `GET /api/v1/admin/reports/:id` - 通報詳細
  - [ ] `PATCH /api/v1/admin/reports/:id/status` - ステータス更新
- [ ] swaggo アノテーション追加
- [ ] テスト

#### 投稿管理
- [ ] Repository作成
  - [ ] `AdminRepository` 作成
    - [ ] 投稿強制削除（論理削除）
    - [ ] 投稿完全削除（物理削除、必要な場合）
- [ ] Service作成
  - [ ] `AdminService` 作成
    - [ ] 投稿削除ロジック
    - [ ] 削除理由記録
- [ ] DTO作成
  - [ ] `AdminDeletePostRequest`
  - [ ] `AdminActionResponse`
- [ ] Handler作成
  - [ ] `DELETE /api/v1/admin/posts/:id` - 投稿削除
- [ ] `admin_actions` テーブルへの記録
  - [ ] アクション履歴保存
- [ ] swaggo アノテーション追加
- [ ] テスト

#### ユーザー管理
- [ ] Repository更新
  - [ ] `AdminRepository` 更新
    - [ ] ユーザー停止
    - [ ] ユーザー停止解除
    - [ ] ユーザー削除（論理削除）
- [ ] Service更新
  - [ ] `AdminService` 更新
    - [ ] ユーザー停止ロジック
    - [ ] 停止解除ロジック
    - [ ] ユーザー削除ロジック
- [ ] DTO更新
  - [ ] `AdminSuspendUserRequest`
  - [ ] `AdminDeleteUserRequest`
- [ ] Handler更新
  - [ ] `POST /api/v1/admin/users/:id/suspend` - ユーザー停止
  - [ ] `POST /api/v1/admin/users/:id/unsuspend` - 停止解除
  - [ ] `DELETE /api/v1/admin/users/:id` - ユーザー削除
- [ ] `admin_actions` テーブルへの記録
- [ ] swaggo アノテーション追加
- [ ] テスト

#### 管理者ダッシュボード
- [ ] Repository作成
  - [ ] 統計情報取得
    - [ ] ユーザー数
    - [ ] 投稿数
    - [ ] 通報数（ステータス別）
    - [ ] アクティブユーザー数
- [ ] Service作成
  - [ ] ダッシュボード統計ロジック
- [ ] DTO作成
  - [ ] `AdminDashboardResponse`
- [ ] Handler作成
  - [ ] `GET /api/v1/admin/dashboard` - ダッシュボード統計
- [ ] swaggo アノテーション追加
- [ ] テスト

### 1.3 アナリティクス機能

#### 投稿閲覧数トラッキング
- [ ] Repository作成
  - [ ] `PostAnalyticsRepository` 作成
    - [ ] 閲覧数記録
    - [ ] 閲覧数取得
    - [ ] ユニークユーザー数取得
- [ ] Service作成
  - [ ] `AnalyticsService` 作成
    - [ ] 閲覧数記録ロジック
    - [ ] 重複閲覧判定（同一ユーザー、一定期間内）
- [ ] DTO作成
  - [ ] `PostAnalyticsResponse`
- [ ] Handler更新
  - [ ] `GET /api/v1/posts/:id` - 閲覧時に閲覧数記録
  - [ ] `GET /api/v1/posts/:id/analytics` - 投稿分析データ取得
- [ ] swaggo アノテーション追加
- [ ] テスト

#### いいね数推移
- [ ] Repository更新
  - [ ] `PostAnalyticsRepository` 更新
    - [ ] 日別いいね数集計
    - [ ] 期間指定いいね数推移取得
- [ ] Service更新
  - [ ] `AnalyticsService` 更新
    - [ ] いいね数推移ロジック
- [ ] DTO更新
  - [ ] `LikeTrendResponse`
- [ ] Handler更新
  - [ ] `GET /api/v1/posts/:id/analytics/likes` - いいね数推移
- [ ] swaggo アノテーション追加
- [ ] テスト

#### ユーザー分析
- [ ] Repository作成
  - [ ] `UserAnalyticsRepository` 作成
    - [ ] フォロワー数推移
    - [ ] 投稿数推移
    - [ ] エンゲージメント率計算
- [ ] Service更新
  - [ ] `AnalyticsService` 更新
    - [ ] ユーザー分析ロジック
- [ ] DTO作成
  - [ ] `UserAnalyticsResponse`
- [ ] Handler作成
  - [ ] `GET /api/v1/users/:username/analytics` - ユーザー分析
- [ ] swaggo アノテーション追加
- [ ] テスト

---

## 2. フロントエンド開発

### 2.1 管理者機能

#### 管理者権限チェック
- [ ] AuthContext 更新
  - [ ] `isAdmin` フラグ追加
  - [ ] 管理者専用ルート保護

#### 管理者ダッシュボード
- [ ] AdminDashboard ページ作成 (`pages/admin/Dashboard.tsx`)
  - [ ] 統計情報表示
    - [ ] ユーザー数
    - [ ] 投稿数
    - [ ] 通報数（ステータス別）
  - [ ] グラフ表示（Chart.js 等使用）
  - [ ] クイックアクションリンク
- [ ] カスタムフック
  - [ ] `useAdminDashboard` - ダッシュボード統計取得

#### 通報管理
- [ ] Reports ページ作成 (`pages/admin/Reports.tsx`)
  - [ ] 通報一覧表示
  - [ ] ステータス別フィルタ
  - [ ] 通報詳細モーダル
  - [ ] ステータス更新機能
    - [ ] 確認中
    - [ ] 解決済み
    - [ ] 却下
- [ ] ReportDetail コンポーネント作成
  - [ ] 通報内容表示
  - [ ] 通報対象表示（投稿、コメント、ユーザー）
  - [ ] アクションボタン
    - [ ] 対象削除
    - [ ] ユーザー停止
    - [ ] 却下
- [ ] カスタムフック
  - [ ] `useAdminReports` - 通報管理機能

#### 投稿管理
- [ ] PostManagement ページ作成 (`pages/admin/Posts.tsx`)
  - [ ] 投稿一覧表示
  - [ ] 検索機能
  - [ ] 削除機能
  - [ ] 削除理由入力
- [ ] カスタムフック
  - [ ] `useAdminPosts` - 投稿管理機能

#### ユーザー管理
- [ ] UserManagement ページ作成 (`pages/admin/Users.tsx`)
  - [ ] ユーザー一覧表示
  - [ ] 検索機能
  - [ ] ユーザー詳細モーダル
  - [ ] アクション
    - [ ] ユーザー停止
    - [ ] 停止解除
    - [ ] ユーザー削除
- [ ] カスタムフック
  - [ ] `useAdminUsers` - ユーザー管理機能

#### ナビゲーション
- [ ] Sidebar 更新
  - [ ] 管理者メニュー追加（管理者のみ表示）
    - [ ] ダッシュボード
    - [ ] 通報管理
    - [ ] 投稿管理
    - [ ] ユーザー管理

### 2.2 アナリティクス機能

#### 投稿分析
- [ ] PostAnalytics コンポーネント作成
  - [ ] 閲覧数表示
  - [ ] いいね数推移グラフ
  - [ ] コメント数推移グラフ
  - [ ] エンゲージメント率表示
- [ ] PostDetail ページ更新
  - [ ] 自分の投稿の場合、分析データ表示
  - [ ] 分析タブ追加
- [ ] カスタムフック
  - [ ] `usePostAnalytics` - 投稿分析データ取得

#### ユーザー分析
- [ ] UserAnalytics コンポーネント作成
  - [ ] フォロワー数推移グラフ
  - [ ] 投稿数推移グラフ
  - [ ] エンゲージメント率推移
- [ ] Profile ページ更新
  - [ ] 自分のプロフィールの場合、分析データ表示
  - [ ] 分析タブ追加
- [ ] カスタムフック
  - [ ] `useUserAnalytics` - ユーザー分析データ取得

#### グラフライブラリ
- [ ] Chart.jsまたはRecharts インストール
- [ ] グラフコンポーネント作成
  - [ ] LineChart (推移グラフ)
  - [ ] BarChart (棒グラフ)
  - [ ] PieChart (円グラフ)

---

## 3. 統合テスト

### 3.1 バックエンドテスト
- [ ] 管理者機能
  - [ ] 通報管理
  - [ ] 投稿削除
  - [ ] ユーザー停止・削除
  - [ ] ダッシュボード統計
- [ ] アナリティクス機能
  - [ ] 閲覧数記録・取得
  - [ ] いいね数推移
  - [ ] ユーザー分析

### 3.2 フロントエンド・E2Eテスト
- [ ] 管理者ログイン → ダッシュボード表示
- [ ] 通報管理フロー
  - [ ] 通報一覧表示
  - [ ] 通報対応（削除、停止、却下）
- [ ] 投稿管理フロー
  - [ ] 投稿検索
  - [ ] 投稿削除
- [ ] ユーザー管理フロー
  - [ ] ユーザー検索
  - [ ] ユーザー停止
- [ ] アナリティクス表示
  - [ ] 投稿分析グラフ表示
  - [ ] ユーザー分析グラフ表示

### 3.3 権限テスト
- [ ] 一般ユーザーが管理者エンドポイントにアクセス → 403エラー
- [ ] 管理者が全機能にアクセス可能

---

## 4. ドキュメント・コードレビュー

- [ ] README.md 更新
  - [ ] 管理者機能説明追加
  - [ ] アナリティクス機能説明追加
- [ ] API仕様書確認（Swagger UI）
- [ ] 管理者マニュアル作成
  - [ ] 通報対応フロー
  - [ ] ユーザー停止基準
- [ ] コードレビュー
- [ ] TODO更新
  - [ ] 完了したタスクをチェック

---

## 5. デプロイ準備

### 5.1 本番環境設定
- [ ] 環境変数設定（本番）
  - [ ] Render/Cloud Run
  - [ ] Firebase Hosting
- [ ] データベースマイグレーション（本番）
- [ ] Firebase Storage設定（本番）

### 5.2 CI/CD設定
- [ ] GitHub Actions 設定
  - [ ] テスト自動実行
  - [ ] ビルド自動化
  - [ ] デプロイ自動化
- [ ] デプロイスクリプト作成

### 5.3 モニタリング・ロギング
- [ ] ログ設定
  - [ ] エラーログ
  - [ ] アクセスログ
  - [ ] 管理者アクションログ
- [ ] モニタリング設定
  - [ ] パフォーマンス監視
  - [ ] エラー監視

---

## 完了条件

- [ ] 管理者機能が正常動作
  - [ ] 通報管理
  - [ ] 投稿削除
  - [ ] ユーザー停止・削除
  - [ ] ダッシュボード表示
- [ ] アナリティクス機能が正常動作
  - [ ] 閲覧数トラッキング
  - [ ] いいね数推移グラフ
  - [ ] ユーザー分析グラフ
- [ ] 権限制御が正しく機能
- [ ] グラフが正しく表示
- [ ] 本番環境へのデプロイ準備完了
- [ ] テストが全てパス

---

**注意事項:**
- コードを変更するたびに、このTODOリストを確認・更新すること
- 完了したタスクは必ずチェックを入れること
- 管理者機能は慎重にテストすること（誤操作によるデータ削除を防ぐ）
- アナリティクスデータはパフォーマンスに影響する可能性があるため、最適化を意識すること
