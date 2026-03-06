# Phase 1 セキュリティ監査レポート

**作成日**: 2026-03-04
**対象**: SNS Application - Phase 1 (基本SNS機能)
**監査範囲**: バックエンド(Go), フロントエンド(React/TypeScript), データベース(PostgreSQL)

---

## 📋 目次

1. [エグゼクティブサマリー](#エグゼクティブサマリー)
2. [重大度レベルの定義](#重大度レベルの定義)
3. [発見された脆弱性](#発見された脆弱性)
4. [推奨される対策](#推奨される対策)
5. [セキュアな実装の確認事項](#セキュアな実装の確認事項)

---

## エグゼクティブサマリー

### 監査結果概要

| 重大度 | 件数 | 状態 |
|--------|------|------|
| 🔴 Critical (緊急) | 2 | 要対処 |
| 🟠 High (高) | 4 | 要対処 |
| 🟡 Medium (中) | 5 | 対処推奨 |
| 🔵 Low (低) | 3 | 改善推奨 |
| **合計** | **14** | - |

### 主要な発見事項

- **Critical**: JWT秘密鍵がハードコードされており、本番環境で危険
- **Critical**: CORS設定が全オリジン許可（`middleware.CORS()`）で制限なし
- **High**: レート制限が未実装（ブルートフォース攻撃に脆弱）
- **High**: SQLインジェクション対策は十分だが、検証バリデーションに課題あり
- **Medium**: XSS対策は基本的に問題ないが、入力サニタイズが不十分

---

## 重大度レベルの定義

| レベル | 説明 | 対処期限 |
|--------|------|----------|
| 🔴 Critical | システム全体の侵害、データ漏洩の危険性が高い | 即時 |
| 🟠 High | 重要なセキュリティホール、攻撃者に悪用される可能性が高い | 1週間以内 |
| 🟡 Medium | セキュリティリスクがあるが、直接的な侵害には至らない | 1ヶ月以内 |
| 🔵 Low | 軽微なリスク、セキュリティのベストプラクティスに反する | 適宜対応 |

---

## 発見された脆弱性

### 🔴 Critical (緊急対応必須)

#### 1. JWT秘密鍵のハードコード

**場所**:
- `backend/internal/middleware/auth_middleware.go:34-36`
- `backend/internal/middleware/auth_middleware.go:98-100`
- `backend/internal/service/auth_service.go:129-132`

**問題**:
```go
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    secret = "your-secret-key" // 開発環境用のデフォルト値
}
```

**リスク**:
- デフォルト秘密鍵が推測可能で、JWTトークンの偽造が容易
- 環境変数が設定されていない場合、本番環境でも脆弱なデフォルト値を使用
- セッションハイジャック、なりすましの危険性

**影響範囲**: 全認証システム

**推奨対策**:
1. デフォルト値を削除し、環境変数が必須であることを強制
2. 起動時に環境変数チェックを実施し、未設定の場合はサーバーを起動させない
3. 本番環境では最低32バイト以上のランダムな秘密鍵を使用

```go
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
if len(secret) < 32 {
    log.Fatal("JWT_SECRET must be at least 32 characters")
}
```

---

#### 2. CORS設定が全オリジン許可

**場所**: `backend/cmd/server/main.go:56`

**問題**:
```go
e.Use(middleware.CORS())
```

**リスク**:
- すべてのオリジンからのリクエストを許可（`Access-Control-Allow-Origin: *`）
- CSRF攻撃、データ盗聴、セッション盗難の危険性
- クロスサイトリクエストフォージェリ（CSRF）攻撃に脆弱

**影響範囲**: 全APIエンドポイント

**推奨対策**:
明示的にオリジンを制限し、認証情報の送信を許可する設定に変更

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{
        "http://localhost:3000",    // 開発環境
        "https://yourdomain.com",   // 本番環境
    },
    AllowMethods: []string{
        http.MethodGet,
        http.MethodPost,
        http.MethodPut,
        http.MethodPatch,
        http.MethodDelete,
    },
    AllowHeaders: []string{
        "Origin",
        "Content-Type",
        "Authorization",
    },
    AllowCredentials: true,
    MaxAge: 86400,
}))
```

---

### 🟠 High (高優先度)

#### 3. レート制限の未実装

**場所**: 全APIエンドポイント

**問題**:
- ログイン、登録、パスワードリセットなどの認証エンドポイントにレート制限がない
- 投稿、コメント、いいね等のリソース作成にレート制限がない

**リスク**:
- ブルートフォース攻撃（パスワード総当たり）
- サービス拒否攻撃（DoS）
- スパム投稿の大量作成
- アカウントの列挙攻撃

**影響範囲**:
- `/api/v1/auth/login`
- `/api/v1/auth/register`
- `/api/v1/posts`
- `/api/v1/posts/{id}/comments`
- `/api/v1/posts/{id}/like`

**推奨対策**:

1. **認証エンドポイント**: IP単位で厳格なレート制限
```go
import "github.com/ulule/limiter/v3"
import "github.com/ulule/limiter/v3/drivers/store/memory"

// 認証用: 5分間に5回まで
authLimiter := middleware.RateLimiter(store, limiter.Rate{
    Period: 5 * time.Minute,
    Limit:  5,
})

auth.POST("/login", authHandler.Login, authLimiter)
auth.POST("/register", authHandler.Register, authLimiter)
```

2. **リソース作成エンドポイント**: ユーザーIDベースでレート制限
```go
// 投稿作成: 1時間に30回まで
postLimiter := middleware.RateLimiter(store, limiter.Rate{
    Period: 1 * time.Hour,
    Limit:  30,
})

posts.POST("", postHandler.CreatePost, appMiddleware.AuthMiddleware(), postLimiter)
```

---

#### 4. バリデーションの不足

**場所**:
- `backend/internal/handler/auth_handler.go:42-44`
- `backend/internal/handler/post_handler.go:52-60`

**問題**:
- メールアドレスの形式検証がない
- パスワードの強度検証がない
- ユーザー名の形式検証が緩い（特殊文字、長さ）
- 投稿内容のサニタイズ処理がない

**リスク**:
- 不正なデータの保存
- XSS攻撃のリスク増加
- データベース汚染
- UXの低下（無効なメールでの登録）

**影響範囲**:
- ユーザー登録
- 投稿作成
- コメント作成

**推奨対策**:

1. **バリデーションライブラリの導入**
```bash
go get github.com/go-playground/validator/v10
```

2. **DTOにバリデーションタグを追加**
```go
// request/auth.go
type RegisterRequest struct {
    Email       string `json:"email" validate:"required,email,max=255"`
    Username    string `json:"username" validate:"required,alphanum,min=3,max=30"`
    DisplayName string `json:"display_name" validate:"required,min=1,max=50"`
    Password    string `json:"password" validate:"required,min=8,max=128,password_strength"`
}
```

3. **カスタムバリデータの実装**
```go
// パスワード強度チェック
func passwordStrength(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    // 最低8文字、大文字・小文字・数字を含む
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

    return hasUpper && hasLower && hasNumber && len(password) >= 8
}
```

---

#### 5. ユーザー列挙攻撃への脆弱性

**場所**: `backend/internal/handler/auth_handler.go:48-50`

**問題**:
```go
if err.Error() == "このメールアドレスは既に使用されています" ||
   err.Error() == "このユーザー名は既に使用されています" {
    return c.JSON(http.StatusConflict, errors.Conflict(err.Error()))
}
```

**リスク**:
- 登録時に既存のメールアドレス・ユーザー名を確認できる
- 攻撃者が有効なアカウントを列挙可能
- プライバシー侵害、標的型攻撃の準備に悪用される

**影響範囲**: `/api/v1/auth/register`

**推奨対策**:

1. **エラーメッセージを一般化**
```go
if err != nil {
    // 既存ユーザーか否かを明示しない
    return c.JSON(http.StatusBadRequest, errors.BadRequest("登録に失敗しました"))
}
```

2. **または、メール確認を必須にする** (Phase 2以降で検討)

---

#### 6. SQLインジェクション対策の確認不足

**場所**: `backend/internal/repository/user_repository.go:86-112`

**問題**:
```go
searchPattern := "%" + query + "%"
searchQuery = searchQuery.Where(
    "username ILIKE ? OR display_name ILIKE ?",
    searchPattern, searchPattern,
)
```

**現状**:
- GORMのプレースホルダ（`?`）を使用しているため、基本的にSQLインジェクションは防げている
- しかし、ワイルドカード文字（`%`, `_`）のエスケープ処理がない

**リスク**:
- 検索結果の意図しない動作
- パフォーマンス問題（`%%`の検索など）

**影響範囲**: `/api/v1/users?q=<query>`

**推奨対策**:

```go
// ワイルドカード文字のエスケープ
func escapeWildcards(s string) string {
    s = strings.ReplaceAll(s, "%", "\\%")
    s = strings.ReplaceAll(s, "_", "\\_")
    return s
}

// 使用例
searchPattern := "%" + escapeWildcards(query) + "%"
```

---

### 🟡 Medium (中優先度)

#### 7. トークンのリフレッシュ機構がない

**場所**: 認証システム全体

**問題**:
- JWTの有効期限が7日間（168時間）と長い
- リフレッシュトークンの仕組みがない
- トークンの無効化（ブラックリスト）機能がない

**リスク**:
- トークン漏洩時の影響範囲が大きい
- ログアウト後もトークンが有効（セッション管理ができない）
- セキュリティとUXのトレードオフが不適切

**影響範囲**: 全認証システム

**推奨対策**:

1. **アクセストークンの有効期限を短縮** (例: 15分)
2. **リフレッシュトークンの導入**
```go
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// アクセストークン: 15分
// リフレッシュトークン: 7日間
```

3. **Redisでトークンブラックリストを管理**

---

#### 8. パスワードリセット機能がない

**場所**: 認証システム

**問題**:
- パスワードを忘れた場合のリカバリー手段がない
- アカウントロックアウト時の対処方法がない

**リスク**:
- ユーザビリティの低下
- サポートコストの増加
- アカウント乗っ取り時の復旧が困難

**影響範囲**: ユーザー管理全体

**推奨対策**: Phase 2以降で実装を検討

---

#### 9. ログイン試行回数の制限がない

**場所**: `/api/v1/auth/login`

**問題**:
- 同一アカウントへの連続ログイン試行を制限していない
- アカウントロックアウト機能がない

**リスク**:
- ブルートフォース攻撃（パスワード総当たり）
- 特定アカウントへの集中攻撃

**影響範囲**: ログインエンドポイント

**推奨対策**:

```go
// Redisでログイン失敗回数を管理
// 5回失敗で15分間ロック
type LoginAttempt struct {
    Email      string
    FailCount  int
    LockedUntil time.Time
}
```

---

#### 10. セキュリティヘッダーの不足

**場所**: `backend/cmd/server/main.go`

**問題**:
- `X-Content-Type-Options` ヘッダーがない
- `X-Frame-Options` ヘッダーがない
- `Content-Security-Policy` ヘッダーがない
- `Strict-Transport-Security` ヘッダーがない

**リスク**:
- クリックジャッキング攻撃
- MIME type sniffing
- XSS攻撃のリスク増加

**影響範囲**: 全HTTPレスポンス

**推奨対策**:

```go
e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
    XSSProtection:         "1; mode=block",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "DENY",
    HSTSMaxAge:            31536000,
    ContentSecurityPolicy: "default-src 'self'",
}))
```

---

#### 11. 機密情報のログ出力

**場所**: 全体（潜在的リスク）

**問題**:
- パスワード、トークンなどの機密情報がログに出力される可能性
- エラーログに詳細なスタックトレースが含まれる可能性

**リスク**:
- ログファイルからの機密情報漏洩
- 攻撃者へのシステム内部情報の提供

**影響範囲**: ログ管理全体

**推奨対策**:

1. **構造化ログの導入**
```go
import "github.com/sirupsen/logrus"

