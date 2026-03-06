# SNS Application - Claude Code プロジェクトルール

## 🚨 最重要原則

### データ保護に関する厳格なルール

#### 絶対に実行してはいけないコマンド
1. **`docker compose down -v`** - ボリュームを削除するため、データベースの全データが消失します
2. **`docker volume rm`** - データベースやファイルが保存されているボリュームを削除します
3. **データベースの DROP/TRUNCATE** - ユーザーの明示的な指示なしに実行禁止

#### データに影響する操作を行う前の必須手順
1. **事前確認**: コマンドがデータに影響する可能性がある場合、必ずユーザーに確認を求める
2. **影響範囲の明示**: 影響範囲を明確に説明する
3. **安全なコマンドの使用**:
   - コンテナ再起動: `docker compose restart [service]`
   - コンテナ停止・起動: `docker compose down` + `docker compose up -d` （-vフラグなし）
   - 再ビルド: `docker compose up -d --build` （-vフラグなし）

---

## 📋 プロジェクト概要

### プロジェクト構造
```
/Users/sugimoto/Desktop/udemy_pj/claudecode/app/
├── frontend/          # React + TypeScript
├── backend/           # Go + Echo
├── docker-compose.yml
├── docs/
│   ├── specs/        # 仕様書
│   └── todo/         # TODOリスト
├── .claude/
│   ├── CLAUDE.md     # 本ファイル
│   └── rules/        # カテゴリ別ルール
└── README.md
```

### 技術スタック
- **フロントエンド**: React, TypeScript, MUI, openapi-typescript
- **バックエンド**: Go, Echo, GORM, swaggo/echo-swagger, Air
- **データベース**: PostgreSQL
- **認証**: JWT
- **ファイルストレージ**: Firebase Storage
- **開発環境**: Docker, Docker Compose
- **デプロイ**: Render/Cloud Run (Backend), Firebase Hosting (Frontend)

---

## 🎯 Phase構成

### Phase 1 (優先度: 高) - 基本SNS機能
- ユーザー認証（メール+パスワード、JWT）
- テキスト投稿機能
- いいね、コメント、ブックマーク
- フォロー/フォロワー
- タイムライン、通知
- ユーザー検索、投稿通報

### Phase 2 (優先度: 中) - メディア機能
- Firebase Storage統合
- 画像/動画アップロード（最大4枚、20MB制限）
- プロフィール画像
- リツイート、ハッシュタグ
- 投稿検索、公開範囲設定
- ブロック機能

### Phase 3 (優先度: 中) - 管理・分析
- 管理者機能（通報管理、ユーザー停止、投稿削除）
- アナリティクス（閲覧数、いいね数推移）

---

## 📖 必須参照ドキュメント

### 仕様書 (`docs/specs/`)
開発を始める前に、必ず該当する仕様書を確認すること：
- `00_OVERVIEW.md` - プロジェクト全体概要
- `01_ARCHITECTURE.md` - システムアーキテクチャ
- `02_DATABASE_SCHEMA.md` - データベース設計
- `03_API_SPECIFICATION.md` - API仕様
- `04_FRONTEND_SPECIFICATION.md` - フロントエンド仕様
- `05_AUTHENTICATION.md` - 認証フロー
- `06_FILE_UPLOAD.md` - ファイルアップロード仕様

### TODOリスト (`docs/todo/`)
**重要: コードを変更するたびに、必ずTODOリストを確認・更新すること**
- `PHASE1_HIGH_PRIORITY.md` - Phase 1 TODO
- `PHASE2_MEDIUM_PRIORITY_PART1.md` - Phase 2 TODO
- `PHASE3_MEDIUM_PRIORITY_PART2.md` - Phase 3 TODO

---

## ✅ TODO管理の厳格なルール

### 必須ルール
1. **コード変更前**: 関連するTODOを確認
2. **コード変更後**: TODOリストを更新
3. **タスク完了時**: 必ずチェックマークを入れる
4. **問題発生時**: TODOに詳細を記録

### TODOの更新方法
```bash
# TODOファイルを開いて、完了したタスクにチェックを入れる
- [x] 完了したタスク
- [ ] 未完了のタスク
```

### TODOリストの活用
- 機能実装前: 該当機能のTODOを確認し、実装内容を把握
- 実装中: 進捗に応じてチェックを入れる
- 実装後: 完了したTODOを確認し、次のタスクへ
- テスト: テスト項目もTODOに含まれているため、必ず実施

