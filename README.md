# SNS Application

Twitter風のSNSアプリケーション - テキスト、画像、動画を投稿し、いいね、コメント、フォローなどの機能を持つソーシャルメディアプラットフォーム

## 📋 プロジェクト概要

このプロジェクトは、モダンな技術スタックを使用して構築されたフルスタックSNSアプリケーションです。

### 主な機能

#### Phase 1 (優先度: 高) - 基本SNS機能
- ✅ ユーザー認証（メール+パスワード、JWT）
- ✅ テキスト投稿（最大250文字）
- ✅ いいね機能（投稿・コメント）
- ✅ コメント機能（ネスト対応）
- ✅ ブックマーク機能
- ✅ フォロー/フォロワー機能
- ✅ タイムライン（時系列、無限スクロール）
- ✅ 通知機能（1分間隔ポーリング）
- ✅ ユーザー検索
- ✅ 投稿通報機能

#### Phase 2 (優先度: 中) - メディア機能
- 🔄 画像/動画アップロード（最大4枚、20MB制限）
- 🔄 プロフィール画像・ヘッダー画像
- 🔄 リツイート機能
- 🔄 ハッシュタグ機能
- 🔄 投稿検索
- 🔄 公開範囲設定（全体公開、フォロワーのみ、自分のみ）
- 🔄 ブロック機能

#### Phase 3 (優先度: 中) - 管理・分析
- 📝 管理者機能（通報管理、投稿削除、ユーザー停止）
- 📝 アナリティクス（閲覧数、いいね数推移）

---

## 🛠 技術スタック

### フロントエンド
- **言語**: TypeScript
- **フレームワーク**: React 18+
- **UIライブラリ**: Material-UI (MUI) v5+
- **状態管理**: Context API + React Query (TanStack Query)
- **ルーティング**: React Router v6+
- **フォーム**: React Hook Form + Yup
- **ビルドツール**: Vite
- **型生成**: openapi-typescript

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Echo
- **ORM**: GORM
- **OpenAPI**: swaggo/echo-swagger
- **ホットリロード**: Air
- **認証**: JWT (golang-jwt/jwt)

### インフラ
- **データベース**: PostgreSQL 15
- **ファイルストレージ**: Firebase Storage
- **開発環境**: Docker + Docker Compose
- **デプロイ**:
  - バックエンド: Render / Google Cloud Run
  - フロントエンド: Firebase Hosting

---

## 📁 プロジェクト構造

```
/Users/sugimoto/Desktop/udemy_pj/claudecode/app/
├── frontend/              # React + TypeScript フロントエンド
│   ├── src/
│   │   ├── api/          # API通信・型定義
│   │   ├── components/   # UIコンポーネント
│   │   ├── pages/        # ページコンポーネント
│   │   ├── hooks/        # カスタムフック
│   │   ├── contexts/     # Context API
│   │   ├── theme/        # MUIテーマ
│   │   └── utils/        # ユーティリティ
│   ├── package.json
│   └── vite.config.ts
├── backend/               # Go + Echo バックエンド
│   ├── cmd/              # エントリーポイント
│   ├── internal/         # 内部パッケージ
│   │   ├── handler/      # HTTPハンドラー
│   │   ├── service/      # ビジネスロジック
│   │   ├── repository/   # データアクセス
│   │   ├── model/        # GORMモデル
│   │   ├── middleware/   # ミドルウェア
│   │   └── config/       # 設定
│   ├── pkg/              # 外部パッケージ
│   ├── docs/             # swaggo生成OpenAPI定義
│   ├── go.mod
│   └── .air.toml
├── docs/
│   ├── specs/            # 仕様書
│   │   ├── 00_OVERVIEW.md
│   │   ├── 01_ARCHITECTURE.md
│   │   ├── 02_DATABASE_SCHEMA.md
│   │   ├── 03_API_SPECIFICATION.md
│   │   ├── 04_FRONTEND_SPECIFICATION.md
│   │   ├── 05_AUTHENTICATION.md
│   │   └── 06_FILE_UPLOAD.md
│   └── todo/             # TODOリスト
│       ├── PHASE1_HIGH_PRIORITY.md
│       ├── PHASE2_MEDIUM_PRIORITY_PART1.md
│       └── PHASE3_MEDIUM_PRIORITY_PART2.md
├── .claude/
│   ├── CLAUDE.md         # プロジェクト全体ルール
│   └── rules/            # カテゴリ別開発ルール
│       ├── 01_database.md
│       ├── 02_api_development.md
│       ├── 03_frontend_development.md
│       ├── 04_docker.md
│       └── 05_todo_management.md
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## 🚀 セットアップ

### 前提条件
- Docker & Docker Compose
- Node.js 18+ (フロントエンド開発用)
- Go 1.21+ (ローカルでバックエンド開発する場合)

### 1. リポジトリクローン
```bash
cd /Users/sugimoto/Desktop/udemy_pj/claudecode/app
```

### 2. 環境変数設定
```bash
# ルートディレクトリに .env ファイルを作成
cp .env.example .env

