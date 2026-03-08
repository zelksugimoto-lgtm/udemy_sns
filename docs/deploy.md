# デプロイガイド

本ドキュメントでは、SNSアプリケーションの本番環境へのデプロイ手順を説明します。

## 📋 目次

- [デプロイ構成の概要](#デプロイ構成の概要)
- [前提条件](#前提条件)
- [Neon（PostgreSQL）のセットアップ](#neonpostgresqlのセットアップ)
- [バックエンド（Cloud Run）のデプロイ](#バックエンドcloud-runのデプロイ)
- [フロントエンド（Firebase Hosting）のデプロイ](#フロントエンドfirebase-hostingのデプロイ)
- [環境変数一覧](#環境変数一覧)
- [ブランチ運用](#ブランチ運用)
- [トラブルシューティング](#トラブルシューティング)

---

## デプロイ構成の概要

```
┌─────────────────────┐
│  Firebase Hosting   │  ← フロントエンド（React）
│  (静的ホスティング)  │
└──────────┬──────────┘
           │ HTTPS
           ▼
┌─────────────────────┐
│    Cloud Run        │  ← バックエンド（Go）
│  (コンテナ実行環境)  │
└──────────┬──────────┘
           │ HTTPS
           ▼
┌─────────────────────┐
│       Neon          │  ← データベース（PostgreSQL）
│  (サーバーレスDB)   │
└─────────────────────┘
```

---

## 前提条件

### 必要なアカウント・サービス

1. **Google Cloud Platform（GCP）アカウント**
   - [GCP Console](https://console.cloud.google.com/)にアクセスできること
   - プロジェクトを作成済みであること
   - 請求先アカウントが有効化されていること

2. **Neon アカウント**
   - [Neon Console](https://console.neon.tech/)にアクセスできること
   - プロジェクトを作成済みであること
   - Free Tierで開始可能

3. **Firebase プロジェクト**
   - [Firebase Console](https://console.firebase.google.com/)にアクセスできること
   - プロジェクトを作成済みであること（GCPプロジェクトと連携）

4. **GitHub リポジトリ**
   - ソースコードがGitHubにプッシュされていること
   - Cloud Runから GitHub連携の権限を付与できること

### 必要なツール

```bash
# Firebase CLI（フロントエンドデプロイ用）
npm install -g firebase-tools

# gcloud CLI（バックエンド管理用・オプション）
# https://cloud.google.com/sdk/docs/install からインストール
```

---

## Neon（PostgreSQL）のセットアップ

Neonは、サーバーレスPostgreSQLデータベースサービスです。自動スケーリング、ブランチング機能、無料枠が特徴です。

### 1. Neon プロジェクトの作成

#### 1-1. Neon にサインアップ/ログイン

1. [Neon Console](https://console.neon.tech/) にアクセス
2. GitHubアカウントでサインアップ/ログイン

#### 1-2. プロジェクトを作成

1. 「Create a project」をクリック
2. 以下の設定を入力:
   - **Project name**: `sns-app`（任意）
   - **PostgreSQL version**: `15` または最新
   - **Region**: `AWS / Tokyo (ap-northeast-1)` （推奨：Cloud Runと近い）
3. 「Create project」をクリック

#### 1-3. データベースを作成

プロジェクト作成時に自動的にデフォルトデータベース（`neondb`）が作成されます。
必要に応じて、別のデータベースを作成できます。

1. Neon Consoleで作成したプロジェクトを開く
2. 左メニューから「Databases」を選択
3. 「Create database」をクリック（オプション）
   - **Database name**: `social_media_db`（任意）
   - **Owner**: デフォルトユーザー
4. 「Create」をクリック

#### 1-4. 接続情報を確認

1. Neon Consoleで「Dashboard」または「Connection Details」を開く
2. 以下の情報をメモ:
   - **Host**: `ep-xxxxx-xxxxx.ap-northeast-1.aws.neon.tech`
   - **Database**: `social_media_db` または `neondb`
   - **User**: デフォルトは `neondb_owner` など
   - **Password**: 自動生成されたパスワード（表示されるのでコピー）
   - **Port**: `5432`

または、「Connection string」から直接コピー:
```
postgresql://user:password@ep-xxxxx-xxxxx.ap-northeast-1.aws.neon.tech/dbname?sslmode=require
```

**重要**:
- Neon接続には **SSL/TLS** が必須です（`sslmode=require`）
- 接続文字列は安全に保管してください

---

## バックエンド（Cloud Run）のデプロイ

### 1. Cloud Runの継続デプロイ設定

#### 1-1. Cloud Runサービスを作成

1. [Cloud Run Console](https://console.cloud.google.com/run) にアクセス
2. 「サービスを作成」をクリック
3. 「ソースリポジトリから新しいリビジョンを継続的にデプロイする」を選択
4. 「Cloud Buildでセットアップ」をクリック

#### 1-2. リポジトリ接続

1. 「リポジトリプロバイダ」で「GitHub」を選択
2. 「リポジトリを認証」をクリックして、GitHubアカウントと連携
3. リポジトリを選択（例: `your-username/sns-app`）
4. 「次へ」をクリック

#### 1-3. ビルド設定

1. **ブランチ**: `^main$`（mainブランチへのpush時に自動デプロイ）
2. **ビルドタイプ**: `Dockerfile`
3. **ソース**: `/backend/Dockerfile.prod`
4. 「保存」をクリック

#### 1-4. サービス設定

1. **サービス名**: `sns-api`（任意）
2. **リージョン**: `asia-northeast1`（東京）
3. **CPU割り当てと料金**:
   - 「リクエストの処理中にのみCPUを割り当てる」を選択
4. **自動スケーリング**:
   - **最小インスタンス数**: 0（コスト削減）または 1（コールドスタート回避）
   - **最大インスタンス数**: 10〜100
5. **認証**:
   - 「未認証の呼び出しを許可」を選択（フロントエンドからアクセスするため）
6. **コンテナ、接続、セキュリティ**:
   - 「コンテナ」タブ:
     - **コンテナポート**: `8080`（環境変数PORTで自動設定）
     - **メモリ**: `512 MiB`〜
     - **CPU**: `1`
     - **リクエストタイムアウト**: `300`秒
     - **最大同時リクエスト数**: `80`
   - 「変数とシークレット」タブ:
     - 環境変数を設定（次のセクション参照）

#### 1-5. 環境変数の設定

「変数とシークレット」タブで以下を設定:

| 変数名 | 値 | 説明 |
|--------|------|------|
| `ENV` | `production` | 実行環境 |
| `DB_HOST` | `ep-xxxxx-xxxxx.ap-northeast-1.aws.neon.tech` | Neonのホスト名 |
| `DB_PORT` | `5432` | PostgreSQLポート |
| `DB_USER` | `neondb_owner` | Neonのデータベースユーザー名 |
| `DB_PASSWORD` | `YOUR_NEON_PASSWORD` | Neonのデータベースパスワード |
| `DB_NAME` | `social_media_db` または `neondb` | データベース名 |
| `DB_SSLMODE` | `require` | SSL接続を強制（Neon必須） |
| `JWT_SECRET` | `強力なランダム文字列` | JWT署名用シークレットキー（32文字以上推奨） |
| `JWT_EXPIRATION` | `168h` | JWTの有効期限（168h = 7日間） |
| `FRONTEND_URL` | `https://your-app.web.app` | フロントエンドのURL（CORS設定用） |
| `ADMIN_ROOT_USERNAME` | `admin` | 管理者ユーザー名 |
| `ADMIN_ROOT_PASSWORD` | `強力なパスワード` | 管理者パスワード |
| `ADMIN_SESSION_SECRET` | `強力なランダム文字列` | セッション署名用シークレットキー（32文字以上推奨） |
| `MAX_VIDEO_SIZE_MB` | `20` | 動画アップロード最大サイズ（MB） |
| `MAX_IMAGE_SIZE_MB` | `5` | 画像アップロード最大サイズ（MB） |
| `MAX_UPLOAD_FILES` | `4` | 一度にアップロード可能なファイル数 |
| `TZ` | `Asia/Tokyo` | タイムゾーン |

**Neon接続情報の確認方法**:
1. [Neon Console](https://console.neon.tech/) を開く
2. プロジェクトを選択
3. 「Dashboard」または「Connection Details」を開く
4. 以下をコピー:
   - **Host**: `DB_HOST`
   - **Database**: `DB_NAME`
   - **User**: `DB_USER`
   - **Password**: `DB_PASSWORD`（「Show password」をクリック）

**重要**: Neonへの接続には **SSL/TLS** が必須です。`DB_SSLMODE=require` を必ず設定してください。

#### 1-6. 作成

「作成」をクリックしてCloud Runサービスを作成します。

### 2. 初回デプロイ確認

1. Cloud Buildが自動的に開始され、Dockerイメージをビルド
2. ビルド完了後、Cloud Runにデプロイ
3. デプロイ完了後、サービスURLが表示される（例: `https://sns-api-xxxxx-an.a.run.app`）
4. 疎通確認:
   ```bash
   curl https://sns-api-xxxxx-an.a.run.app/health
   ```
   レスポンス: `{"status":"ok"}`

### 3. データベースマイグレーション

初回デプロイ後、Cloud Runコンソールから「ログ」タブを開き、マイグレーションログを確認:

```bash
# Cloud Runのログを確認（マイグレーションログが出力されているはず）
gcloud run services logs read sns-api --region asia-northeast1 --limit 50
```

アプリケーションは起動時に自動的にGORM AutoMigrateを実行します。
手動でマイグレーションを実行する必要はありません。

### 4. 継続デプロイの動作確認

1. `main`ブランチにコードをpush
2. Cloud Buildが自動的にトリガーされる
3. [Cloud Build コンソール](https://console.cloud.google.com/cloud-build) でビルド状況を確認
4. ビルド完了後、Cloud Runに自動デプロイ

---

## フロントエンド（Firebase Hosting）のデプロイ

### 1. Firebase プロジェクトのセットアップ

#### 1-1. Firebase プロジェクトを作成

1. [Firebase Console](https://console.firebase.google.com/) にアクセス
2. 「プロジェクトを追加」をクリック
3. 既存のGCPプロジェクトを選択、または新規作成
4. Google Analyticsは任意で有効化

#### 1-2. Firebase CLI にログイン

```bash
firebase login
```

ブラウザが開き、Googleアカウントでログインします。

#### 1-3. Firebaseプロジェクトを初期化

```bash
cd frontend
firebase init hosting
```

以下の質問に回答:
- **プロジェクトを選択**: 作成したFirebaseプロジェクトを選択
- **公開ディレクトリ**: `dist`
- **シングルページアプリとして設定**: `Yes`
- **GitHub Actionsのセットアップ**: `No`（手動デプロイ）
- **既存のindex.htmlを上書き**: `No`

設定完了後、`firebase.json` と `.firebaserc` が作成されます。

#### 1-4. firebase.json の確認

`frontend/firebase.json` の内容を確認:

```json
{
  "hosting": {
    "public": "dist",
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ],
    "rewrites": [
      {
        "source": "**",
        "destination": "/index.html"
      }
    ]
  }
}
```

### 2. 環境変数の設定

#### 2-1. .env.production を編集

`frontend/.env.production` を開き、Cloud RunのURLに変更:

```bash
VITE_API_BASE_URL=https://sns-api-xxxxx-an.a.run.app/api/v1
```

**重要**: Cloud RunのサービスURLは以下で確認できます:
1. [Cloud Run Console](https://console.cloud.google.com/run) を開く
2. `sns-api` サービスをクリック
3. 「URL」をコピー
4. `/api/v1` を末尾に追加

### 3. デプロイ

#### 3-1. Makefileコマンドでデプロイ

```bash
cd /Users/sugimoto/Desktop/udemy_pj/app
make deploy-frontend
```

#### 3-2. 手動デプロイ（スクリプト直接実行）

```bash
cd frontend
./scripts/deploy.sh
```

#### 3-3. デプロイ完了

デプロイが成功すると、Firebase HostingのURLが表示されます:
```
Hosting URL: https://your-app.web.app
```

### 4. カスタムドメインの設定（オプション）

1. [Firebase Console](https://console.firebase.google.com/) を開く
2. 「Hosting」→「カスタムドメインを追加」
3. 所有するドメインを入力（例: `example.com`）
4. DNSレコードを設定（Firebase側で表示される指示に従う）
5. SSL証明書が自動的にプロビジョニングされる

---

## 環境変数一覧

### バックエンド（Cloud Run）

| 変数名 | 必須 | デフォルト | 説明 | 例 |
|--------|------|-----------|------|-----|
| `ENV` | Yes | - | 実行環境（`production`, `development`, `test`） | `production` |
| `PORT` | No | `8080` | Cloud Runが自動設定 | `8080` |
| `DB_HOST` | Yes | - | Neonのホスト名 | `ep-xxxxx-xxxxx.ap-northeast-1.aws.neon.tech` |
| `DB_PORT` | Yes | - | PostgreSQLポート | `5432` |
| `DB_USER` | Yes | - | Neonのユーザー名 | `neondb_owner` |
| `DB_PASSWORD` | Yes | - | Neonのパスワード | `xxxxx` |
| `DB_NAME` | Yes | - | データベース名 | `social_media_db` |
| `DB_SSLMODE` | Yes | - | SSL接続モード（Neon必須） | `require` |
| `JWT_SECRET` | Yes | - | JWT署名用シークレットキー（32文字以上推奨） | `your-super-secret-key-12345678` |
| `JWT_EXPIRATION` | No | `168h` | JWTの有効期限（時間単位） | `168h` |
| `FRONTEND_URL` | Yes | - | フロントエンドのURL（CORS設定用） | `https://your-app.web.app` |
| `ADMIN_ROOT_USERNAME` | Yes | - | 管理者ユーザー名 | `admin` |
| `ADMIN_ROOT_PASSWORD` | Yes | - | 管理者パスワード | `strong-password` |
| `ADMIN_SESSION_SECRET` | Yes | - | セッション署名用シークレットキー（32文字以上推奨） | `session-secret-key-12345678` |
| `MAX_VIDEO_SIZE_MB` | No | `20` | 動画アップロード最大サイズ（MB） | `20` |
| `MAX_IMAGE_SIZE_MB` | No | `5` | 画像アップロード最大サイズ（MB） | `5` |
| `MAX_UPLOAD_FILES` | No | `4` | 一度にアップロード可能なファイル数 | `4` |
| `TZ` | No | `UTC` | タイムゾーン | `Asia/Tokyo` |

### フロントエンド（Firebase Hosting）

| 変数名 | 必須 | 説明 | 例 |
|--------|------|------|-----|
| `VITE_API_BASE_URL` | Yes | バックエンドAPIのベースURL | `https://sns-api-xxxxx-an.a.run.app/api/v1` |

**設定方法**: `frontend/.env.production` ファイルに記述

---

## ブランチ運用

### 推奨ブランチ戦略

```
main (本番環境) ← Cloud Runが継続デプロイ
  ↑
  └─ develop (開発環境) ← PRをマージ
       ↑
       └─ feature/xxx (機能開発ブランチ)
```

### デプロイフロー

1. **機能開発**:
   ```bash
   git checkout -b feature/new-feature
   # 開発作業
   git add .
   git commit -m "feat: add new feature"
   git push origin feature/new-feature
   ```

2. **Pull Request作成**:
   - `feature/new-feature` → `develop` にPR作成
   - レビュー・テスト
   - マージ

3. **本番デプロイ準備**:
   - `develop` → `main` にPR作成
   - 最終確認・テスト
   - マージ

4. **自動デプロイ**:
   - `main`へのマージをトリガーに、Cloud Buildが自動実行
   - Cloud Runに自動デプロイ

### リリース管理

- **タグ付け**: 本番リリース時にタグを付けることを推奨
  ```bash
  git checkout main
  git pull origin main
  git tag -a v1.0.0 -m "Release v1.0.0"
  git push origin v1.0.0
  ```

- **ロールバック**: 問題が発生した場合
  1. Cloud Run コンソールを開く
  2. 「リビジョン」タブを選択
  3. 前のリビジョンを選択し、「すべてのトラフィックを送信」

---

## トラブルシューティング

### バックエンド（Cloud Run）

#### ビルドが失敗する

**原因**: Dockerfile.prod の設定ミス、依存関係のエラー

**対処法**:
1. [Cloud Build Console](https://console.cloud.google.com/cloud-build) でビルドログを確認
2. エラーメッセージを確認
3. ローカルで Dockerfile.prod をビルドしてテスト:
   ```bash
   cd backend
   docker build -f Dockerfile.prod -t sns-api:test .
   docker run -p 8080:8080 --env-file ../.env sns-api:test
   ```

#### データベース接続エラー

**原因**: Neon接続情報の設定ミス、SSL設定漏れ

**対処法**:
1. Cloud Runの環境変数を確認:
   - `DB_HOST`: Neonのホスト名（例: `ep-xxxxx.ap-northeast-1.aws.neon.tech`）
   - `DB_USER`: Neonのユーザー名
   - `DB_PASSWORD`: 正しいパスワード
   - `DB_NAME`: データベース名
   - `DB_SSLMODE`: `require` が設定されているか **必須**
2. Neon Consoleで接続情報を再確認:
   - [Neon Console](https://console.neon.tech/) → Dashboard → Connection Details
3. Neonプロジェクトがアクティブであることを確認
4. Cloud Runのログでエラー詳細を確認:
   ```bash
   gcloud run services logs read sns-api --region asia-northeast1 --limit 50
   ```

#### CORS エラー

**原因**: `FRONTEND_URL` が正しく設定されていない

**対処法**:
1. Cloud Runの環境変数 `FRONTEND_URL` を確認
2. フロントエンドのURLと完全一致しているか確認（末尾の `/` に注意）
3. 例: `https://your-app.web.app` （末尾スラッシュなし）

#### 認証エラー

**原因**: JWT_SECRETが設定されていない、Cookieが正しく送信されていない

**対処法**:
1. Cloud Runの環境変数 `JWT_SECRET` が設定されているか確認
2. ブラウザの開発者ツールで、リクエストヘッダーにCookieが含まれているか確認
3. CORS設定で `credentials: true` が有効か確認

### フロントエンド（Firebase Hosting）

#### デプロイが失敗する

**原因**: Firebase CLIが未インストール、ログインしていない

**対処法**:
1. Firebase CLIをインストール:
   ```bash
   npm install -g firebase-tools
   ```
2. ログイン:
   ```bash
   firebase login
   ```
3. プロジェクトを確認:
   ```bash
   firebase projects:list
   ```

#### ビルドエラー

**原因**: 依存関係のエラー、型定義の不一致

**対処法**:
1. 依存関係を再インストール:
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   ```
2. 型定義を再生成:
   ```bash
   npm run generate:types
   ```
3. ビルドを再実行:
   ```bash
   npm run build
   ```

#### APIリクエストが失敗する

**原因**: `.env.production` のAPI URLが間違っている

**対処法**:
1. `frontend/.env.production` を確認:
   ```bash
   VITE_API_BASE_URL=https://sns-api-xxxxx-an.a.run.app/api/v1
   ```
2. Cloud RunのURLを確認（Cloud Run Console）
3. 末尾の `/api/v1` を確認
4. ビルドとデプロイを再実行:
   ```bash
   make deploy-frontend
   ```

#### ルーティングエラー（404 on refresh）

**原因**: `firebase.json` の `rewrites` 設定が不足

**対処法**:
`frontend/firebase.json` に以下を追加:
```json
{
  "hosting": {
    "public": "dist",
    "rewrites": [
      {
        "source": "**",
        "destination": "/index.html"
      }
    ]
  }
}
```

---

## セキュリティ設定の推奨事項

### 1. Cloud Runのセキュリティ

- **サービスアカウント**: 最小権限の原則に従い、必要最小限の権限のみ付与
- **シークレット管理**: Secret Managerを使用して、環境変数（特にDB_PASSWORDやJWT_SECRET）をより安全に管理

### 2. Neonのセキュリティ

- **SSL/TLS接続**: 必ず `DB_SSLMODE=require` を設定（Neon必須）
- **IP制限**: Neon Consoleで許可するIPアドレスを制限（オプション、Pro以上）
- **パスワード管理**: Neonのパスワードは安全に保管し、定期的にローテーション

### 3. Firebase Hostingのセキュリティ

- **HTTPSのみ**: 自動的にHTTPSが有効（Firebase Hostingのデフォルト）
- **セキュリティルール**: Firebase Storageを使用する場合、適切なセキュリティルールを設定

### 4. CORS設定

- `FRONTEND_URL` を正確に設定し、特定のドメインからのみアクセスを許可

---

## コスト最適化のヒント

### Cloud Run

- **最小インスタンス数**: 0に設定してコスト削減（コールドスタート許容）
- **CPU割り当て**: リクエスト処理中のみCPUを割り当てる設定
- **リソース調整**: メモリとCPUを適切に設定（過剰割り当てを避ける）

### Neon

- **Free Tier**:
  - 1プロジェクト、0.5GB RAM、3GBストレージまで無料
  - 月間100時間のコンピュート時間
- **Auto-suspend**: アイドル状態で自動的にコンピュートを一時停止（Free Tierでは5分）
- **ストレージ最適化**: 不要なブランチやデータを定期的に削除
- **スケールトゥゼロ**: 未使用時は自動的にスケールダウンされるため、コスト効率が高い

### Firebase Hosting

- 無料枠: 月10GBまで無料（個人開発には十分）

---

## まとめ

本ドキュメントに従ってデプロイを実施することで、以下が実現されます:

✅ **バックエンド**: Cloud RunでDockerコンテナを実行、GitHub連携で自動デプロイ
✅ **データベース**: Neonでサーバーレス PostgreSQLを運用（自動スケーリング、高可用性）
✅ **フロントエンド**: Firebase Hostingで静的ファイルをホスティング
✅ **継続的デプロイ**: mainブランチへのpushで自動的にデプロイ
✅ **セキュアな運用**: SSL/TLS通信、環境変数による設定管理、HTTPS通信

何か問題が発生した場合は、トラブルシューティングセクションを参照してください。
