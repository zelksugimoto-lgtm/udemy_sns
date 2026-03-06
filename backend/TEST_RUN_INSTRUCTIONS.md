# 認証関連ユニットテストの実行手順

## 作成したテストファイル

1. `/Users/sugimoto/Desktop/udemy_pj/app/backend/internal/service/auth_service_test.go`
2. `/Users/sugimoto/Desktop/udemy_pj/app/backend/internal/handler/auth_handler_test.go`
3. `/Users/sugimoto/Desktop/udemy_pj/app/backend/internal/middleware/auth_middleware_test.go`

## テスト実行コマンド

### 1. Dockerコンテナを起動

```bash
cd /Users/sugimoto/Desktop/udemy_pj/app
docker compose up -d
```

### 2. すべてのテストを実行

```bash
docker compose exec api go test -v ./internal/service -run TestAuthService -count=1
docker compose exec api go test -v ./internal/handler -run TestAuthHandler -count=1
docker compose exec api go test -v ./internal/middleware -run TestAuthMiddleware -count=1
```

### 3. 個別のテスト実行

#### auth_service_test.go

```bash
# Registerテスト
docker compose exec api go test -v ./internal/service -run TestAuthService_Register -count=1

# Loginテスト
docker compose exec api go test -v ./internal/service -run TestAuthService_Login -count=1

# GetMeテスト
docker compose exec api go test -v ./internal/service -run TestAuthService_GetMe -count=1
```

#### auth_handler_test.go

```bash
# Registerテスト
docker compose exec api go test -v ./internal/handler -run TestAuthHandler_Register -count=1

# Loginテスト
docker compose exec api go test -v ./internal/handler -run TestAuthHandler_Login -count=1

# GetMeテスト
docker compose exec api go test -v ./internal/handler -run TestAuthHandler_GetMe -count=1
```

#### auth_middleware_test.go

```bash
# AuthMiddlewareテスト
docker compose exec api go test -v ./internal/middleware -run TestAuthMiddleware -count=1

# OptionalAuthMiddlewareテスト
docker compose exec api go test -v ./internal/middleware -run TestOptionalAuthMiddleware -count=1

# GetUserIDテスト
docker compose exec api go test -v ./internal/middleware -run TestGetUserID -count=1
```

### 4. すべての認証関連テストを一括実行

```bash
docker compose exec api go test -v ./internal/service ./internal/handler ./internal/middleware -count=1
```

## テストの概要

### auth_service_test.go

- **TestAuthService_Register**
  - 正常系: ユーザー登録成功、パスワードハッシュ化確認、JWT生成
  - 異常系: メールアドレス重複
  - 異常系: ユーザー名重複

- **TestAuthService_Login**
  - 正常系: ログイン成功、JWT生成
  - 異常系: ユーザーが存在しない
  - 異常系: パスワードが間違っている

- **TestAuthService_GetMe**
  - 正常系: ユーザー情報取得成功
  - 異常系: ユーザーが存在しない

### auth_handler_test.go

- **TestAuthHandler_Register**
  - 正常系: 有効なデータで登録成功、201、JWTトークン返却
  - 異常系: 必須項目不足（email, username, display_name, password）、400
  - 異常系: メールアドレス重複、409
  - 異常系: ユーザー名重複、409

- **TestAuthHandler_Login**
  - 正常系: 正しいメール・パスワードでログイン成功、200、JWTトークン返却
  - 異常系: 必須項目不足（email, password）、400
  - 異常系: メールアドレス不一致、401
  - 異常系: パスワード不一致、401

- **TestAuthHandler_GetMe**
  - 正常系: 有効なJWTで自分の情報取得、200
  - 異常系: JWTなし、401
  - 異常系: 無効なJWT、401

### auth_middleware_test.go

- **TestAuthMiddleware**
  - 正常系: 有効なJWTで認証成功
  - 異常系: Authorizationヘッダーなし、401
  - 異常系: Bearer形式でない、401
  - 異常系: 無効なトークン、401
  - 異常系: 有効期限切れトークン、401

- **TestOptionalAuthMiddleware**
  - 正常系: 有効なJWTで認証成功
  - 正常系: Authorizationヘッダーなしでも次へ進む
  - 正常系: 無効なトークンでも次へ進む

- **TestGetUserID**
  - 正常系: user_idが設定されている
  - 異常系: user_idが設定されていない

## テストの特徴

1. **データクリーンアップ**: 各テストケースの開始時に `testutil.CleanupTestData` でデータをクリーンアップ
2. **日本語メッセージ検証**: エラーメッセージが日本語であることを確認
3. **並列実行制限**: `-count=1` で並列実行を制御
4. **テストヘルパー活用**: `testutil` パッケージを使用してテストデータ作成とJWT生成を簡素化

## 注意事項

- テスト実行前に必ずDockerコンテナが起動していることを確認してください
- 環境変数 `JWT_SECRET` はミドルウェアテストで設定されます
- データベース接続情報は `.env` ファイルから読み込まれます