# .env を編集（必要な値を設定）
# - JWT_SECRET
# - FIREBASE_SERVICE_ACCOUNT_KEY（画像アップロード機能を使う場合）
# - FIREBASE_PROJECT_ID
# - FIREBASE_STORAGE_BUCKET
# - その他の設定
```

**画像アップロード機能を使用する場合**:
Firebase Storageの設定が必要です。詳細は [Firebase セットアップガイド](docs/FIREBASE_SETUP.md) を参照してください。

### 3. Docker環境起動
```bash
# コンテナをバックグラウンドで起動
docker compose up -d

# ログを確認
docker compose logs -f
```

### 4. アクセス確認
- **バックエンドAPI**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **PostgreSQL**: localhost:5432
- **フロントエンド** (別途起動): http://localhost:3000

---

## 🔧 開発

### バックエンド開発

#### Airホットリロード
Dockerコンテナ内でAirが動作し、コード変更時に自動的に再ビルド・再起動します。

```bash
# ログでホットリロードを確認
docker compose logs -f backend
```

#### OpenAPI定義生成
```bash
# コンテナ内でswag initを実行
docker compose exec backend swag init -g cmd/server/main.go -o docs

# 生成されたSwagger UIにアクセス
open http://localhost:8080/swagger/index.html
```

#### データベース操作
```bash
# psqlでデータベースに接続
docker compose exec db psql -U postgres -d social_media_db

# マイグレーション実行（アプリ起動時に自動実行される）
```

### フロントエンド開発

#### セットアップ
```bash
cd frontend

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev
```

#### 型生成
```bash
# バックエンドのOpenAPI定義から型を生成
npm run generate-types

# または手動で
npx openapi-typescript http://localhost:8080/swagger/doc.json -o src/api/generated/schema.ts
```

#### ビルド
```bash
# プロダクションビルド
npm run build

# プレビュー
npm run preview
```

---

## 📚 ドキュメント

### 仕様書
開発前に必ず該当する仕様書を確認してください：

- **[プロジェクト概要](docs/specs/00_OVERVIEW.md)**: 全体像、Phase構成
- **[システムアーキテクチャ](docs/specs/01_ARCHITECTURE.md)**: レイヤー構成、データフロー
- **[データベース設計](docs/specs/02_DATABASE_SCHEMA.md)**: テーブル定義、ER図
- **[API仕様](docs/specs/03_API_SPECIFICATION.md)**: エンドポイント一覧
- **[フロントエンド仕様](docs/specs/04_FRONTEND_SPECIFICATION.md)**: コンポーネント設計
- **[認証フロー](docs/specs/05_AUTHENTICATION.md)**: JWT認証の詳細
- **[ファイルアップロード](docs/specs/06_FILE_UPLOAD.md)**: Firebase Storage統合

### TODOリスト
**重要**: 開発はTODOリストに従って進めてください。

- **[Phase 1 TODO](docs/todo/PHASE1_HIGH_PRIORITY.md)**: 基本SNS機能
- **[Phase 2 TODO](docs/todo/PHASE2_MEDIUM_PRIORITY_PART1.md)**: メディア機能
- **[Phase 3 TODO](docs/todo/PHASE3_MEDIUM_PRIORITY_PART2.md)**: 管理・分析機能

### 開発ルール
- **[全体ルール](.claude/CLAUDE.md)**: プロジェクト全体の重要ルール
- **[データベース操作](.claude/rules/01_database.md)**: DB操作のルール
- **[API開発](.claude/rules/02_api_development.md)**: API開発のルール
- **[フロントエンド開発](.claude/rules/03_frontend_development.md)**: フロントエンド開発のルール
- **[Docker操作](.claude/rules/04_docker.md)**: Docker操作のルール
- **[TODO管理](.claude/rules/05_todo_management.md)**: TODO管理のルール

---

## 🔒 重要な注意事項

### ⚠️ データ保護
**絶対に実行してはいけないコマンド**:
```bash
docker compose down -v  # ❌ ボリューム削除（データベース全削除）
docker volume rm postgres_data  # ❌ ボリューム削除
```

**安全なコマンド**:
```bash
docker compose restart  # ✅ 再起動（データ保持）
docker compose down && docker compose up -d  # ✅ 停止・起動（データ保持）
```

### 🔐 セキュリティ
- `.env` ファイルは `.gitignore` に追加済み（コミット禁止）
- パスワードは必ずbcryptでハッシュ化
- JWTシークレットは本番環境で強力なものを使用
- Firebase サービスアカウントキーは `.gitignore` に追加

---

## 🧪 テスト

### テスト環境について
テスト環境は開発環境とは完全に独立しています：
- テスト用PostgreSQL（db_test）: ポート5433
- テスト用APIサーバー（api_test）: ポート8081
- テスト用DBはtmpfs（メモリ）を使用し、開発環境のデータに一切影響を与えません

### Makefile経由でのテスト実行（推奨）

```bash
# バックエンド単体テスト
make test-backend

