# API仕様書

## ベースURL

### 開発環境
```
http://localhost:8080/api/v1
```

### 本番環境
```
https://your-app.render.com/api/v1
```

## 認証

### 認証方式
JWT (JSON Web Token) を使用

### 認証ヘッダー
```
Authorization: Bearer <token>
```

### 認証が不要なエンドポイント
- `POST /auth/register`
- `POST /auth/login`
- `GET /posts` (公開投稿の閲覧)
- `GET /users/:username` (公開プロフィールの閲覧)

## レスポンス形式

### 成功レスポンス
```json
{
  "data": { ... },
  "message": "Success"
}
```

### エラーレスポンス
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": [ ... ]
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

## Phase 1 API エンドポイント

### 認証 (Authentication)

#### ユーザー登録
```
POST /auth/register
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "username": "johndoe",
  "display_name": "John Doe"
}
```
**Response:** `201 Created`
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe",
      "display_name": "John Doe",
      "bio": "",
      "created_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

#### ログイン
```
POST /auth/login
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
**Response:** `200 OK`
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { ... }
  }
}
```

#### 現在のユーザー情報取得
```
GET /auth/me
Authorization: Bearer <token>
```
**Response:** `200 OK`
```json
{
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "display_name": "John Doe",
    "bio": "Hello!",
    "follower_count": 100,
    "following_count": 50,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### ユーザー (Users)

#### ユーザープロフィール取得
```
GET /users/:username
```
**Response:** `200 OK`
```json
{
  "data": {
    "id": "uuid",
    "username": "johndoe",
    "display_name": "John Doe",
    "bio": "Hello!",
    "avatar_url": null,
    "header_url": null,
    "follower_count": 100,
    "following_count": 50,
    "is_following": false,
    "is_followed_by": false,
    "is_blocked": false,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### ユーザープロフィール更新
```
PATCH /users/me
Authorization: Bearer <token>
```
**Request Body:**
```json
{
  "display_name": "John Doe Updated",
  "bio": "New bio"
}
```
**Response:** `200 OK`

#### ユーザー検索
```
GET /users?q=john&limit=20&offset=0
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "username": "johndoe",
      "display_name": "John Doe",
      "avatar_url": null,
      "follower_count": 100
    }
  ],
  "pagination": {
    "total": 5,
    "limit": 20,
    "offset": 0,
    "has_more": false
  }
}
```

---

### 投稿 (Posts)

#### 投稿作成
```
POST /posts
Authorization: Bearer <token>
```
**Request Body (Phase 1: テキストのみ):**
```json
{
  "content": "Hello, world!"
}
```
**Request Body (Phase 2: メディアあり):**
```json
{
  "content": "Check this out!",
  "media": [
    {
      "media_type": "image",
      "media_url": "https://storage.firebase.com/..."
    }
  ],
  "visibility": "public"
}
```
**Response:** `201 Created`
```json
{
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "content": "Hello, world!",
    "visibility": "public",
    "like_count": 0,
    "comment_count": 0,
    "retweet_count": 0,
    "is_liked": false,
    "is_bookmarked": false,
    "user": { ... },
    "media": [],
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### 投稿取得
```
GET /posts/:id
```
**Response:** `200 OK`

#### 投稿更新
```
PATCH /posts/:id
Authorization: Bearer <token>
```
**Request Body:**
```json
{
  "content": "Updated content"
}
```
**Response:** `200 OK`

#### 投稿削除
```
DELETE /posts/:id
Authorization: Bearer <token>
```
**Response:** `204 No Content`

#### タイムライン取得
```
GET /timeline?limit=20&offset=0
Authorization: Bearer <token>
```
フォロー中のユーザー + 自分の投稿を時系列で取得

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "user": { ... },
      "content": "Post content",
      "media": [],
      "like_count": 10,
      "comment_count": 5,
      "is_liked": true,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ],
  "pagination": { ... }
}
```

#### ユーザーの投稿一覧取得
```
GET /users/:username/posts?limit=20&offset=0
```
**Response:** `200 OK`

---

### コメント (Comments)

#### コメント作成
```
POST /posts/:id/comments
Authorization: Bearer <token>
```
**Request Body:**
```json
{
  "content": "Great post!",
  "parent_comment_id": null
}
```
**Response:** `201 Created`

#### コメント取得
```
GET /posts/:id/comments?limit=20&offset=0
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "post_id": "uuid",
      "user": { ... },
      "content": "Great post!",
      "parent_comment_id": null,
      "child_comments": [
        {
          "id": "uuid",
          "content": "I agree!",
          "user": { ... },
          "like_count": 2,
          "created_at": "2025-01-01T00:00:00Z"
        }
      ],
      "like_count": 5,
      "is_liked": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### コメント削除
```
DELETE /comments/:id
Authorization: Bearer <token>
```
**Response:** `204 No Content`

---

### いいね (Likes)

#### 投稿にいいね
```
POST /posts/:id/like
Authorization: Bearer <token>
```
**Response:** `201 Created`

#### 投稿のいいね解除
```
DELETE /posts/:id/like
Authorization: Bearer <token>
```
**Response:** `204 No Content`

#### コメントにいいね
```
POST /comments/:id/like
Authorization: Bearer <token>
```
**Response:** `201 Created`

#### コメントのいいね解除
```
DELETE /comments/:id/like
Authorization: Bearer <token>
```
**Response:** `204 No Content`

---

### ブックマーク (Bookmarks)

#### ブックマーク追加
```
POST /posts/:id/bookmark
Authorization: Bearer <token>
```
**Response:** `201 Created`

#### ブックマーク解除
```
DELETE /posts/:id/bookmark
Authorization: Bearer <token>
```
**Response:** `204 No Content`

#### ブックマーク一覧取得
```
GET /bookmarks?limit=20&offset=0
Authorization: Bearer <token>
```
**Response:** `200 OK`

---

### フォロー (Follows)

#### フォロー
```
POST /users/:username/follow
Authorization: Bearer <token>
```
**Response:** `201 Created`

#### フォロー解除
```
DELETE /users/:username/follow
Authorization: Bearer <token>
```
**Response:** `204 No Content`

#### フォロワー一覧取得
```
GET /users/:username/followers?limit=20&offset=0
```
**Response:** `200 OK`

#### フォロー中一覧取得
```
GET /users/:username/following?limit=20&offset=0
```
**Response:** `200 OK`

---

### 通知 (Notifications)

#### 通知一覧取得
```
GET /notifications?limit=20&offset=0
Authorization: Bearer <token>
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "notification_type": "like",
      "actor": {
        "id": "uuid",
        "username": "janedoe",
        "display_name": "Jane Doe",
        "avatar_url": null
      },
      "target_type": "Post",
      "target_id": "uuid",
      "is_read": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ],
  "unread_count": 5
}
```

#### 通知を既読にする
```
PATCH /notifications/:id/read
Authorization: Bearer <token>
```
**Response:** `200 OK`

#### すべての通知を既読にする
```
POST /notifications/read-all
Authorization: Bearer <token>
```
**Response:** `200 OK`

---

### 通報 (Reports)

#### 通報作成
```
POST /reports
Authorization: Bearer <token>
```
**Request Body:**
```json
{
  "reportable_type": "Post",
  "reportable_id": "uuid",
  "reason": "spam",
  "details": "This is spam content"
}
```
**Response:** `201 Created`

**通報理由 (reason):**
- `spam`: スパム
- `inappropriate`: 不適切なコンテンツ
- `harassment`: ハラスメント
- `other`: その他

---

## Phase 2 追加API

### リツイート (Retweets)
```
POST /posts/:id/retweet
DELETE /posts/:id/retweet
GET /posts/:id/retweets?limit=20&offset=0
```

### ハッシュタグ (Hashtags)
```
GET /hashtags/trending?limit=10
GET /hashtags/:name/posts?limit=20&offset=0
```

### 投稿検索 (Search)
```
GET /posts/search?q=keyword&limit=20&offset=0
```

### ブロック (Blocks)
```
POST /users/:username/block
DELETE /users/:username/block
GET /blocks?limit=20&offset=0
```

### メディアアップロード (File Upload)
```
POST /upload/image
POST /upload/video
```

---

## レート制限

- **認証なし**: 100リクエスト/分
- **認証あり**: 500リクエスト/分

レート制限超過時: `429 Too Many Requests`

## OpenAPI仕様書

詳細なAPI仕様は swaggo により自動生成される OpenAPI 定義ファイルを参照:
- `backend/docs/swagger.json`
- `backend/docs/swagger.yaml`

Swagger UI: `http://localhost:8080/swagger/index.html`
