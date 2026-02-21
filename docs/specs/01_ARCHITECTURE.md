# システムアーキテクチャ

## 全体構成図

```
┌─────────────────────────────────────────────────────────────┐
│                        クライアント                          │
│              (PC / Tablet / Smartphone)                     │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTPS
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Firebase Hosting                          │
│              (React SPA - フロントエンド)                    │
└──────────────────────┬──────────────────────────────────────┘
                       │ REST API (HTTPS)
                       │ JWT Token
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              Render / Cloud Run                              │
│          (Go + Echo - バックエンドAPI)                       │
│                                                              │
│  ┌──────────────────────────────────────────────┐           │
│  │  Echo Router (Middleware)                    │           │
│  │  - CORS                                      │           │
│  │  - Rate Limiting                             │           │
│  │  - JWT Authentication                        │           │
│  │  - Request Logging                           │           │
│  └────────────┬─────────────────────────────────┘           │
│               ▼                                              │
│  ┌──────────────────────────────────────────────┐           │
│  │  Handlers (Controllers)                      │           │
│  │  - User Handler                              │           │
│  │  - Post Handler                              │           │
│  │  - Comment Handler                           │           │
│  │  - Like Handler                              │           │
│  │  - Follow Handler                            │           │
│  │  - Notification Handler                      │           │
│  │  - Report Handler                            │           │
│  └────────────┬─────────────────────────────────┘           │
│               ▼                                              │
│  ┌──────────────────────────────────────────────┐           │
│  │  Services (Business Logic)                   │           │
│  │  - User Service                              │           │
│  │  - Post Service                              │           │
│  │  - Comment Service                           │           │
│  │  - Like Service                              │           │
│  │  - Follow Service                            │           │
│  │  - Notification Service                      │           │
│  │  - Report Service                            │           │
│  │  - Storage Service (Firebase)                │           │
│  └────────────┬─────────────────────────────────┘           │
│               ▼                                              │
│  ┌──────────────────────────────────────────────┐           │
│  │  Repositories (Data Access Layer)            │           │
│  │  - GORM Models                               │           │
│  │  - Database Queries                          │           │
│  └────────────┬─────────────────────────────────┘           │
└───────────────┼──────────────────────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                       │
│                  (メインデータストア)                        │
└─────────────────────────────────────────────────────────────┘

                ┌───────────────────────┐
                │  Firebase Storage     │
                │  (画像/動画ストレージ)│
                └───────────────────────┘
```

## レイヤー構成

### フロントエンド (React + TypeScript)

```
frontend/
├── src/
│   ├── components/          # UIコンポーネント
│   │   ├── common/         # 共通コンポーネント
│   │   ├── post/           # 投稿関連
│   │   ├── user/           # ユーザー関連
│   │   ├── timeline/       # タイムライン
│   │   └── notification/   # 通知
│   ├── pages/              # ページコンポーネント
│   │   ├── Home.tsx
│   │   ├── Profile.tsx
│   │   ├── Login.tsx
│   │   └── Register.tsx
│   ├── hooks/              # カスタムフック
│   │   ├── useAuth.ts
│   │   ├── usePosts.ts
│   │   └── useNotifications.ts
│   ├── api/                # API通信層
│   │   ├── client.ts       # Axios設定
│   │   └── generated/      # openapi-typescript生成
│   ├── store/              # 状態管理 (Context API or Redux)
│   │   ├── auth/
│   │   ├── posts/
│   │   └── notifications/
│   ├── utils/              # ユーティリティ
│   │   ├── validation.ts
│   │   └── formatters.ts
│   ├── theme/              # MUIテーマ
│   │   ├── lightTheme.ts
│   │   └── darkTheme.ts
│   └── App.tsx
├── public/
├── package.json
└── tsconfig.json
```

