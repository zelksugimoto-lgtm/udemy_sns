# 非機能要件診断レポート

**作成日**: 2026-03-07
**対象**: SNS Application (Phase 1)
**診断者**: Claude Code

---

## 📊 診断サマリー

| カテゴリ | 達成率 | 重要度 |
|---------|-------|--------|
| 性能 | 🔴 20% (1/5) | High |
| 運用・保守性 | 🟡 67% (2/3) | Medium |
| セキュリティ | 🔴 20% (1/5) | High |

**総合評価**: 🔴 **要改善** (36% 達成)

---

## 1. 性能要件

### 1.1 レスポンスタイム

**ステータス**: ❌ **未達成**

**要件**: 一覧取得API：500ms以内

**現状**:
- タイムアウト設定なし
- レスポンスタイムの測定・監視なし
- パフォーマンステストなし

**根拠**:
- `backend/cmd/server/main.go`: タイムアウトミドルウェアの設定なし
- `backend/internal/config/database.go:31-33`: データベース接続タイムアウトの明示的設定なし

**不足点**:
- Echoサーバーのタイムアウト設定（ReadTimeout, WriteTimeout）が未設定
- データベースクエリタイムアウトの設定なし
- レスポンスタイム計測ミドルウェアなし

**改善提案**:
```go
// cmd/server/main.go
e.Server.ReadTimeout = 30 * time.Second
e.Server.WriteTimeout = 30 * time.Second

// タイムアウトミドルウェア追加
e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
    Timeout: 30 * time.Second,
}))

// レスポンスタイム計測ミドルウェア
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        err := next(c)
        duration := time.Since(start)
        log.Printf("Request: %s %s - Duration: %v", c.Request().Method, c.Path(), duration)
        return err
    }
})
```

**優先度**: **High**

---

### 1.2 N+1クエリ対策

**ステータス**: ❌ **未達成（重大な問題あり）**

**要件**: N+1クエリを発生させない

**現状**:
- **深刻なN+1クエリ問題**が複数箇所で発生
- タイムライン取得時に1件の投稿につき4クエリが追加発行される
- 20件取得時: 1 + (20 × 4) = **81クエリ**

**根拠**:
`backend/internal/service/post_service.go:160-168`:
```go
for i, post := range posts {
    likesCount, _ := s.postRepo.CountLikes(post.ID)       // N+1
    commentsCount, _ := s.postRepo.CountComments(post.ID) // N+1
    isLiked, _ := s.likeRepo.Exists(userID, "Post", post.ID) // N+1
    isBookmarked, _ := s.bookmarkRepo.Exists(userID, post.ID) // N+1
    // ...
}
```

同様の問題が以下のメソッドにも存在:
- `GetUserPosts` (post_service.go:199-211)
- `GetUserLikedPosts` (post_service.go:243-255)

**不足点**:
- バッチクエリでのカウント取得が未実装
- Preloadやサブクエリでの一括取得が未実装

**改善提案**:
```go
// リポジトリに一括取得メソッドを追加
func (r *postRepository) CountLikesInBatch(postIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
    type Result struct {
        PostID uuid.UUID
        Count  int64
    }
    var results []Result
    err := r.db.Model(&model.Like{}).
        Select("likeable_id as post_id, COUNT(*) as count").
        Where("likeable_type = ? AND likeable_id IN ?", "Post", postIDs).
        Group("likeable_id").
        Scan(&results).Error

    counts := make(map[uuid.UUID]int64)
    for _, r := range results {
        counts[r.PostID] = r.Count
    }
    return counts, err
}

// サービス層で一括取得
postIDs := make([]uuid.UUID, len(posts))
for i, post := range posts {
    postIDs[i] = post.ID
}

likesCounts, _ := s.postRepo.CountLikesInBatch(postIDs)
commentsCounts, _ := s.postRepo.CountCommentsInBatch(postIDs)
// ...
```

**優先度**: **High（最優先）**

---

### 1.3 アップロード制限

**ステータス**: 🟡 **部分的達成**

**要件**: リクエストボディ上限：画像5MB、動画50MB（ミドルウェアで制限）