詳細は `.claude/rules/05_todo_management.md` を参照

---

## 🔧 開発ルール

### データベース操作
詳細は `.claude/rules/01_database.md` を参照

**主要ルール:**
- 論理削除（soft delete）を使用（`deleted_at` フィールド）
- マイグレーションはGORM AutoMigrateを使用
- すべてのテーブルに `id` (UUID), `created_at`, `updated_at`, `deleted_at` を含める
- 外部キー制約を適切に設定
- インデックスは検索に使用されるカラムに設定

### API開発
詳細は `.claude/rules/02_api_development.md` を参照

**主要ルール:**
- OpenAPI仕様書を先に作成し、swaggoで自動生成
- すべてのエンドポイントにswaggoアノテーションを追加
- JWT認証が必要なエンドポイントには `AuthMiddleware` を適用
- エラーレスポンスは統一形式を使用
- ページネーションは `limit`/`offset` 方式

### フロントエンド開発
詳細は `.claude/rules/03_frontend_development.md` を参照

**主要ルール:**
- openapi-typescriptで型を自動生成し、使用する
- React Queryでデータフェッチとキャッシング
- MUIコンポーネントを使用
- レスポンシブデザイン必須（PC・スマホ対応）
- 無限スクロールは20件ずつ読み込み

### Docker操作
詳細は `.claude/rules/04_docker.md` を参照

**主要ルール:**
- **絶対禁止**: `docker compose down -v`（データが消える）
- 安全な再起動: `docker compose restart`
- ホットリロード: Airを使用（バックエンド）
- ボリュームは必ず永続化設定

### Go開発（Docker必須）
**🚨 重要: ホストOSでGoコマンドを直接実行しない**

**必須ルール:**
- すべてのGoコマンドは `docker compose exec api` 内で実行
- ホストOSに Go をインストールする必要はない
- 開発環境はDockerコンテナで完結させる

**例:**
```bash
# ❌ 禁止: ホストOSで直接実行
go mod init
go get -u github.com/labstack/echo/v4

# ✅ 正しい: コンテナ内で実行
docker compose exec api go mod init
docker compose exec api go get -u github.com/labstack/echo/v4
docker compose exec api swag init -g cmd/server/main.go -o docs
```

### Swagger/OpenAPI ドキュメント生成ルール
**🚨 重要: Swaggerドキュメント生成後は必ずOpenAPI 3.0形式への変換も実行する**

**必須手順:**
1. swagでSwagger 2.0ドキュメントを生成
2. swagger2openapiでOpenAPI 3.0形式に変換

**実行コマンド:**
```bash
# 1. Swagger 2.0ドキュメント生成（コンテナ内）
docker compose exec api swag init -g cmd/server/main.go -o docs

# 2. OpenAPI 3.0形式への変換（ホストOS）
npx swagger2openapi backend/docs/swagger.json -o backend/docs/openapi.yaml -y
```

**ルール:**
- swagドキュメントを生成・更新したら**必ず**openapi.yamlへの変換も実行すること
- openapi.yamlはフロントエンド開発で型生成に使用される
- 両方のファイル（swagger.json, openapi.yaml）をGitにコミットすること

---

## 🔒 セキュリティルール

### 認証・認可
- パスワードはbcryptでハッシュ化（cost: 10）
- JWTトークンは7日間有効
- トークンは `Authorization: Bearer <token>` ヘッダーで送信
- 認証エラー時は401を返却
- 権限エラー時は403を返却

### データ保護
- 論理削除を使用（物理削除は最終手段）
- ユーザー削除時、関連データは「削除されたユーザー」として表示
- 機密情報（パスワード、トークン）はログに出力しない
- `.env` ファイルは `.gitignore` に追加

### ファイルアップロード
- ファイル形式を検証（MIME type）
- ファイルサイズ制限（画像: 5MB, 動画: 20MB）
- Firebase Storageのセキュリティルール設定
- URLの検証（Firebase StorageのURLか確認）

---

## 🧪 テストルール

### テスト環境構成

**テスト環境は開発環境とは完全に独立しています:**
- テスト用PostgreSQL（db_test）: ポート5433
- テスト用APIサーバー（api_test）: ポート8081
- テスト用DBはtmpfs（メモリ）を使用し、開発環境のデータに一切影響を与えない

### テストコマンド（Makefile経由）

**必ずMakefile経由でテストを実行すること:**