### バックエンド (Go + Echo)

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── handler/                 # HTTPハンドラー
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   ├── comment_handler.go
│   │   ├── like_handler.go
│   │   ├── follow_handler.go
│   │   ├── notification_handler.go
│   │   └── report_handler.go
│   ├── service/                 # ビジネスロジック
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── comment_service.go
│   │   ├── like_service.go
│   │   ├── follow_service.go
│   │   ├── notification_service.go
│   │   ├── report_service.go
│   │   └── storage_service.go
│   ├── repository/              # データアクセス層
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── comment_repository.go
│   │   ├── like_repository.go
│   │   ├── follow_repository.go
│   │   ├── notification_repository.go
│   │   └── report_repository.go
│   ├── model/                   # GORMモデル
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── comment.go
│   │   ├── like.go
│   │   ├── follow.go
│   │   ├── notification.go
│   │   ├── report.go
│   │   └── base.go              # 共通フィールド
│   ├── middleware/              # ミドルウェア
│   │   ├── auth.go              # JWT認証
│   │   ├── cors.go
│   │   ├── rate_limit.go
│   │   └── logger.go
│   ├── dto/                     # Data Transfer Objects
│   │   ├── request/
│   │   └── response/
│   └── config/                  # 設定
│       ├── database.go
│       ├── firebase.go
│       └── env.go
├── pkg/                         # 外部パッケージ
│   ├── auth/                    # JWT生成・検証
│   ├── validation/              # バリデーション
│   └── errors/                  # エラーハンドリング
├── docs/                        # swaggo生成
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── migrations/                  # DBマイグレーション
├── go.mod
├── go.sum
└── .air.toml                    # Airホットリロード設定
```

## データフロー

### 1. ユーザー登録フロー

```
[クライアント]
    ↓ POST /api/v1/auth/register
    ↓ { email, password, username, display_name }
[Handler]
    ↓ リクエストバリデーション
[Service]
    ↓ メールアドレス・ユーザー名の重複チェック
    ↓ パスワードハッシュ化 (bcrypt)
[Repository]
    ↓ ユーザーをDBに保存
[Service]
    ↓ JWT生成
[Handler]
    ↓ レスポンス返却
    ↓ { token, user }
[クライアント]
```

### 2. ログインフロー

```
[クライアント]
    ↓ POST /api/v1/auth/login
    ↓ { email, password }
[Handler]
    ↓ リクエストバリデーション
[Service]
    ↓ ユーザー検索（メールアドレス）
    ↓ パスワード検証 (bcrypt)
    ↓ JWT生成
[Handler]
    ↓ レスポンス返却
    ↓ { token, user }
[クライアント]
    ↓ トークンをlocalStorageに保存
```

### 3. 投稿作成フロー (Phase 2: 画像あり)

```
[クライアント]
    ↓ 画像選択
    ↓ POST /api/v1/posts
    ↓ Authorization: Bearer <JWT>
    ↓ { content, media_files[] }
[Middleware]
    ↓ JWT検証
    ↓ ユーザー情報取得
[Handler]
    ↓ リクエストバリデーション
    ↓ ファイルサイズ・形式チェック
[Service]
    ↓ Firebase Storageにアップロード
    ↓ URL取得
    ↓ 投稿データ作成
[Repository]
    ↓ 投稿をDBに保存
[Service]
    ↓ フォロワーに通知作成（非同期）
[Handler]
    ↓ レスポンス返却
    ↓ { post }
[クライアント]
```

### 4. タイムライン取得フロー

```
[クライアント]
    ↓ GET /api/v1/timeline?limit=20&offset=0
    ↓ Authorization: Bearer <JWT>
[Middleware]
    ↓ JWT検証
[Handler]
    ↓ ページネーションパラメータ検証
[Service]
    ↓ フォロー中ユーザーの投稿取得
    ↓ 自分の投稿取得
    ↓ 時系列ソート（新しい順）
[Repository]
    ↓ LIMIT/OFFSET クエリ実行
    ↓ いいね数、コメント数も取得
[Handler]
    ↓ レスポンス返却
    ↓ { posts[], has_more, total }
[クライアント]
    ↓ 無限スクロール表示
