# テスト実行ガイド

このドキュメントでは、SNSアプリケーションのテスト環境のセットアップとテスト実行手順を説明します。

---

## 📋 目次

1. [テスト環境の概要](#テスト環境の概要)
2. [前提条件](#前提条件)
3. [テストコマンド一覧](#テストコマンド一覧)
4. [バックエンドテストの実行](#バックエンドテストの実行)
5. [フロントエンドE2Eテストの実行](#フロントエンドe2eテストの実行)
6. [全テストの実行](#全テストの実行)
7. [トラブルシューティング](#トラブルシューティング)

---

## テスト環境の概要

### 環境構成

**テスト環境は開発環境とは完全に独立しています:**

| 項目 | 開発環境 | テスト環境 |
|------|---------|-----------|
| データベース | db (ポート5432) | db_test (ポート5433) |
| APIサーバー | api (ポート8080) | api_test (ポート8081) |
| データ永続化 | Docker Volumeで永続化 | tmpfs（メモリ）で揮発性 |
| プロファイル | dev | test |

**重要:** テスト環境のデータベースはメモリ上に構築されるため、開発環境のデータに**一切影響を与えません**。

---

## 前提条件

### 必要なツール

- Docker Desktop（起動済み）
- Make
- Node.js（フロントエンドE2Eテスト用）

### 初回セットアップ

#### 1. Goモジュールの依存関係をインストール

```bash
cd backend
docker compose --profile test run --rm api_test go mod download
```

#### 2. フロントエンドの依存関係をインストール

```bash
cd frontend
npm install
npx playwright install chromium
```

**注意:** Playwrightブラウザがインストールされていない場合、E2Eテストは失敗します。初回は上記のコマンドを実行してください。

---

## テストコマンド一覧

すべてのテストコマンドはプロジェクトルート（`/Users/sugimoto/Desktop/udemy_pj/app/`）で実行します。

### 開発環境の操作

```bash
make dev-up         # 開発環境を起動
make dev-down       # 開発環境を停止
make dev-logs       # 開発環境のログを表示
```

### テスト実行

```bash
make test-backend   # バックエンドの単体テストを実行
make test-e2e       # フロントエンドのE2Eテストを実行
make test           # 全テストを実行（backend + e2e）
```

### テスト環境の手動操作

```bash
make test-setup     # テスト環境を起動
make test-teardown  # テスト環境を停止・削除
```

### ヘルプ表示

```bash
make help           # 全コマンドの一覧を表示
```

---

## バックエンドテストの実行

### 基本的な実行方法

```bash
# プロジェクトルートで実行
make test-backend
```

このコマンドは以下の処理を自動で行います:

1. テスト環境を起動（db_test、api_test）
2. データベースマイグレーションを実行
3. 単体テストを実行（`-parallel 2`で並列実行）
4. テスト環境をクリーンアップ

### テスト結果の見方

```
=== RUN   TestAuthService_Register
=== RUN   TestAuthService_Register/ユーザー登録成功
=== RUN   TestAuthService_Register/メールアドレス重複エラー
...
--- PASS: TestAuthService_Register (0.15s)
    --- PASS: TestAuthService_Register/ユーザー登録成功 (0.05s)
    --- PASS: TestAuthService_Register/メールアドレス重複エラー (0.05s)
...
PASS
coverage: 65.2% of statements
ok      github.com/yourusername/sns-app/internal/service    0.753s
```

- `PASS`: テストが成功
- `FAIL`: テストが失敗
- `coverage: XX%`: コードカバレッジ（目標: 60%以上）

### カバレッジレポートの生成

```bash
# テスト環境を起動
make test-setup

# カバレッジレポート付きでテスト実行
docker compose run --rm api_test go test -v -cover -coverprofile=coverage.out ./...

# HTMLレポート生成
docker compose run --rm api_test go tool cover -html=coverage.out -o coverage.html

# テスト環境をクリーンアップ
make test-teardown
```

生成された `backend/coverage.html` をブラウザで開いて確認できます。

### 特定のテストのみ実行

```bash
# テスト環境を起動
make test-setup

# 認証関連のテストのみ実行
docker compose run --rm api_test go test -v ./internal/handler -run TestAuthHandler

# 投稿関連のテストのみ実行
docker compose run --rm api_test go test -v ./internal/handler -run TestPostHandler

# テスト環境をクリーンアップ
make test-teardown
```

---

## フロントエンドE2Eテストの実行

### 基本的な実行方法

```bash
# プロジェクトルートで実行
make test-e2e
```

このコマンドは以下の処理を自動で行います:

1. テスト環境を起動（db_test、api_test）
2. データベースマイグレーションを実行
3. E2Eテストを実行（Playwright）
4. テスト環境をクリーンアップ

### UIモードでテスト実行（推奨）

UIモードを使うと、テストの進行をビジュアルで確認できます。

```bash
# テスト環境を起動
make test-setup

# フロントエンド開発サーバーを**テストモード**で起動（別ターミナル）
cd frontend
npm run dev:test

# UIモードでテスト実行（さらに別ターミナル）
cd frontend
npm run test:e2e:ui

# 完了後、テスト環境をクリーンアップ
make test-teardown
```

**重要:** テストモード（`npm run dev:test`）では、APIのベースURLが `http://localhost:8081/api/v1`（テスト用）に設定されます。通常の `npm run dev` ではテストが失敗します。

### デバッグモードでテスト実行

```bash
# テスト環境を起動
make test-setup

# フロントエンド開発サーバーを**テストモード**で起動（別ターミナル）
cd frontend
npm run dev:test

# デバッグモードでテスト実行（さらに別ターミナル）
cd frontend
npm run test:e2e:debug

# 完了後、テスト環境をクリーンアップ
make test-teardown
```

### 特定のテストファイルのみ実行

```bash
# 認証テストのみ実行
cd frontend
npx playwright test tests/e2e/auth.spec.ts

# 投稿テストのみ実行
npx playwright test tests/e2e/post.spec.ts
```

### テストレポートの表示

```bash
cd frontend
npx playwright show-report
```

---

## 全テストの実行

```bash
# プロジェクトルートで実行
make test
```

このコマンドは以下の順序で実行します:

1. バックエンドテスト（`make test-backend`）
2. フロントエンドE2Eテスト（`make test-e2e`）

**注意:** E2Eテスト実行時は、フロントエンド開発サーバー（`npm run dev:test`）が別ターミナルで**テストモード**起動している必要があります。

---

## トラブルシューティング

### 1. テスト環境が起動しない

**症状:**
```
ERROR: No such service: api_test
```

**解決方法:**
```bash
# docker-compose.ymlを確認
cat docker-compose.yml | grep "api_test"

# Dockerイメージを再ビルド
docker compose --profile test build
```

### 2. データベース接続エラー

**症状:**
```
データベース接続失敗: connection refused
```

**解決方法:**
```bash
# テスト用DBが起動しているか確認
docker compose ps

# テスト用DBが起動していない場合
make test-setup

# ログ確認
docker compose logs db_test
```

### 3. ポート競合エラー

**症状:**
```
Error: Port 5433 is already in use
```

**解決方法:**
```bash
# 使用中のポートを確認
lsof -i :5433

# テスト環境を完全にクリーンアップ
make test-teardown

# 必要に応じて、該当のプロセスを停止してから再度テスト実行
```

### 4. E2Eテストがタイムアウトする（Network Error）

**症状:**
```
Timeout 30000ms exceeded.
```
または画面に「Network Error」が表示される

**原因:**
フロントエンドが間違ったAPIサーバーに接続しようとしている

**解決方法:**
```bash
# 1. フロントエンド開発サーバーを**テストモード**で起動しているか確認
cd frontend
# ❌ 間違い: npm run dev
# ✅ 正しい: npm run dev:test
npm run dev:test

# 2. テスト用APIサーバーが起動しているか確認
docker compose ps | grep api_test

# 3. .env.testファイルが存在するか確認
cat frontend/.env.test
# VITE_API_BASE_URL=http://localhost:8081/api/v1 が設定されているべき

# 4. タイムアウト時間を延長（playwright.config.ts）
# timeout: 30000 → timeout: 60000
```

**重要:** E2Eテストでは必ず `npm run dev:test` を使用してください。通常の `npm run dev` ではポート8080のAPIに接続しようとして失敗します。

### 5. テスト実行後にコンテナが残っている

**症状:**
テスト終了後も `api_test` や `db_test` コンテナが起動したまま。

**解決方法:**
```bash
# テスト環境を完全にクリーンアップ
make test-teardown

# 強制的にすべてのコンテナを停止・削除
docker compose --profile test down -v
```

### 6. Go 1.23のツールチェーンエラー

**症状:**
```
go: download go1.23 for darwin/amd64: toolchain not available
```

**解決方法:**
```bash
# Dockerコンテナ内でのみGoを実行するため、ホストOSのGoは不要
# テストはDockerコンテナ内で実行される
make test-backend
```

### 7. Playwright インストールエラー

**症状:**
```
sh: playwright: command not found
```
または
```
Executable doesn't exist at /path/to/playwright
```

**解決方法:**
```bash
cd frontend
npm install
npx playwright install chromium
# 必要に応じてシステム依存関係もインストール
npx playwright install-deps
```

**補足:** Playwrightブラウザは初回のみインストールが必要です。すでにインストール済みの場合、このエラーは発生しません。

---

## テスト実行のベストプラクティス

### 1. テスト実行前のチェックリスト

- [ ] Dockerが起動している
- [ ] 開発環境のコンテナは停止している（`docker compose ps`で確認）
- [ ] ポート5433と8081が空いている

### 2. テスト失敗時の対応

1. エラーメッセージを詳しく確認
2. ログを確認（`docker compose logs api_test`）
3. テスト環境をクリーンアップして再実行（`make test-teardown` → `make test-backend`）
4. それでも解決しない場合は、このガイドのトラブルシューティングを参照

### 3. 開発時のテスト実行

新しい機能を実装した際は、以下の順序でテストを実行することを推奨します:

```bash
# 1. 実装したレイヤーのユニットテストを実行
make test-setup
docker compose run --rm api_test go test -v ./internal/handler -run TestXXX
make test-teardown

# 2. すべてのバックエンドテストを実行
make test-backend

# 3. E2Eテストを実行（該当する場合）
make test-e2e

# 4. カバレッジを確認
make test-setup
docker compose run --rm api_test go test -cover ./...
make test-teardown
```

---

## 参考リンク

- [Makefile](/Users/sugimoto/Desktop/udemy_pj/app/Makefile) - テストコマンドの定義
- [docker-compose.yml](/Users/sugimoto/Desktop/udemy_pj/app/docker-compose.yml) - テスト環境の構成
- [CLAUDE.md](/Users/sugimoto/Desktop/udemy_pj/app/.claude/CLAUDE.md) - プロジェクト全体のルール
- [Playwright設定](/Users/sugimoto/Desktop/udemy_pj/app/frontend/playwright.config.ts) - E2Eテストの設定

---

## サポート

問題が発生した場合は、以下の情報を添えて報告してください:

1. 実行したコマンド
2. エラーメッセージの全文
3. `docker compose ps` の出力
4. `docker compose logs api_test` の出力（該当する場合）
5. 環境情報（OS、Dockerバージョン）

---

**最後に:** テストは品質保証の要です。機能を実装したら必ずテストを書き、定期的にすべてのテストを実行して、アプリケーションの品質を維持しましょう。