# フロントエンドE2Eテスト
make test-e2e

# 全テスト実行
make test

# テスト環境の手動操作
make test-setup     # テスト環境を起動
make test-teardown  # テスト環境を停止・削除
```

### バックエンドテスト（直接実行）
```bash
# テスト実行
cd backend
go test ./...

# カバレッジ付き
go test -cover ./...
```

### フロントエンドE2Eテスト（直接実行）
```bash
cd frontend

# E2Eテスト（テスト用APIサーバーが起動している必要があります）
npm run test:e2e

# UIモード
npm run test:e2e:ui

# デバッグモード
npm run test:e2e:debug
```

---

## 🔄 CI/CD (GitHub Actions)

### 自動実行されるテスト

`main` ブランチへの `push` または `Pull Request` 時に以下が自動実行されます：

#### 1. バックエンド単体テスト
- PostgreSQL 15（サービスコンテナ）
- Go 1.23
- `go test -v -cover ./...`

#### 2. フロントエンドビルド確認
- Node.js 20
- `npm run build`

#### 3. E2Eテスト
- PostgreSQL 15（サービスコンテナ）
- バックエンドAPIサーバー（ポート8081）
- フロントエンド開発サーバー（ポート3001）
- Playwright with Chromium
- `npx playwright test`

### CI実行ルール

**必須ルール:**
1. **mainへのマージ前にCIが全てグリーンであること**
2. **push前に必ずローカルでテストを実行すること**
   ```bash
   make test  # 全テスト実行
   ```

### テスト失敗時のデバッグ

**Playwrightテスト失敗時:**
1. GitHub Actionsの「Artifacts」タブからレポートとスクリーンショットをダウンロード
2. `playwright-report/index.html` をブラウザで開く
3. 失敗したテストケースのスクリーンショットとトレースを確認

**バックエンドテスト失敗時:**
1. GitHub Actionsのログで失敗したテストケースを特定
2. ローカルで該当テストを実行: `make test-backend`

詳細は [CI/CDルール](.claude/CLAUDE.md#-cicd継続的インテグレーション) を参照してください。

---

## 📦 デプロイ

### バックエンド (Render / Cloud Run)
```bash
# Render
git push origin main  # 自動デプロイ

# Cloud Run
gcloud builds submit --tag gcr.io/PROJECT_ID/backend
gcloud run deploy backend --image gcr.io/PROJECT_ID/backend
```

### フロントエンド (Firebase Hosting)
```bash
cd frontend

# ビルド
npm run build

# デプロイ
firebase deploy --only hosting
```

---

## 🐛 トラブルシューティング

### コンテナが起動しない
```bash
# ログを確認
docker compose logs backend
docker compose logs db

# コンテナの状態を確認
docker compose ps

# 再起動
docker compose restart
```

### データベース接続エラー
```bash
# データベースが起動しているか確認
docker compose ps db

# 接続テスト
docker compose exec db pg_isready -U postgres
```

### Airのホットリロードが動かない
```bash
# ログを確認
docker compose logs -f backend

# コンテナ再起動
docker compose restart backend
```

---

## 📞 サポート・質問

問題が発生した場合：
1. 該当する `.claude/rules/` のルールファイルを確認
2. `docs/specs/` の仕様書を確認
3. `docs/todo/` のTODOリストを確認
4. エラーログを詳細に記録
5. GitHubのIssueを作成

---

## 📄 ライセンス

このプロジェクトは教育目的で作成されています。

---

## 🎯 開発の進め方

1. **Phase 1から順番に開発**
2. **TODOリストに従って実装**
3. **仕様書を常に参照**
4. **コード変更後はTODO更新**
5. **テストを必ず実施**

詳細は [TODO管理ルール](.claude/rules/05_todo_management.md) を参照してください。

---

**Happy Coding! 🚀**