**現状**:
- 環境変数で制限値を定義済み（`docker-compose.yml:50-52`）
- **ミドルウェアでの実装なし**（Phase 2の機能のため未実装）

**根拠**:
```yaml
# docker-compose.yml
MAX_VIDEO_SIZE_MB: 20
MAX_IMAGE_SIZE_MB: 5
MAX_UPLOAD_FILES: 4
```

- `backend/cmd/server/main.go`: ボディサイズ制限ミドルウェアなし

**不足点**:
- ミドルウェアでのボディサイズ制限が未実装
- ファイルタイプのバリデーションなし

**改善提案** (Phase 2で実装予定):
```go
// cmd/server/main.go
e.Use(middleware.BodyLimit("50M")) // 最大50MB

// ファイルアップロード専用ミドルウェア
func FileSizeLimit(maxSize int64) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if c.Request().ContentLength > maxSize {
                return c.JSON(http.StatusRequestEntityTooLarge,
                    errors.BadRequest("ファイルサイズが大きすぎます"))
            }
            return next(c)
        }
    }
}
```

**優先度**: **Medium（Phase 2で対応）**

---

### 1.4 ページネーション

**ステータス**: ✅ **達成**

**要件**: 一覧取得で全件取得を禁止。カーソルベースまたはオフセット方式で20件単位

**現状**:
- オフセット方式で実装済み
- デフォルト20件
- 全件取得は不可

**根拠**:
`backend/internal/handler/post_handler.go:230-239`:
```go
limit, _ := strconv.Atoi(c.QueryParam("limit"))
offset, _ := strconv.Atoi(c.QueryParam("offset"))

if limit <= 0 {
    limit = 20  // デフォルト20件
}
if offset < 0 {
    offset = 0
}
```

`backend/internal/repository/post_repository.go:63-67`:
```go
err := query.Preload("User").
    Order("created_at DESC").
    Limit(limit).Offset(offset).  // ページネーション実装済み
    Find(&posts).Error
```

**優先度**: N/A（達成済み）

---

### 1.5 タイムアウト

**ステータス**: ❌ **未達成**

**要件**: APIリクエスト：30秒以内、DB接続：5秒以内

**現状**:
- APIリクエストタイムアウト: 未設定
- DB接続タイムアウト: 未設定（GORMデフォルト）

**根拠**:
- `backend/cmd/server/main.go:51-166`: タイムアウト設定なし
- `backend/internal/config/database.go:15-44`: 接続タイムアウト未設定

**不足点**:
- HTTPサーバーのタイムアウト設定なし
- データベース接続プールのタイムアウト設定なし

**改善提案**:
```go
// cmd/server/main.go - サーバータイムアウト
e.Server.ReadTimeout = 30 * time.Second
e.Server.WriteTimeout = 30 * time.Second
e.Server.IdleTimeout = 120 * time.Second

// config/database.go - DB接続タイムアウト
sqlDB, err := db.DB()
if err != nil {
    return nil, err
}

sqlDB.SetConnMaxLifetime(time.Minute * 3)
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxIdleTime(time.Minute * 5)
```

**優先度**: **High**

---

## 2. 運用・保守性要件

### 2.1 ログ

**ステータス**: 🟡 **部分的達成**

**要件**: アクセスログ・エラーログを構造化ログ(JSON形式)で記録

**現状**:
- アクセスログ: Echoのデフォルトロガー（非構造化）
- エラーログ: 標準出力（非構造化）
- データベースログ: GORMのInfoレベル（非構造化）

**根拠**:
`backend/cmd/server/main.go:54-55`:
```go
e.Use(middleware.Logger())  // デフォルトロガー（非JSON）
e.Use(middleware.Recover())
```

`backend/internal/config/database.go:31-33`:
```go
config := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // 非構造化ログ
}
```

**不足点**:
- JSON形式の構造化ログではない
- ログレベルの設定が柔軟でない
- リクエストID、ユーザーIDなどのコンテキスト情報がログに含まれていない

**改善提案**:
```go
// zapやlogrusを使用した構造化ログ
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
    LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
        logger.Info("request",
            zap.String("method", v.Method),
            zap.String("uri", v.URI),
            zap.Int("status", v.Status),
            zap.Duration("latency", v.Latency),
            zap.String("remote_ip", v.RemoteIP),
        )
        return nil
    },
}))
```

