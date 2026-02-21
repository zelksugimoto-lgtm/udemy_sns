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

### バックエンドテスト
- すべてのエンドポイントに対してテストを作成
- 正常系・異常系の両方をテスト
- 認証が必要なエンドポイントは認証ありなしの両方をテスト

### フロントエンドテスト
- 主要なユーザーフローをE2Eテスト
- レスポンシブデザインを確認
- エラーハンドリングを確認

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
