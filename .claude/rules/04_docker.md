# Docker操作ルール

## 🚨 絶対に禁止されているコマンド

### データ削除の危険があるコマンド
```bash
# ❌ 絶対に実行してはいけない
docker compose down -v           # ボリュームを削除（データベース全削除）
docker volume rm postgres_data   # ボリュームを削除（データベース全削除）
docker volume prune             # 未使用ボリューム全削除
```

**理由**: これらのコマンドは PostgreSQL のデータボリュームを削除し、すべてのデータが失われます。

---

## ✅ 安全なDocker操作

### コンテナの起動
```bash
docker compose up -d
```

### コンテナの停止
```bash
docker compose down  # -v フラグは付けない！
```

### コンテナの再起動
```bash
# すべてのコンテナを再起動
docker compose restart

# 特定のコンテナのみ再起動
docker compose restart backend
docker compose restart db
```

### コンテナの再ビルド
```bash
# コードを変更した場合
docker compose up -d --build

# 特定のコンテナのみ再ビルド
docker compose up -d --build backend
```

---

## docker-compose.yml 構成例

```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: sns_db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: social_media_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - sns_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: sns_backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=social_media_db
      - JWT_SECRET=${JWT_SECRET}
      - FIREBASE_PROJECT_ID=${FIREBASE_PROJECT_ID}
      - FIREBASE_STORAGE_BUCKET=${FIREBASE_STORAGE_BUCKET}
    volumes:
      - ./backend:/app
      - /app/tmp  # Airの一時ファイル用
    depends_on:
      db:
        condition: service_healthy
    networks:
      - sns_network

volumes:
  postgres_data:  # データベースデータの永続化

networks:
  sns_network:
    driver: bridge
```

---

## バックエンド Dockerfile (Air使用)

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Air（ホットリロード）のインストール
RUN go install github.com/cosmtrek/air@latest

# 依存関係のコピーとインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# Airで起動
CMD ["air", "-c", ".air.toml"]
```

### .air.toml 設定例

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

---

## よく使うコマンド

### コンテナの状態確認
```bash
docker compose ps
```

### ログの確認
```bash
# すべてのコンテナのログ
docker compose logs

# 特定のコンテナのログ
docker compose logs backend
docker compose logs db

# リアルタイムでログを見る
docker compose logs -f backend
```

### コンテナ内でコマンド実行
```bash
# データベースに接続
docker compose exec db psql -U postgres -d social_media_db

# バックエンドコンテナでbashを起動
docker compose exec backend sh

# マイグレーション実行（例）
docker compose exec backend go run cmd/migrate/main.go
```

### データベースバックアップ
```bash
# バックアップ作成
docker compose exec db pg_dump -U postgres social_media_db > backup.sql

# バックアップから復元
docker compose exec -T db psql -U postgres social_media_db < backup.sql
```

---

## トラブルシューティング

### ポートが既に使用されている
```bash
# ポートを使用しているプロセスを確認
lsof -i :5432
lsof -i :8080

# プロセスを終了
kill -9 <PID>

# または docker-compose.yml のポートを変更
ports:
  - "5433:5432"  # ホスト側のポートを変更
```

### コンテナが起動しない
```bash
# ログを確認
docker compose logs backend
docker compose logs db

# コンテナの状態を確認
docker compose ps

# コンテナを完全に削除して再作成
docker compose down
docker compose up -d
```

### データベース接続エラー
```bash
# データベースが起動しているか確認
docker compose ps db

# データベースのヘルスチェック
docker compose exec db pg_isready -U postgres

# 環境変数を確認
docker compose exec backend env | grep DB_
```

### Airのホットリロードが動かない
```bash
# バックエンドコンテナのログを確認
docker compose logs -f backend

# .air.toml の設定を確認
# ファイルの変更が検知されているか確認

# コンテナを再起動
docker compose restart backend
```

### ボリュームのクリーンアップ（注意！）
**警告**: 以下のコマンドはデータを削除します。必要な場合のみ実行。

```bash
# 未使用のボリュームを確認
docker volume ls

# 特定のボリュームを削除（データベースのデータも削除される）
# ユーザーに必ず確認を取ること
docker volume rm <volume_name>
```

---

## 開発フロー

### 1. 初回セットアップ
```bash
# リポジトリクローン後
cd app

# 環境変数設定
cp .env.example .env
# .env を編集

# コンテナ起動
docker compose up -d

# ログ確認
docker compose logs -f
```

### 2. 日常の開発
```bash
# コンテナ起動（既に起動している場合は不要）
docker compose up -d

# コード編集
# → Airが自動的に再ビルド・再起動

# ログ確認
docker compose logs -f backend
```

### 3. 終了時
```bash
# コンテナ停止（データは保持される）
docker compose down

# 次回起動時
docker compose up -d
```

---

## データベース操作

### psql でデータベースに接続
```bash
docker compose exec db psql -U postgres -d social_media_db
```

### SQL実行
```sql
-- テーブル一覧
\dt

-- テーブル構造確認
\d users

-- データ確認
SELECT * FROM users;

-- 終了
\q
```

---

## 環境変数管理

### .env ファイル
```bash
# Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=social_media_db

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=168h

# Firebase
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_STORAGE_BUCKET=your-bucket

# Server
PORT=8080
ENV=development
```

### .env.example
- `.env.example` をリポジトリにコミット
- `.env` は `.gitignore` に追加
- チームメンバーは `.env.example` をコピーして `.env` を作成

---

## ベストプラクティス

### 1. データ永続化は必須
- ボリュームを使用してデータベースデータを永続化
- `-v` フラグは絶対に使用しない

### 2. ホットリロードの活用
- Air を使用してバックエンドのホットリロード
- コード変更時の再起動を自動化

### 3. ログの確認
- エラーが発生したら必ずログを確認
- `docker compose logs -f` でリアルタイム監視

### 4. 定期的なクリーンアップ
```bash
# 未使用のイメージを削除
docker image prune

# 未使用のコンテナを削除
docker container prune

# ⚠️ ボリュームは削除しない
```

---

## 参照
- Docker Compose: https://docs.docker.com/compose/
- Air: https://github.com/cosmtrek/air
- PostgreSQL Docker: https://hub.docker.com/_/postgres