```

### 5. 通知取得フロー（1分間隔ポーリング）

```
[クライアント]
    ↓ 1分ごとに実行
    ↓ GET /api/v1/notifications?limit=20&offset=0
    ↓ Authorization: Bearer <JWT>
[Middleware]
    ↓ JWT検証
[Handler]
    ↓ ページネーションパラメータ検証
[Service]
    ↓ 未読通知の有無チェック
    ↓ 通知リスト取得
[Repository]
    ↓ LIMIT/OFFSET クエリ実行
[Handler]
    ↓ レスポンス返却
    ↓ { notifications[], unread_count }
[クライアント]
    ↓ 未読バッジ更新
    ↓ 通知リスト表示
```

## 認証フロー

### JWT構造

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": "uuid-string",
    "email": "user@example.com",
    "exp": 1234567890,
    "iat": 1234567890
  },
  "signature": "..."
}
```

### トークン管理
- **保存場所**: localStorage（フロントエンド）
- **有効期限**: 7日間
- **リフレッシュ**: Phase 1では未実装、Phase 2以降で検討
- **送信方法**: `Authorization: Bearer <token>` ヘッダー

## データベース接続

### 接続プール設定
```go
// GORM設定例
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### トランザクション管理
- 複数テーブルへの書き込みはトランザクションで保護
- 例: 投稿削除時（投稿、いいね、コメント、通知を一括削除）

## エラーハンドリング

### HTTPステータスコード
- `200 OK`: 成功
- `201 Created`: リソース作成成功
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証エラー
- `403 Forbidden`: 権限エラー
- `404 Not Found`: リソースが存在しない
- `409 Conflict`: 重複エラー（メール、ユーザー名等）
- `429 Too Many Requests`: レート制限
- `500 Internal Server Error`: サーバーエラー

### エラーレスポンス形式
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

## パフォーマンス最適化

### データベース
- 適切なインデックス設定
- N+1問題の回避（GORM Preload使用）
- ページネーションによるデータ量制限

### API
- レスポンスのgzip圧縮
- 適切なキャッシュヘッダー設定
- レート制限による負荷軽減

### フロントエンド
- コード分割（React.lazy）
- 画像の遅延読み込み（Lazy Loading）
- 仮想スクロール（無限スクロール最適化）
- SWRやReact Queryによるキャッシング

## セキュリティ対策

### バックエンド
- パスワードのbcryptハッシュ化
- JWT署名検証
- CORS設定（許可されたオリジンのみ）
- SQLインジェクション対策（GORM使用）
- XSS対策（入力のサニタイゼーション）
- ファイルアップロード検証
  - ファイル形式チェック（MIME type）
  - ファイルサイズ制限
  - ファイル名のサニタイゼーション

### フロントエンド
- XSS対策（Reactのエスケープ機能）
- CSRF対策（JWTによるステートレス認証）
- 安全なトークン保存（localStorage）
- HTTPSのみの通信

## デプロイ構成

### 開発環境（ローカル）
```
Docker Compose
├── backend (Go + Air)       ポート: 8080
├── db (PostgreSQL)          ポート: 5432
└── frontend (開発サーバー)   ポート: 3000
```

### 本番環境
```
Firebase Hosting (フロントエンド)
    ↓ API呼び出し
Render/Cloud Run (バックエンド)
    ↓ DB接続
PostgreSQL (マネージドDB)

Firebase Storage (画像/動画)
```

## 環境変数

### バックエンド (.env)
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=social_media_db

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=168h  # 7日

# Server
PORT=8080
ENV=development  # development, staging, production

# Firebase
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_STORAGE_BUCKET=your-bucket

# File Upload
MAX_VIDEO_SIZE_MB=20
MAX_IMAGE_SIZE_MB=5
MAX_UPLOAD_FILES=4

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
```

### フロントエンド (.env)
```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_FIREBASE_API_KEY=your-api-key
VITE_FIREBASE_AUTH_DOMAIN=your-auth-domain
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-bucket
```

## 次のステップ
データベーススキーマの詳細設計（02_DATABASE_SCHEMA.md）
