# プロジェクト概要

## プロジェクト名
**Social Media App** (仮称)

## 概要
テキスト、画像、動画を投稿し、ユーザー同士で交流できるTwitterライクなSNSアプリケーション

## 目的
- ユーザーが自由にコンテンツを投稿・共有できるプラットフォームの提供
- いいね、コメント、フォロー機能によるソーシャルな交流の促進
- レスポンシブデザインによるマルチデバイス対応

## 技術スタック

### フロントエンド
- **言語**: TypeScript
- **フレームワーク**: React
- **UIライブラリ**: Material-UI (MUI)
- **型生成**: openapi-typescript (OpenAPI定義から自動生成)
- **ホスティング**: Firebase Hosting

### バックエンド
- **言語**: Go
- **フレームワーク**: Echo
- **ORM**: GORM
- **OpenAPI生成**: swaggo/echo-swagger
- **ホットリロード**: Air
- **ホスティング**: Render または Google Cloud Run

### データベース
- **RDBMS**: PostgreSQL

### 認証
- **方式**: JWT (JSON Web Token)

### ファイルストレージ
- **サービス**: Firebase Storage

### 開発環境
- **コンテナ**: Docker & Docker Compose
- **バックエンドコンテナ**: Go + Air (ホットリロード)
- **データベースコンテナ**: PostgreSQL

## Phase構成

### Phase 1 (優先度: 高)
基本的なSNS機能の実装
- ユーザー認証（メール+パスワード、JWT）
- ユーザープロフィール（名前、ユーザー名、自己紹介）
- テキスト投稿機能
- 投稿の編集・削除
- いいね機能
- コメント機能（ネスト対応、いいね可能）
- ブックマーク機能
- フォロー/フォロワー機能
- タイムライン（時系列、無限スクロール、20件/ページ）
- 通知機能（1分間隔ポーリング）
- ユーザー検索
- 投稿の通報機能
- 文字数制限: 250文字

### Phase 2 (優先度: 中) - Firebase Storage準備
メディア投稿と追加機能
- Firebase Storage統合
- 画像/動画アップロード機能（最大4枚）
- 動画サイズ制限（環境変数で設定可能、デフォルト20MB）
- ユーザーアイコン・ヘッダー画像
- リツイート（シェア）機能
- ハッシュタグ機能
- 投稿検索機能
- 投稿の公開範囲設定（全体公開、フォロワーのみ、自分のみ）
- ブロック機能
- ユーザー通報機能

### Phase 3 (優先度: 中) - 追加機能
管理・分析機能
- コメント通報機能
- 管理者機能（投稿削除、ユーザー停止等）
- アナリティクス（閲覧数、いいね数推移等）

### Phase 4以降 (優先度: 低)
将来的な拡張
- ソーシャルログイン（Google, Twitter等）
- ダイレクトメッセージ（DM）
- アルゴリズムによるおすすめタイムライン
- サムネイル自動生成（動画・大画像）

## デザイン要件
- **レスポンシブ対応**: PC・スマートフォン両対応
- **テーマ切り替え**: MUIカスタムテーマを複数用意し、UI上で切り替え可能
- **アクセシビリティ**: WCAG 2.1 Level AA準拠を目指す

## セキュリティ要件
- パスワードのハッシュ化（bcrypt推奨）
- JWT トークンの安全な管理
- CORS設定
- XSS、CSRF対策
- SQLインジェクション対策（ORMの適切な使用）
- ファイルアップロードのバリデーション（ファイル形式、サイズ）
- レート制限（API呼び出し制限）

## パフォーマンス要件
- タイムライン読み込み: 2秒以内
- 投稿作成: 3秒以内（メディアなし）
- 画像アップロード: 10秒以内（5MB画像）
- 無限スクロールのスムーズな動作

## データ保持ポリシー
- **論理削除**: ユーザー、投稿、コメント等は論理削除（deleted_atフィールド）
- **ユーザー削除時の対応**: 削除されたユーザーの投稿・コメントは「削除されたユーザー」として表示（データは保持）

## 開発フロー
1. OpenAPI仕様書の作成（バックエンド）
2. swaggoによるOpenAPI定義ファイル自動生成
3. openapi-typescriptによる型定義自動生成（フロントエンド）
4. バックエンド実装（Go + Echo + GORM）
5. フロントエンド実装（React + TypeScript + MUI）
6. Docker Composeでのローカル開発・テスト
7. デプロイ（Render/Cloud Run + Firebase Hosting）

## プロジェクト構成
```
/Users/sugimoto/Desktop/udemy_pj/claudecode/app/
├── frontend/              # フロントエンド
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── tsconfig.json
├── backend/               # バックエンド
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── docs/             # swaggo生成のOpenAPI定義
│   ├── go.mod
│   ├── go.sum
│   └── .air.toml         # Airホットリロード設定
├── docker-compose.yml
├── docs/
│   ├── specs/            # 仕様書
│   └── todo/             # TODOリスト
├── .claude/
│   ├── CLAUDE.md         # プロジェクト全体ルール
│   └── rules/            # カテゴリ別ルール
├── .env.example          # 環境変数サンプル
└── README.md
```

## ドキュメント構成
- `docs/specs/00_OVERVIEW.md`: 本ドキュメント（プロジェクト概要）
- `docs/specs/01_ARCHITECTURE.md`: システムアーキテクチャ
- `docs/specs/02_DATABASE_SCHEMA.md`: データベース設計
- `docs/specs/03_API_SPECIFICATION.md`: API仕様概要
- `docs/specs/04_FRONTEND_SPECIFICATION.md`: フロントエンド仕様
- `docs/specs/05_AUTHENTICATION.md`: 認証フロー
- `docs/specs/06_FILE_UPLOAD.md`: ファイルアップロード仕様

## 次のステップ
各仕様書の詳細を確認し、Phase 1から順次開発を進める