**優先度**: **Medium**

---

### 2.2 自動テスト

**ステータス**: ✅ **達成**

**要件**: ユニットテストおよびE2Eテストで品質担保

**現状**:
- バックエンド単体テスト実装済み
- E2Eテスト実装済み
- Makefileでテスト自動化済み

**根拠**:
**バックエンドテスト:**
- `backend/internal/service/auth_service_test.go`
- `backend/internal/middleware/auth_middleware_test.go`
- `backend/internal/handler/*_test.go` (post, comment, like, follow, auth)

**E2Eテスト:**
- `frontend/tests/e2e/auth.spec.ts`
- `frontend/tests/e2e/post.spec.ts`
- `frontend/tests/e2e/interaction.spec.ts`
- `frontend/tests/e2e/follow.spec.ts`

**Makefile:**
```makefile
test-backend:  # バックエンド単体テスト
test-e2e:      # E2Eテスト
test:          # 全テスト実行
```

**優先度**: N/A（達成済み）

---

### 2.3 エラーハンドリング

**ステータス**: ✅ **達成**

**要件**: ユーザーに適切で安全なエラーメッセージを返却

**現状**:
- 日本語エラーメッセージ実装済み
- 適切なHTTPステータスコード返却
- バリデーションエラーも日本語化

**根拠**:
`backend/pkg/validator/validator.go:30-58`:
```go
func translateError(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return fmt.Sprintf("%sは必須です", fieldName)
    case "email":
        return fmt.Sprintf("%sの形式が正しくありません", fieldName)
    // ...
    }
}
```

`backend/internal/handler/post_handler.go:99-102`:
```go
if err.Error() == "record not found" || err.Error() == "投稿が見つかりません" {
    return c.JSON(http.StatusNotFound, errors.NotFound("投稿が見つかりません"))
}
```

**優先度**: N/A（達成済み）

---

## 3. セキュリティ要件

### 3.1 CORS

**ステータス**: 🟡 **部分的達成**

**要件**: フロントエンドのオリジンのみ許可。credentials: true

**現状**:
- CORSミドルウェア使用
- **デフォルト設定のまま**（全オリジン許可の可能性）
- credentials設定なし

**根拠**:
`backend/cmd/server/main.go:56`:
```go
e.Use(middleware.CORS())  // デフォルト設定
```

**不足点**:
- 許可オリジンの明示的設定なし
- credentials: true の設定なし
- 許可メソッド・ヘッダーの制限なし

**改善提案**:
```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
    AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
    AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
    AllowCredentials: true,
    MaxAge:           86400, // 24時間
}))
```

**優先度**: **High**

---

### 3.2 認証

**ステータス**: ❌ **未達成**

**要件**: HttpOnly Cookie（SameSite=None; Secure）で管理。アクセストークン1時間、リフレッシュ7日

**現状**:
- **Bearer Token認証**（Authorizationヘッダー）
- Cookieは使用していない
- トークン有効期限: 7日（アクセストークンのみ、リフレッシュトークンなし）

**根拠**:
`backend/internal/service/auth_service.go:134-140`:
```go
claims := jwt.MapClaims{
    "user_id": userID.String(),
    "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7日間
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString([]byte(secret))
```

`backend/internal/middleware/auth_middleware.go:19-22`:
```go
authHeader := c.Request().Header.Get("Authorization")  // Headerから取得
if authHeader == "" {
    return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("認証が必要です"))
}
```

**不足点**:
- HttpOnly Cookieでのトークン管理なし
- アクセストークン（1時間）とリフレッシュトークン（7日）の分離なし
- SameSite, Secureフラグの設定なし
- リフレッシュトークンのエンドポイントなし

