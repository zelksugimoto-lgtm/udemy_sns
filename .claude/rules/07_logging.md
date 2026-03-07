# ログ設計ルール

## 構造化ログの必須使用

### ログライブラリ
- **zerolog** を使用すること
- 全てのログをJSON形式で出力
- 環境変数 `ENV` に応じてログレベルと出力形式を切り替え
  - `production`: JSON形式、INFO以上
  - `development`: コンソール形式（色付き）、DEBUG以上
  - `test`: JSON形式、WARN以上

### ロガー初期化
```go
import (
    appLogger "github.com/yourusername/sns-app/pkg/logger"
    "github.com/rs/zerolog/log"
)

func main() {
    env := os.Getenv("ENV")
    appLogger.Init(env)

    log.Info().Str("key", "value").Msg("message")
}
```

## ログレベルの使い分け

### DEBUG
- データベースクエリ（開発時のみ）
- 詳細なデバッグ情報
- 関数の入出力値

```go
log.Debug().
    Str("sql", sql).
    Dur("elapsed", elapsed).
    Msg("database_query")
```

### INFO
- サーバー起動/停止
- アクセスログ
- 正常系の重要な処理

```go
log.Info().
    Str("request_id", requestID).
    Str("method", method).
    Str("path", path).
    Int("status", status).
    Dur("response_time", responseTime).
    Msg("access")
```

### WARN
- リトライ可能なエラー
- 想定内の異常系
- リソース使用率の警告

```go
log.Warn().
    Str("request_id", requestID).
    Err(err).
    Msg("retry_scheduled")
```

### ERROR
- エラー発生時（スタックトレース含む）
- リトライ不可能なエラー
- データ不整合

```go
log.Error().
    Str("request_id", requestID).
    Err(err).
    Str("user_id", userID).
    Msg("failed_to_process_request")
```

### FATAL
- 起動時の致命的エラー（データベース接続失敗など）
- プログラムを終了させる必要がある場合のみ

```go
log.Fatal().Err(err).Msg("データベース接続失敗")
```

## アクセスログ

### 必須ミドルウェア
1. **RequestID** - リクエストごとにユニークなIDを生成
2. **AccessLog** - 全APIリクエストのログを出力

### 適用順序
```go
e.Use(middleware.Recover())                  // パニック時のリカバリー
e.Use(appMiddleware.RequestID())             // リクエストID生成（最優先）
e.Use(appMiddleware.AccessLog())             // アクセスログ（リクエストID取得後）
e.Use(middleware.CORS(...))                  // CORS設定
```

### 出力項目
- タイムスタンプ（自動）
- リクエストID
- HTTPメソッド
- パス
- ステータスコード
- レスポンスタイム
- リモートIP
- ユーザーID（認証済みの場合のみ）
- クエリパラメータ（存在する場合のみ）

### 出力例（JSON）
```json
{
  "level": "info",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/posts",
  "status": 201,
  "response_time": 45.2,
  "remote_ip": "192.168.1.1",
  "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "time": "2026-03-07T12:34:56+09:00",
  "message": "access"
}
```

## リクエストID

### 生成ルール
- UUID v4を使用
- リクエストヘッダー `X-Request-ID` が存在する場合はそれを使用
- 存在しない場合は新規生成

### レスポンスヘッダー
- `X-Request-ID` をレスポンスヘッダーに含める
- クライアントがログと紐付けできるようにする

### コンテキストへの保存
```go
// ミドルウェアで保存
c.Set("request_id", requestID)

// ハンドラーで取得
requestID := middleware.GetRequestID(c)
```

## エラーログとユーザーへのレスポンス

### エラー発生時の処理
1. **詳細なエラーログを出力**（スタックトレース含む）
2. **ユーザーには内部情報を見せない**
3. **リクエストIDのみ含める**

### 実装例
```go
// ハンドラー内でエラー発生時
if err != nil {
    requestID := middleware.GetRequestID(c)

    // 詳細なエラーログ
    log.Error().
        Str("request_id", requestID).
        Err(err).
        Str("user_id", userID).
        Str("post_id", postID).
        Msg("failed_to_create_post")

    // ユーザーへのレスポンス（内部情報を含まない）
    return c.JSON(http.StatusInternalServerError,
        errors.InternalErrorWithRequestID(requestID))
}
```

### エラーレスポンス例
```json
{
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "サーバー内部エラーが発生しました。お手数ですが、リクエストID（550e8400-e29b-41d4-a716-446655440000）をお伝えください。",
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

## 禁止事項

### ❌ 標準ライブラリのlogパッケージを使わない
```go
// ❌ 禁止
log.Println("message")
log.Printf("value: %v", value)
log.Fatal("error")

// ✅ 正しい
log.Info().Msg("message")
log.Info().Any("value", value).Msg("log_value")
log.Fatal().Msg("error")
```

### ❌ fmt.Printlnを使わない
```go
// ❌ 禁止
fmt.Println("debug:", value)

// ✅ 正しい
log.Debug().Any("value", value).Msg("debug")
```

### ❌ ユーザーに内部エラー情報を返さない
```go
// ❌ 禁止
return c.JSON(500, map[string]string{
    "error": err.Error(), // 内部情報が漏洩
})

// ✅ 正しい
requestID := middleware.GetRequestID(c)
log.Error().Err(err).Str("request_id", requestID).Msg("internal_error")
return c.JSON(500, errors.InternalErrorWithRequestID(requestID))
```

### ❌ ログに機密情報を含めない
```go
// ❌ 禁止
log.Info().Str("password", password).Msg("user_login")
log.Debug().Str("token", token).Msg("auth_token")

// ✅ 正しい
log.Info().Str("user_id", userID).Msg("user_login")
log.Debug().Str("token_type", "JWT").Msg("auth_token_generated")
```

## ログの検索とトラブルシューティング

### リクエストIDでログを検索
```bash
# JSON形式のログからリクエストIDで検索
cat app.log | jq 'select(.request_id == "550e8400-...")'
```

### エラーログのみ抽出
```bash
# ERRORレベルのログのみ
cat app.log | jq 'select(.level == "error")'
```

### 特定ユーザーのアクセスログ
```bash
# ユーザーIDで検索
cat app.log | jq 'select(.user_id == "a1b2c3d4-...")'
```

### レスポンスタイムが長いリクエスト
```bash
# 1秒以上かかったリクエスト
cat app.log | jq 'select(.response_time > 1000)'
```

## テスト環境でのログ

### テスト時のログレベル
- `ENV=test` の場合、WARN以上のみ出力
- テストの出力が見やすくなる
- 重要なエラーのみ確認できる

### テストでのログ確認
```go
func TestSomething(t *testing.T) {
    // テスト用にロガーを初期化
    appLogger.Init("test")

    // テスト実行
    // ERRORログが出力されたらテスト失敗とみなす
}
```

## ログローテーション（本番環境）

### 推奨設定
- ファイルサイズ: 100MB
- 保持期間: 30日
- 圧縮: gzip
- ローテーション後の削除: 自動

### 実装例（logrotate）
```
/var/log/sns-app/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0644 app app
    sharedscripts
    postrotate
        systemctl reload sns-app
    endscript
}
```

## まとめ

- **必ず構造化ログを使う**（zerolog）
- **リクエストIDを全ログに含める**
- **ユーザーには内部情報を見せない**
- **適切なログレベルを使い分ける**
- **機密情報をログに出力しない**
- **fmt.PrintlnやlogパッケージのPrintlnは使わない**