```bash
# 開発環境
make dev-up         # 開発環境を起動
make dev-down       # 開発環境を停止
make dev-logs       # 開発環境のログを表示

# テスト実行
make test-backend   # バックエンドの単体テストを実行
make test-e2e       # フロントエンドのE2Eテストを実行
make test           # 全テストを実行（backend + e2e）

# テスト環境の手動操作
make test-setup     # テスト環境を起動
make test-teardown  # テスト環境を停止・削除
```

### バックエンドテスト

**実装ルール:**
1. **新しい機能を実装したら対応するテストも書くこと**
2. テストコマンドは **必ずMakefile** を経由して実行
3. テスト用コンテナは **api_test** を使う（開発用 `api` コンテナではない）
4. 各テストケースの開始時に `testutil.CleanupTestData` でデータをクリーンアップ
5. テストケース名は **日本語** で記述
6. エラーメッセージが **日本語** であることを検証
7. 並列実行数は制限（`-parallel 2`）、`t.Parallel()` は使わない

**テスト対象レイヤー:**
- **ハンドラー層**: リクエストのバリデーション、レスポンスのステータスコードとメッセージ、エラーメッセージが日本語で返ること
- **サービス層**: ビジネスロジック、認証（ユーザー登録、ログイン、JWT生成・検証）、認可（他人のリソースへのアクセス制御）
- **ミドルウェア層**: JWT認証ミドルウェア、認証が必要なエンドポイントにトークンなしでアクセスした場合の挙動

**カバレッジ目標:**
- 60%以上を目標
- `go test -cover` で確認

**テストファイル配置:**
- `backend/internal/handler/*_test.go`
- `backend/internal/service/*_test.go`
- `backend/internal/middleware/*_test.go`
- `backend/internal/testutil/testutil.go`（テストヘルパー）

### フロントエンドE2Eテスト

**実装ルール:**
1. Playwrightを使用
2. テスト先は **テスト用APIサーバー** (`http://localhost:8081`) に向ける
3. テスト項目名は **日本語** で記述
4. セレクタは **data-testid属性** を使用
5. バックエンドのバリデーションを確認して、正常系が通る入力値でテスト
6. テスト並列実行数は1に制限（`workers: 1`）
7. ループ処理は検証に必要な最小回数にとどめる

**テストシナリオ:**
- ユーザー登録 → ログイン
- 投稿 → タイムラインに表示されること
- 投稿の編集・削除
- いいね → いいね解除
- コメント → コメント削除
- コメントのいいね → コメントのいいね解除
- コメントの返信 → 返信として表示されること
- フォロー → フォロー解除 → フォロー一覧に反映
- プロフィール編集
- ログアウト → 認証が必要なページにアクセス → ログインページにリダイレクト

**テストファイル配置:**
- `frontend/tests/e2e/*.spec.ts`
- `frontend/tests/helpers/test-helpers.ts`（テストヘルパー）
- `frontend/playwright.config.ts`（Playwright設定）

### テスト実行時の重要な注意事項

**環境別の動作要件を必ず実装すること:**
- 「本番環境でのみ」「開発環境では無効」などの環境別の動作要件がある場合、必ず実装すること
- テスト環境では `ENV=test` が設定されているため、環境変数で動作を制御

**実装前の確認:**
- 実装前に要件をすべて箇条書きで列挙し、1つも漏らさず実装すること
- テストの並列実行数はマシン負荷を抑えるため制限すること（go test: `-parallel 2`、Playwright: `workers: 1`）
- テストケース内のループは検証に必要な最小回数にとどめること

**エラーメッセージの検証:**
- エラーメッセージはユーザーにとって適切な内容で表示されるようにテストで保証すること
- 満たせていない場合は実装を修正すること

---

## 🔄 CI/CD（継続的インテグレーション）

### GitHub Actions CI

**自動実行されるテスト:**
- **バックエンド単体テスト**: `go test -v -cover ./...`
- **フロントエンドビルド確認**: `npm run build`
- **E2Eテスト**: `npx playwright test`

**トリガー:**
- `main` ブランチへの `push`
- `main` ブランチへの Pull Request

### CI実行環境

**バックエンド単体テスト:**
- PostgreSQL 15（サービスコンテナ）
- Go 1.23
- テスト用DB: `social_media_db_test`
- 環境変数: `ENV=test`, `JWT_SECRET=test-secret-key`

**E2Eテスト:**
- PostgreSQL 15（サービスコンテナ）
- Go 1.23（APIサーバー: ポート8081）
- Node.js 20（Playwright）
- Chromiumブラウザ
- ヘルスチェック: `http://localhost:8081/health`