**改善提案**:
```go
// ハンドラーでCookieに設定
func (h *AuthHandler) Login(c echo.Context) error {
    // ... ログイン処理 ...

    // アクセストークン (1時間)
    accessToken, _ := generateAccessToken(user.ID)
    c.SetCookie(&http.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        HttpOnly: true,
        Secure:   true,              // HTTPS必須
        SameSite: http.SameSiteNoneMode,
        MaxAge:   3600,              // 1時間
        Path:     "/",
    })

    // リフレッシュトークン (7日)
    refreshToken, _ := generateRefreshToken(user.ID)
    c.SetCookie(&http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteNoneMode,
        MaxAge:   604800,            // 7日間
        Path:     "/api/v1/auth/refresh",
    })

    return c.JSON(http.StatusOK, result)
}

// ミドルウェアでCookieから取得
func AuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            cookie, err := c.Cookie("access_token")
            if err != nil {
                return c.JSON(http.StatusUnauthorized,
                    pkgErrors.Unauthorized("認証が必要です"))
            }
            tokenString := cookie.Value
            // ... JWT検証 ...
        }
    }
}
```

**優先度**: **High**

---

### 3.3 入力検証

**ステータス**: ✅ **達成**

**要件**: 全入力に対してバリデーション・サニタイズを実施

**現状**:
- go-playground/validator で実装済み
- 日本語エラーメッセージ
- リクエストDTOにバリデーションタグ設定済み

**根拠**:
`backend/pkg/validator/validator.go:17-28`:
```go
func Validate(s interface{}) error {
    if err := validate.Struct(s); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            for _, e := range validationErrors {
                return fmt.Errorf("%s", translateError(e))
            }
        }
        return err
    }
    return nil
}
```

`backend/internal/dto/request/post.go:4-7`:
```go
type CreatePostRequest struct {
    Content    string `json:"content" validate:"required,min=1,max=250"`
    Visibility string `json:"visibility" validate:"required,oneof=public followers private"`
}
```

**優先度**: N/A（達成済み）

---

### 3.4 レートリミット

**ステータス**: ❌ **未達成**

**要件**: 認証系API：5回/分、一般API：60回/分

**現状**:
- レートリミットミドルウェアなし
- IP制限なし
- ユーザー単位の制限なし

**根拠**:
- `backend/cmd/server/main.go`: レートリミットミドルウェアの設定なし

**不足点**:
- レートリミットの実装が全くない
- DDoS対策なし

**改善提案**:
```go
import "golang.org/x/time/rate"

// IP単位のレートリミット
var limiter = rate.NewLimiter(rate.Limit(60), 120) // 60回/分、バースト120

func RateLimitMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if !limiter.Allow() {
                return c.JSON(http.StatusTooManyRequests,
                    errors.TooManyRequests("リクエストが多すぎます。しばらく待ってから再試行してください"))
            }
            return next(c)
        }
    }
}

// 認証エンドポイント専用（5回/分）
func AuthRateLimitMiddleware() echo.MiddlewareFunc {
    limiter := rate.NewLimiter(rate.Limit(5), 10)
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if !limiter.Allow() {
                return c.JSON(http.StatusTooManyRequests,
                    errors.TooManyRequests("ログイン試行回数が多すぎます"))
            }
            return next(c)
        }
    }
}

// 使用例
auth := api.Group("/auth")
auth.Use(AuthRateLimitMiddleware())
auth.POST("/register", authHandler.Register)
auth.POST("/login", authHandler.Login)
```

**優先度**: **High**

---

### 3.5 脆弱性対策

**ステータス**: ❌ **未達成**

**要件**: security-review による定期的なセキュリティレビュー実施

**現状**:
- セキュリティレビューの実施なし
- 脆弱性スキャンツール未導入
- 依存関係の脆弱性チェックなし

**根拠**:
- プロジェクト全体を調査したが、セキュリティレビュープロセスの記録なし

**不足点**:
- 定期的なセキュリティレビューのプロセスがない
- 依存関係の脆弱性スキャンなし
- OWASP Top 10への対応状況が不明

**改善提案**:
```bash
# 依存関係の脆弱性スキャン
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# GitHub ActionsでCIに組み込み
- name: Run govulncheck
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...

# フロントエンド
npm audit
npm audit fix

# セキュリティレビューチェックリスト作成
docs/SECURITY_REVIEW_CHECKLIST.md
```

