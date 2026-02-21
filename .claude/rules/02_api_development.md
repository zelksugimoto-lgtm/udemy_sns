# API開発ルール

## OpenAPI仕様優先開発

### 開発フロー
1. **OpenAPI仕様書作成** (swaggoアノテーション)
2. **swagger生成** (`swag init`)
3. **型生成** (フロントエンド: openapi-typescript)
4. **実装** (バックエンド → フロントエンド)

### swaggoアノテーション必須
すべてのエンドポイントに以下を記述：
```go
// @Summary      ユーザー登録
// @Description  新しいユーザーを登録
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "登録情報"
// @Success      201      {object}  dto.AuthResponse
// @Failure      400      {object}  errors.ErrorResponse
// @Failure      409      {object}  errors.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
    // ...
}
```

---

## エンドポイント設計

### URL構造
```
/api/v1/{resource}[/{id}][/{action}]
```

**例:**
- `POST /api/v1/posts` - 投稿作成
- `GET /api/v1/posts/:id` - 投稿取得
- `PATCH /api/v1/posts/:id` - 投稿更新
- `DELETE /api/v1/posts/:id` - 投稿削除
- `POST /api/v1/posts/:id/like` - 投稿にいいね
- `GET /api/v1/timeline` - タイムライン取得

### HTTPメソッド
- `GET`: リソース取得
- `POST`: リソース作成、アクション実行
- `PATCH`: リソース部分更新
- `DELETE`: リソース削除

---

## レスポンス形式

### 成功レスポンス
```json
{
  "data": { ... }
}
```

### エラーレスポンス
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

### ページネーションレスポンス
```json
{
  "data": [ ... ],
  "pagination": {
    "total": 100,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

---

## HTTPステータスコード

### 成功
- `200 OK`: 取得・更新成功
- `201 Created`: 作成成功
- `204 No Content`: 削除成功

### クライアントエラー
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証エラー（トークンなし、無効）
- `403 Forbidden`: 権限エラー（認証済みだが権限なし）
- `404 Not Found`: リソースが存在しない
- `409 Conflict`: 重複エラー（メール、ユーザー名等）
- `429 Too Many Requests`: レート制限

### サーバーエラー
- `500 Internal Server Error`: サーバーエラー

---

## 認証・認可

### JWT認証
```go
// 認証が必要なエンドポイント
api.Use(middleware.AuthMiddleware(jwtSecret))
api.POST("/posts", postHandler.Create)

// 認証が不要なエンドポイント
public.POST("/auth/login", authHandler.Login)
public.POST("/auth/register", authHandler.Register)
```

### ユーザー情報取得
```go
func (h *PostHandler) Create(c echo.Context) error {
    // ミドルウェアで設定されたユーザー情報を取得
    user := c.Get("user").(*model.User)
    // ...
}
```

---

## バリデーション

### リクエストバリデーション
```go
import "github.com/go-playground/validator/v10"

type CreatePostRequest struct {
    Content string `json:"content" validate:"required,min=1,max=250"`
}

func (h *PostHandler) Create(c echo.Context) error {
    req := new(CreatePostRequest)

    if err := c.Bind(req); err != nil {
        return c.JSON(400, errors.BadRequest("Invalid request"))
    }

    if err := h.validator.Struct(req); err != nil {
        return c.JSON(400, errors.ValidationError(err))
    }

    // ...
}
```

### バリデーションルール
- **必須**: `required`
- **文字数**: `min=1,max=250`
- **メール**: `email`
- **URL**: `url`

---

## ページネーション

### パラメータ
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `offset`: 開始位置（デフォルト: 0）

### 実装例
```go
func (h *PostHandler) GetTimeline(c echo.Context) error {
    limit := getIntParam(c, "limit", 20, 100)
    offset := getIntParam(c, "offset", 0, math.MaxInt)

    posts, total, err := h.service.GetTimeline(user.ID, limit, offset)
    if err != nil {
        return c.JSON(500, errors.InternalError(err))
    }

    return c.JSON(200, map[string]interface{}{
        "data": posts,
        "pagination": map[string]interface{}{
            "total":    total,
            "limit":    limit,
            "offset":   offset,
            "has_more": offset+limit < total,
        },
    })
}
```

---

## エラーハンドリング

### カスタムエラー
```go
// pkg/errors/errors.go
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Details []ValidationError `json:"details,omitempty"`
}

func BadRequest(message string) *ErrorResponse {
    return &ErrorResponse{
        Error: ErrorDetail{
            Code:    "BAD_REQUEST",
            Message: message,
        },
    }
}

func Unauthorized(message string) *ErrorResponse {
    return &ErrorResponse{
        Error: ErrorDetail{
            Code:    "UNAUTHORIZED",
            Message: message,
        },
    }
}
```

### エラーログ
```go
import "github.com/labstack/echo/v4"

func (h *PostHandler) Create(c echo.Context) error {
    // ...

    if err != nil {
        c.Logger().Error("Failed to create post", err)
        return c.JSON(500, errors.InternalError("Failed to create post"))
    }

    // ...
}
```

---

## ミドルウェア

### 認証ミドルウェア
```go
func AuthMiddleware(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            token := extractToken(c)
            claims, err := auth.ValidateToken(token, secret)
            if err != nil {
                return c.JSON(401, errors.Unauthorized("Invalid token"))
            }

            user, err := getUserByID(claims.UserID)
            if err != nil {
                return c.JSON(401, errors.Unauthorized("User not found"))
            }

            c.Set("user", user)
            return next(c)
        }
    }
}
```

### CORS ミドルウェア
```go
import "github.com/labstack/echo/v4/middleware"

e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000", "https://your-app.web.app"},
    AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
}))
```

### レート制限ミドルウェア
```go
e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))
```

---

## レイヤー構成

### Handler (コントローラー)
- HTTPリクエスト/レスポンス処理
- バリデーション
- Serviceの呼び出し

### Service (ビジネスロジック)
- ビジネスルールの実装
- Repositoryの呼び出し
- トランザクション管理

### Repository (データアクセス)
- データベースCRUD操作
- クエリ実装

### 依存関係
```
Handler → Service → Repository → Database
```

---

## テスト

### ユニットテスト
```go
func TestCreatePost(t *testing.T) {
    // Setup
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(`{"content":"test"}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Test
    handler := NewPostHandler(mockService)
    err := handler.Create(c)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, rec.Code)
}
```

---

## OpenAPI生成

### swag init 実行
```bash
cd backend
swag init -g cmd/server/main.go -o docs
```

### Swagger UI アクセス
```
http://localhost:8080/swagger/index.html
```

---

## 環境変数

```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=social_media_db

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=168h

# Firebase
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_STORAGE_BUCKET=your-bucket

# File Upload
MAX_VIDEO_SIZE_MB=20
MAX_IMAGE_SIZE_MB=5
MAX_UPLOAD_FILES=4
```

---

## 参照
- Echo: https://echo.labstack.com/
- swaggo: https://github.com/swaggo/swag
- validator: https://github.com/go-playground/validator