// 機密情報をマスク
log.WithFields(logrus.Fields{
    "email": maskEmail(user.Email),
    "ip":    c.RealIP(),
}).Info("User login attempt")
```

2. **本番環境ではDEBUGログを無効化**

---

### 🔵 Low (低優先度)

#### 12. フロントエンドのトークン保存方法

**場所**: `frontend/src/utils/storage.ts:4-14`

**問題**:
```typescript
export const tokenStorage = {
  get: (): string | null => {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
  },
  set: (token: string): void => {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token);
  },
```

**リスク**:
- `localStorage`はXSS攻撃に脆弱
- JavaScriptからアクセス可能なため、悪意あるスクリプトによるトークン盗難の危険性
- より安全な代替手段が存在

**影響範囲**: フロントエンド認証

**推奨対策**:

1. **HttpOnly Cookieの使用** (最も安全)
```go
// バックエンド
c.SetCookie(&http.Cookie{
    Name:     "access_token",
    Value:    token,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode,
    MaxAge:   900, // 15分
})
```

2. **またはsessionStorageの使用** (次善策)
```typescript
sessionStorage.setItem(STORAGE_KEYS.TOKEN, token);
```

---

#### 13. エラーメッセージが詳細すぎる

**場所**:
- `backend/internal/handler/auth_handler.go:82`
- `backend/internal/service/auth_service.go:94`

**問題**:
```go
return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
```

**現状**: 良好（メールアドレスとパスワードを区別していない）

**リスク**: 低（現状問題なし）

**改善点**:
- ログイン失敗時のレスポンス時間を一定にする（タイミング攻撃対策）

```go
// タイミング攻撃対策
time.Sleep(time.Millisecond * 200)
return c.JSON(http.StatusUnauthorized, errors.Unauthorized("メールアドレスまたはパスワードが正しくありません"))
```

---

#### 14. 環境変数の管理

**場所**: `docker-compose.yml:44`

**問題**:
```yaml
JWT_SECRET: your-secret-key-change-in-production
```

**リスク**:
- 本番環境で変更し忘れる危険性
- docker-compose.ymlがGitにコミットされている（`.gitignore`で除外されていない）

**影響範囲**: デプロイメント

**推奨対策**:

1. **.env.exampleを作成し、実際の`.env`ファイルは`.gitignore`に追加**
```bash
# .gitignore
.env
docker-compose.override.yml
```

2. **本番環境では環境変数を外部から注入**
```yaml
# docker-compose.prod.yml
environment:
  JWT_SECRET: ${JWT_SECRET}  # 環境変数から取得