**定期レビュー項目**:
- [ ] OWASP Top 10 への対応確認
- [ ] SQLインジェクション対策
- [ ] XSS対策
- [ ] CSRF対策
- [ ] 認証・認可の脆弱性
- [ ] センシティブ情報の漏洩
- [ ] 依存関係の脆弱性

**優先度**: **Medium**

---

## 📋 改善TODO一覧

### 優先度: High（最優先）

1. **N+1クエリ問題の解決**
   - [ ] `CountLikesInBatch`, `CountCommentsInBatch` メソッドをリポジトリに追加
   - [ ] `ExistsInBatch` メソッドを実装（いいね、ブックマーク）
   - [ ] `GetTimeline`, `GetUserPosts`, `GetUserLikedPosts` でバッチクエリを使用
   - [ ] パフォーマンステストで改善を確認

2. **タイムアウト設定**
   - [ ] HTTPサーバーのタイムアウト設定 (ReadTimeout, WriteTimeout)
   - [ ] データベース接続プールの設定 (ConnMaxLifetime, MaxOpenConns)
   - [ ] タイムアウトミドルウェアの追加

3. **CORS設定の強化**
   - [ ] 許可オリジンを明示的に設定
   - [ ] `AllowCredentials: true` を設定
   - [ ] 環境変数でオリジンを管理

4. **Cookie認証への移行**
   - [ ] アクセストークン（1時間）とリフレッシュトークン（7日）を分離
   - [ ] HttpOnly Cookie での管理
   - [ ] SameSite=None, Secure フラグ設定
   - [ ] リフレッシュトークンエンドポイント実装
   - [ ] 認証ミドルウェアをCookie対応に変更

5. **レートリミット実装**
   - [ ] IP単位のレートリミットミドルウェア（60回/分）
   - [ ] 認証エンドポイント専用レートリミット（5回/分）
   - [ ] ユーザー単位のレートリミット（オプション）

### 優先度: Medium

6. **構造化ログの導入**
   - [ ] zap または logrus の導入
   - [ ] JSON形式のログ出力
   - [ ] リクエストID、ユーザーIDなどのコンテキスト情報追加
   - [ ] ログレベルの環境変数管理

7. **アップロード制限（Phase 2）**
   - [ ] ボディサイズ制限ミドルウェア
   - [ ] ファイルタイプバリデーション
   - [ ] ファイルサイズチェック（画像5MB、動画50MB）

8. **セキュリティレビュープロセス**
   - [ ] `govulncheck` の定期実行
   - [ ] `npm audit` の定期実行
   - [ ] セキュリティレビューチェックリスト作成
   - [ ] GitHub ActionsでCIに組み込み

### 優先度: Low

9. **レスポンスタイム計測**
   - [ ] レスポンスタイム計測ミドルウェア
   - [ ] パフォーマンス監視ダッシュボード（将来的に）

---

## 📈 推奨実装順序

### フェーズ1: 性能・セキュリティの重大問題解決（1週間）
1. N+1クエリ問題の解決
2. タイムアウト設定
3. CORS設定の強化
4. レートリミット実装

### フェーズ2: 認証の強化（3日）
5. Cookie認証への移行
6. リフレッシュトークン実装

### フェーズ3: 運用性の向上（2日）
7. 構造化ログの導入
8. セキュリティレビュープロセスの確立

### フェーズ4: Phase 2機能（別途スケジュール）
9. アップロード制限の実装

---

## 🎯 達成目標

### 短期目標（2週間以内）
- N+1クエリ問題を解決し、性能を大幅改善
- タイムアウト設定とレートリミットで安定性・セキュリティを確保
- CORS設定を適切に設定

### 中期目標（1ヶ月以内）
- Cookie認証への完全移行
- 構造化ログで運用性向上
- セキュリティレビュープロセスの確立

### 最終目標
- 全非機能要件を80%以上達成
- 本番環境デプロイの準備完了

---

## 📞 サポート・質問

このレポートについて質問がある場合は、以下を参照してください：
- `.claude/CLAUDE.md`: プロジェクトルール
- `docs/specs/`: 仕様書
- `docs/todo/`: TODOリスト

---

**次のステップ**: このレポートを基に、優先度Highの項目から順に実装を進めてください。
