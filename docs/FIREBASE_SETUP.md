# Firebase Storage セットアップガイド

## 概要

このアプリケーションは、画像アップロード機能にFirebase Storageを使用します。
ローカル環境での動作確認には、Firebaseプロジェクトとサービスアカウントキーの設定が必要です。

## 前提条件

- Firebaseアカウント
- Firebaseプロジェクト（まだ作成していない場合）

---

## 1. Firebaseプロジェクトの作成

1. [Firebase Console](https://console.firebase.google.com/) にアクセス
2. 「プロジェクトを追加」をクリック
3. プロジェクト名を入力（例: `sns-app-dev`）
4. Google Analyticsの設定（任意）
5. プロジェクトを作成

---

## 2. Firebase Storageの有効化

1. Firebase Consoleで作成したプロジェクトを開く
2. 左メニューから「Storage」を選択
3. 「始める」をクリック
4. セキュリティルールは一旦デフォルトのまま「次へ」
5. ロケーションを選択（例: `asia-northeast1` - 東京）
6. 「完了」をクリック

---

## 3. サービスアカウントキーの取得

### 3.1 サービスアカウントキーをダウンロード

1. Firebase Consoleで「プロジェクトの設定」（歯車アイコン）を開く
2. 「サービス アカウント」タブを選択
3. 「新しい秘密鍵の生成」をクリック
4. 「キーを生成」をクリック
5. JSONファイル（例: `your-project-id-firebase-adminsdk-xxxxx.json`）がダウンロードされる

### 3.2 ⚠️ セキュリティ上の注意

**絶対にサービスアカウントキーをGitにコミットしないでください！**

- `.gitignore` に `*.json` を追加済み
- サービスアカウントキーはローカルでのみ使用
- 本番環境ではSecret Manager等で管理

---

## 4. 環境変数の設定

### 4.1 JSON文字列への変換

ダウンロードしたJSONファイルの内容を、**1行の文字列**に変換して環境変数にセットします。

#### macOS / Linux の場合

```bash
# JSONファイルのパスを指定
JSON_FILE=path/to/your-project-id-firebase-adminsdk-xxxxx.json

# 1行の文字列に変換して表示（手動でコピー）
cat $JSON_FILE | tr -d '\n'
```

#### 環境変数に直接セット（ローカル開発）

``bash
# .envファイルを作成・編集
cp .env.example .env
vi .env
```

`.env` ファイルに以下を追加：

```bash
# Firebase サービスアカウントキー（JSON文字列）
FIREBASE_SERVICE_ACCOUNT_KEY='{"type":"service_account","project_id":"your-project","private_key_id":"...","private_key":"-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n","client_email":"...","client_id":"...","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url":"..."}'

# Firebaseプロジェクト情報
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_STORAGE_BUCKET=your-bucket.appspot.com
```

**重要**:
- シングルクォート `'...'` で囲む
- 改行を含まない1行の文字列にする
- `private_key` 内の改行は `\n` として含まれる

### 4.2 Docker環境での設定

`docker-compose.yml` で環境変数を読み込むため、`.env` ファイルの設定だけで動作します。

```bash
# Dockerコンテナ再起動
docker compose restart api
```

---

## 5. フロントエンド環境変数の設定

フロントエンドでFirebase SDKを使用するため、以下の環境変数も設定します。

### `frontend/.env` に追加

```bash
# Firebase設定（Webアプリ用）
VITE_FIREBASE_API_KEY=your-api-key
VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-bucket.appspot.com
```

### 取得方法

1. Firebase Console > プロジェクトの設定
2. 「全般」タブ
3. 「マイアプリ」セクションで「ウェブアプリ」を追加（まだの場合）
4. 表示された設定値をコピー

---

## 6. Firebase Storage セキュリティルールの設定

### 推奨ルール（開発環境）

Firebase Console > Storage > ルール タブで以下を設定：

```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    // 投稿画像（誰でも読み取り可、アップロードは認証必要なし - バックエンドで管理）
    match /posts/{userId}/{filename} {
      allow read: if true;
      allow write: if true;  // 開発環境のみ。本番では認証必須
    }

    // プロフィール画像（Phase 2）
    match /avatars/{userId}/{filename} {
      allow read: if true;
      allow write: if true;  // 開発環境のみ
    }

    // ヘッダー画像（Phase 2）
    match /headers/{userId}/{filename} {
      allow read: if true;
      allow write: if true;  // 開発環境のみ
    }
  }
}
```

**本番環境では認証を必須にすること！**

---

## 7. 動作確認

### 7.1 バックエンドの起動確認

```bash
# Dockerログでエラーがないか確認
docker compose logs api

# エラー例: FIREBASE_SERVICE_ACCOUNT_KEY環境変数が設定されていません
# → .envファイルを確認
```

### 7.2 画像アップロード動作確認

1. アプリケーションを起動
   ```bash
   docker compose up -d
   cd frontend && npm run dev
   ```

2. ブラウザで http://localhost:3000 にアクセス

3. ログイン後、投稿作成画面で画像を選択

4. 投稿作成後、画像が表示されることを確認

5. Firebase Console > Storage でファイルがアップロードされていることを確認

---

## 8. トラブルシューティング

### エラー: `FIREBASE_SERVICE_ACCOUNT_KEY環境変数が設定されていません`

**原因**: 環境変数が設定されていない、または形式が不正

**対処法**:
1. `.env` ファイルに `FIREBASE_SERVICE_ACCOUNT_KEY` が設定されているか確認
2. シングルクォートで囲んでいるか確認
3. JSON文字列が1行になっているか確認
4. Docker コンテナを再起動: `docker compose restart api`

### エラー: `invalid firebase storage URL`

**原因**: クライアントから送信されたURLがFirebase Storage URLではない

**対処法**:
1. フロントエンドの `.env` 設定を確認
2. `VITE_FIREBASE_STORAGE_BUCKET` が正しいか確認
3. 画像アップロード時のURLを確認（DevToolsのNetworkタブ）

### エラー: `Firebase初期化エラー`

**原因**: サービスアカウントキーの形式が不正、またはプロジェクトIDが間違っている

**対処法**:
1. JSONファイルの内容をコピーし直す
2. `FIREBASE_PROJECT_ID` が正しいか確認
3. サービスアカウントキーの権限を確認（Firebase Console）

---

## 9. 本番環境へのデプロイ

### 9.1 環境変数の設定（Render / Cloud Run）

**Render の場合:**
1. Render Dashboard > サービス設定 > Environment
2. `FIREBASE_SERVICE_ACCOUNT_KEY` を追加
3. 値: JSON文字列（シングルクォート不要）

**Cloud Run の場合:**
1. Secret Managerにサービスアカウントキーを登録
2. Cloud Runサービスの環境変数でSecretを参照

### 9.2 セキュリティルールの更新

本番環境では必ず認証を必須にする：

```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    match /posts/{userId}/{filename} {
      allow read: if true;
      allow write: if request.auth != null && request.auth.uid == userId;
    }
  }
}
```

---

## まとめ

- ✅ Firebaseプロジェクト作成
- ✅ Firebase Storage有効化
- ✅ サービスアカウントキー取得
- ✅ 環境変数設定（バックエンド・フロントエンド）
- ✅ セキュリティルール設定
- ✅ 動作確認

**重要**: サービスアカウントキーは秘密情報です。Gitにコミットしないよう注意してください。