```

---

## 推奨される対策

### 優先度1: 即時対応（Critical）

1. **JWT秘密鍵の強制チェック** (影響度: ★★★★★)
   - デフォルト値の削除
   - 起動時の環境変数検証
   - 本番用の強力な秘密鍵生成

2. **CORS設定の厳格化** (影響度: ★★★★★)
   - オリジン制限の実装
   - 認証情報の適切な管理
   - Preflight requestの適切な処理

### 優先度2: 1週間以内（High）

3. **レート制限の実装** (影響度: ★★★★☆)
   - 認証エンドポイント: 5分/5回
   - 投稿作成: 1時間/30回
   - コメント: 1時間/100回

4. **バリデーションの強化** (影響度: ★★★★☆)
   - validator/v10の導入
   - パスワード強度チェック
   - メールアドレス形式検証

5. **ユーザー列挙対策** (影響度: ★★★☆☆)
   - エラーメッセージの一般化

6. **SQLワイルドカードエスケープ** (影響度: ★★★☆☆)

### 優先度3: 1ヶ月以内（Medium）

7. **トークンリフレッシュ機構** (Phase 2で実装検討)
8. **パスワードリセット機能** (Phase 2で実装検討)
9. **ログイン試行制限**
10. **セキュリティヘッダーの追加**
11. **ログマスキング機能**

### 優先度4: 適宜対応（Low）

12. **HttpOnly Cookie検討**
13. **タイミング攻撃対策**
14. **環境変数管理の改善**

---

## セキュアな実装の確認事項

### ✅ 良好な実装

以下の点は適切に実装されています：

1. **パスワードハッシュ化**
   - bcryptを使用（`bcrypt.DefaultCost = 10`）
   - `backend/internal/service/auth_service.go:56`

2. **SQLインジェクション対策**
   - GORMのプレースホルダを適切に使用
   - 直接的なSQLクエリは使用していない

3. **XSS対策（基本）**
   - Reactのデフォルトエスケープを活用
   - `dangerouslySetInnerHTML`の使用なし

4. **認証ミドルウェア**
   - JWT検証が適切に実装
   - `backend/internal/middleware/auth_middleware.go`

5. **論理削除**
   - GORMの論理削除を使用
   - データの完全性が保たれる

6. **HTTPS対応準備**
   - HTTPSを前提とした設計（本番環境で必須）

---

## 次のステップ

### Phase 2 開発前の必須対応

- [ ] JWT秘密鍵のハードコード削除
- [ ] CORS設定の厳格化
- [ ] レート制限の実装
- [ ] バリデーションの強化

### Phase 2 で実装検討

- [ ] トークンリフレッシュ機構
- [ ] パスワードリセット機能
- [ ] 二要素認証（2FA）
- [ ] HttpOnly Cookieへの移行
- [ ] セキュリティログの強化

### 継続的なモニタリング

- [ ] セキュリティヘッダーのスキャン
- [ ] 依存関係の脆弱性チェック（`npm audit`, `go list -m all`)
- [ ] ペネトレーションテスト実施
- [ ] OWASP Top 10 準拠確認

---

## 参考資料

- [OWASP Top 10 (2021)](https://owasp.org/www-project-top-ten/)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [Go Secure Coding Practices](https://github.com/OWASP/Go-SCP)
- [React Security Best Practices](https://snyk.io/blog/10-react-security-best-practices/)
- [JWT Security Best Practices](https://tools.ietf.org/html/rfc8725)

---

**レポート作成者**: Claude Code
**最終更新**: 2026-03-04