### CI関連の厳格なルール

**必須ルール:**
1. **mainへのマージ前にCIが全てグリーンであること**
   - すべてのジョブ（backend-test, frontend-build, e2e-test）がPASSしていること
   - CIが失敗している場合は、原因を特定し修正してからマージ

2. **CIで実行されるテストがローカルでも通ることを確認してからpushすること**
   - push前に必ずローカルでテストを実行
   - ローカルテストコマンド:
     ```bash
     make test-backend   # バックエンド単体テスト
     make test-e2e       # E2Eテスト
     make test           # 全テスト実行
     ```

3. **CIエラー時の対応**
   - GitHub Actionsのログを確認
   - テスト失敗時はアーティファクト（Playwrightレポート、スクリーンショット）をダウンロードして原因調査
   - 環境差異がある場合は、CIの環境変数や設定を確認

4. **Pull Request作成時の確認事項**
   - [ ] ローカルで全テストがPASS
   - [ ] CIが全てグリーン
   - [ ] コードレビュー準備完了
   - [ ] TODOリスト更新済み

### CIキャッシュ設定

**高速化のためのキャッシュ:**
- Goモジュール: `~/go/pkg/mod`, `~/.cache/go-build`
- Node.js依存関係: `node_modules`（npm ci使用）
- Playwrightブラウザ: `npx playwright install chromium --with-deps`

### テスト失敗時のデバッグ

**Playwrightテスト失敗時:**
1. GitHub Actionsの「Artifacts」タブからレポートとスクリーンショットをダウンロード
2. `playwright-report/index.html` をブラウザで開く
3. 失敗したテストケースのスクリーンショットとトレースを確認
4. ローカルで再現して修正

**バックエンドテスト失敗時:**
1. GitHub Actionsのログで失敗したテストケースを特定
2. ローカルで該当テストを実行: `make test-backend`
3. 必要に応じて環境変数やDB状態を確認

---

## 📝 コミットルール

### コミットメッセージ形式
```
<type>: <subject>

<body>
```

**Type:**
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント更新
- `style`: コードフォーマット
- `refactor`: リファクタリング
- `test`: テスト追加・修正
- `chore`: その他（ビルド、設定等）

**例:**
```
feat: add user registration endpoint

- POST /api/v1/auth/register を実装
- パスワードのbcryptハッシュ化
- JWT生成機能追加
```

---

## 🚀 デプロイルール

### デプロイ前チェックリスト
- [ ] すべてのテストがパス
- [ ] 環境変数が正しく設定されている
- [ ] データベースマイグレーションが完了
- [ ] Firebase Storageの設定が完了
- [ ] Swagger UIでAPI仕様が確認できる
- [ ] フロントエンドのビルドが成功

### デプロイ手順
1. バックエンドデプロイ（Render/Cloud Run）
2. データベースマイグレーション実行
3. フロントエンドビルド
4. Firebase Hostingへデプロイ
5. 動作確認

---

## 📚 追加リソース

### 公式ドキュメント
- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [React](https://react.dev/)
- [Material-UI](https://mui.com/)
- [Firebase](https://firebase.google.com/docs)

### ツール
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **PostgreSQL**: `localhost:5432`
- **フロントエンド開発サーバー**: `localhost:3000`
- **バックエンドAPI**: `localhost:8080`

---

## ⚠️ トラブルシューティング

### データベース接続エラー
1. Dockerコンテナが起動しているか確認: `docker compose ps`
2. 環境変数が正しいか確認: `.env` ファイル
3. データベースログ確認: `docker compose logs db`

### Airホットリロードが動かない
1. `.air.toml` 設定確認
2. コンテナ再起動: `docker compose restart backend`
3. ログ確認: `docker compose logs backend`

### Firebase Storage アップロードエラー
1. Firebase設定確認: `.env` ファイル
2. セキュリティルール確認
3. サービスアカウントキー確認

---

## 📞 サポート

問題が発生した場合:
1. 該当する `.claude/rules/` のルールファイルを確認
2. `docs/specs/` の仕様書を確認
3. `docs/todo/` のTODOリストを確認
4. エラーログを詳細に記録
5. 必要に応じてユーザーに質問

---

**最後に**: このプロジェクトの成功は、ルールの遵守とTODOリストの適切な管理にかかっています。常にドキュメントを参照し、慎重に開発を進めてください。
