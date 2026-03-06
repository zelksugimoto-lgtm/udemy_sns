# E2Eテスト実行ガイド

## 概要

このプロジェクトでは、Playwrightを使用したE2Eテストを実装しています。

## 前提条件

1. **バックエンドAPIサーバーが起動していること**
   - テスト用APIサーバー: `http://localhost:8081`（環境変数で変更可能）
   - フロントエンド開発サーバー: `http://localhost:3000`

2. **依存関係のインストール**
   ```bash
   cd frontend
   npm install
   ```

3. **Playwrightブラウザのインストール**
   ```bash
   npx playwright install
   ```

## テスト実行

### 通常実行
```bash
npm run test:e2e
```

### UIモードで実行（推奨）
```bash
npm run test:e2e:ui
```

UIモードでは、テストの進行状況を視覚的に確認でき、デバッグがしやすくなります。

### デバッグモードで実行
```bash
npm run test:e2e:debug
```

## テスト構成

### テストファイル
- `tests/e2e/auth.spec.ts` - 認証機能のテスト
- `tests/e2e/post.spec.ts` - 投稿機能のテスト
- `tests/e2e/interaction.spec.ts` - いいね・コメント機能のテスト
- `tests/e2e/follow.spec.ts` - フォロー・プロフィール機能のテスト

### ヘルパー関数
- `tests/helpers/test-helpers.ts` - テストで共通的に使用する関数

## テストシナリオ

### 認証テスト (`auth.spec.ts`)
- ユーザー登録 → ログイン
- 登録失敗（メールアドレス重複）
- ログイン失敗（パスワード間違い）
- ログアウト → 認証が必要なページにアクセス → ログインページにリダイレクト

### 投稿テスト (`post.spec.ts`)
- 投稿作成 → タイムラインに表示される
- 投稿編集 → 内容が更新される（未実装のためスキップ）
- 投稿削除 → タイムラインから消える

### インタラクションテスト (`interaction.spec.ts`)
- いいね → いいね解除
- コメント作成 → コメント削除
- コメントのいいね → コメントのいいね解除
- コメントの返信 → 返信として表示される

### フォロー・プロフィールテスト (`follow.spec.ts`)
- フォロー → フォロー解除 → フォロー一覧に反映
- プロフィール編集

## 環境変数

### テスト対象のURLを変更する場合
```bash
export PLAYWRIGHT_BASE_URL=http://localhost:3000
npm run test:e2e
```

## トラブルシューティング

### テストが失敗する場合
1. **APIサーバーが起動しているか確認**
   ```bash
   curl http://localhost:8081/health
   ```

2. **フロントエンド開発サーバーが起動しているか確認**
   ```bash
   curl http://localhost:3000
   ```

3. **テストレポートを確認**
   ```bash
   npx playwright show-report
   ```

### データベースがリセットされない場合
テストはクリーンな状態から開始する必要があります。各テストで新しいユーザーを作成しています。

## 注意事項

- テストは並列実行数を1に制限しています（`workers: 1`）
- 各テストは独立して実行されるため、テスト間でデータが共有されることはありません
- `data-testid`属性を使用してDOM要素を特定しています
- ループ処理は必要最小限にとどめています

## 参考
- [Playwright公式ドキュメント](https://playwright.dev/)
- [Playwrightベストプラクティス](https://playwright.dev/docs/best-practices)
