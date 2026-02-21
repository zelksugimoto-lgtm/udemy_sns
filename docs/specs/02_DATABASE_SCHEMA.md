# データベーススキーマ設計

## 概要
PostgreSQLを使用したリレーショナルデータベース設計

## 基本ルール
1. **主キー**: すべてのテーブルは `id` (UUID) を主キーとする
2. **論理削除**: `deleted_at` フィールドで論理削除を管理
3. **タイムスタンプ**: すべてのテーブルに `created_at`, `updated_at` を含める
4. **外部キー制約**: リレーションシップには外部キー制約を設定
5. **インデックス**: 検索に使用されるカラムにはインデックスを設定

## 共通フィールド (BaseModel)

すべてのテーブルに含まれる共通フィールド：

```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    CreatedAt time.Time      `gorm:"not null" json:"created_at"`
    UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

## テーブル一覧

### Phase 1
1. `users` - ユーザー
2. `posts` - 投稿
3. `post_media` - 投稿メディア（Phase 2で使用）
4. `comments` - コメント
5. `likes` - いいね（投稿・コメント共通）
6. `bookmarks` - ブックマーク
7. `follows` - フォロー関係
8. `notifications` - 通知
9. `reports` - 通報

### Phase 2
10. `retweets` - リツイート
11. `hashtags` - ハッシュタグ
12. `post_hashtags` - 投稿とハッシュタグの関連
13. `blocks` - ブロック

### Phase 3
14. `admin_actions` - 管理者アクション
15. `post_analytics` - 投稿分析

---

## テーブル定義

### 1. users テーブル

ユーザー情報を管理

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    avatar_url VARCHAR(500),        -- Phase 2
    header_url VARCHAR(500),        -- Phase 2
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT username_format CHECK (username ~ '^[a-zA-Z0-9_]{3,50}$')
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

**フィールド説明:**
- `email`: メールアドレス（ログイン用、一意）
- `username`: ユーザー名（@username形式、一意、3-50文字、英数字とアンダースコアのみ）
- `display_name`: 表示名（1-100文字）
- `password_hash`: bcryptハッシュ化されたパスワード
- `bio`: 自己紹介（最大1000文字）
- `avatar_url`: プロフィール画像URL（Firebase Storage）
- `header_url`: ヘッダー画像URL（Firebase Storage）
- `is_admin`: 管理者フラグ

**GORM Model:**
```go
type User struct {
    BaseModel
    Email        string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email,where:deleted_at IS NULL" json:"email"`
    Username     string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_users_username,where:deleted_at IS NULL" json:"username"`
    DisplayName  string         `gorm:"type:varchar(100);not null" json:"display_name"`
    PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
    Bio          string         `gorm:"type:text" json:"bio"`
    AvatarURL    string         `gorm:"type:varchar(500)" json:"avatar_url"`
    HeaderURL    string         `gorm:"type:varchar(500)" json:"header_url"`
    IsAdmin      bool           `gorm:"default:false" json:"is_admin"`

    // Relations
    Posts         []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
    Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
    Likes         []Like         `gorm:"foreignKey:UserID" json:"likes,omitempty"`
    Bookmarks     []Bookmark     `gorm:"foreignKey:UserID" json:"bookmarks,omitempty"`
    Following     []Follow       `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
    Followers     []Follow       `gorm:"foreignKey:FollowedID" json:"followers,omitempty"`
    Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}
```

---

### 2. posts テーブル

投稿情報を管理

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',  -- Phase 2: public, followers, private
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT content_length CHECK (LENGTH(content) <= 250),
    CONSTRAINT visibility_check CHECK (visibility IN ('public', 'followers', 'private'))
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);
CREATE INDEX idx_posts_visibility ON posts(visibility);
```

**フィールド説明:**
- `user_id`: 投稿者のユーザーID
- `content`: 投稿本文（最大250文字）
- `visibility`: 公開範囲（public: 全体公開, followers: フォロワーのみ, private: 自分のみ）

**GORM Model:**
```go
type Post struct {
    BaseModel
    UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    Content    string    `gorm:"type:text;not null" json:"content"`
    Visibility string    `gorm:"type:varchar(20);not null;default:'public';index" json:"visibility"`

    // Relations
    User       User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Media      []PostMedia  `gorm:"foreignKey:PostID" json:"media,omitempty"`
    Comments   []Comment    `gorm:"foreignKey:PostID" json:"comments,omitempty"`
    Likes      []Like       `gorm:"foreignKey:PostID;polymorphic:Likeable" json:"likes,omitempty"`
    Bookmarks  []Bookmark   `gorm:"foreignKey:PostID" json:"bookmarks,omitempty"`
    Retweets   []Retweet    `gorm:"foreignKey:PostID" json:"retweets,omitempty"`
    Hashtags   []Hashtag    `gorm:"many2many:post_hashtags" json:"hashtags,omitempty"`

    // Computed fields (not in DB, populated by queries)
    LikeCount     int  `gorm:"-" json:"like_count"`
    CommentCount  int  `gorm:"-" json:"comment_count"`
    RetweetCount  int  `gorm:"-" json:"retweet_count"`
    IsLiked       bool `gorm:"-" json:"is_liked"`
    IsBookmarked  bool `gorm:"-" json:"is_bookmarked"`
    IsRetweeted   bool `gorm:"-" json:"is_retweeted"`
}
```

---

### 3. post_media テーブル

投稿に添付されたメディア（画像/動画）を管理（Phase 2）

```sql
CREATE TABLE post_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    media_type VARCHAR(20) NOT NULL,  -- image, video
    media_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),       -- Phase 4: サムネイル
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT media_type_check CHECK (media_type IN ('image', 'video')),
    CONSTRAINT max_media_per_post CHECK (display_order < 4)
);

CREATE INDEX idx_post_media_post_id ON post_media(post_id);
CREATE INDEX idx_post_media_display_order ON post_media(post_id, display_order);
```

**フィールド説明:**
- `post_id`: 投稿ID
- `media_type`: メディアタイプ（image or video）
- `media_url`: メディアのURL（Firebase Storage）
- `thumbnail_url`: サムネイルURL（動画用、Phase 4）
- `display_order`: 表示順序（0-3、最大4枚）

**GORM Model:**
```go
type PostMedia struct {
    BaseModel
    PostID       uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
    MediaType    string    `gorm:"type:varchar(20);not null" json:"media_type"`
    MediaURL     string    `gorm:"type:varchar(500);not null" json:"media_url"`
    ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url"`
    DisplayOrder int       `gorm:"not null;default:0" json:"display_order"`

    // Relations
    Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
```

---

### 4. comments テーブル

コメント（返信）を管理、ネスト構造対応

```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT content_length CHECK (LENGTH(content) <= 500)
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);
```

**フィールド説明:**
- `post_id`: コメント対象の投稿ID
- `user_id`: コメント投稿者のユーザーID
- `parent_comment_id`: 親コメントID（ネスト用、NULL = トップレベルコメント）
- `content`: コメント本文（最大500文字）

**GORM Model:**
```go
type Comment struct {
    BaseModel
    PostID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"post_id"`
    UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
    ParentCommentID *uuid.UUID `gorm:"type:uuid;index" json:"parent_comment_id,omitempty"`
    Content         string     `gorm:"type:text;not null" json:"content"`

    // Relations
    Post           Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
    User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    ParentComment  *Comment  `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
    ChildComments  []Comment `gorm:"foreignKey:ParentCommentID" json:"child_comments,omitempty"`
    Likes          []Like    `gorm:"foreignKey:CommentID;polymorphic:Likeable" json:"likes,omitempty"`

    // Computed fields
    LikeCount int  `gorm:"-" json:"like_count"`
    IsLiked   bool `gorm:"-" json:"is_liked"`
}
```

---

### 5. likes テーブル

いいね（投稿・コメント共通、ポリモーフィック）

```sql
CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    likeable_type VARCHAR(50) NOT NULL,  -- Post, Comment
    likeable_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT likeable_type_check CHECK (likeable_type IN ('Post', 'Comment')),
    CONSTRAINT unique_like UNIQUE (user_id, likeable_type, likeable_id)
);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_likeable ON likes(likeable_type, likeable_id);
```

**フィールド説明:**
- `user_id`: いいねしたユーザーID
- `likeable_type`: いいね対象のタイプ（Post or Comment）
- `likeable_id`: いいね対象のID
- **制約**: 同一ユーザーが同一対象に複数回いいねできない

**GORM Model:**
```go
type Like struct {
    BaseModel
    UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    LikeableType string    `gorm:"type:varchar(50);not null;index:idx_likes_likeable" json:"likeable_type"`
    LikeableID   uuid.UUID `gorm:"type:uuid;not null;index:idx_likes_likeable" json:"likeable_id"`

    // Relations
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

---

### 6. bookmarks テーブル

ブックマーク

```sql
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT unique_bookmark UNIQUE (user_id, post_id)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_post_id ON bookmarks(post_id);
CREATE INDEX idx_bookmarks_created_at ON bookmarks(created_at DESC);
```

**フィールド説明:**
- `user_id`: ブックマークしたユーザーID
- `post_id`: ブックマーク対象の投稿ID
- **制約**: 同一ユーザーが同一投稿に複数回ブックマークできない

**GORM Model:**
```go
type Bookmark struct {
    BaseModel
    UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    PostID uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`

    // Relations
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
```

---

### 7. follows テーブル

フォロー関係

```sql
CREATE TABLE follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT no_self_follow CHECK (follower_id != followed_id),
    CONSTRAINT unique_follow UNIQUE (follower_id, followed_id)
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_followed_id ON follows(followed_id);
CREATE INDEX idx_follows_created_at ON follows(created_at DESC);
```

**フィールド説明:**
- `follower_id`: フォローしているユーザーID
- `followed_id`: フォローされているユーザーID
- **制約**:
  - 自分自身をフォローできない
  - 同一ユーザーペアに複数回フォローできない

**GORM Model:**
```go
type Follow struct {
    BaseModel
    FollowerID uuid.UUID `gorm:"type:uuid;not null;index" json:"follower_id"`
    FollowedID uuid.UUID `gorm:"type:uuid;not null;index" json:"followed_id"`

    // Relations
    Follower User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
    Followed User `gorm:"foreignKey:FollowedID" json:"followed,omitempty"`
}
```

---

### 8. notifications テーブル

通知

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL,
    target_type VARCHAR(50),
    target_id UUID,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT notification_type_check CHECK (notification_type IN ('like', 'comment', 'follow', 'retweet'))
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
```

**フィールド説明:**
- `user_id`: 通知を受け取るユーザーID
- `actor_id`: 通知のアクションを起こしたユーザーID
- `notification_type`: 通知タイプ（like, comment, follow, retweet）
- `target_type`: 対象のタイプ（Post, Comment等）
- `target_id`: 対象のID
- `is_read`: 既読フラグ

**通知タイプ別の値:**
| Type | Target Type | Target ID | 意味 |
|------|-------------|-----------|------|
| like | Post | post_id | 投稿にいいね |
| like | Comment | comment_id | コメントにいいね |
| comment | Post | post_id | 投稿にコメント |
| comment | Comment | comment_id | コメントに返信 |
| follow | NULL | NULL | フォロー |
| retweet | Post | post_id | リツイート（Phase 2） |

**GORM Model:**
```go
type Notification struct {
    BaseModel
    UserID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
    ActorID          uuid.UUID  `gorm:"type:uuid;not null" json:"actor_id"`
    NotificationType string     `gorm:"type:varchar(50);not null" json:"notification_type"`
    TargetType       string     `gorm:"type:varchar(50)" json:"target_type,omitempty"`
    TargetID         *uuid.UUID `gorm:"type:uuid" json:"target_id,omitempty"`
    IsRead           bool       `gorm:"default:false;index:idx_notifications_is_read" json:"is_read"`

    // Relations
    User  User `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Actor User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}
```

---

### 9. reports テーブル

通報

```sql
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reportable_type VARCHAR(50) NOT NULL,
    reportable_id UUID NOT NULL,
    reason VARCHAR(100) NOT NULL,
    details TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT reportable_type_check CHECK (reportable_type IN ('Post', 'Comment', 'User')),
    CONSTRAINT status_check CHECK (status IN ('pending', 'reviewing', 'resolved', 'dismissed'))
);

CREATE INDEX idx_reports_reporter_id ON reports(reporter_id);
CREATE INDEX idx_reports_reportable ON reports(reportable_type, reportable_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_created_at ON reports(created_at DESC);
```

**フィールド説明:**
- `reporter_id`: 通報者のユーザーID
- `reportable_type`: 通報対象のタイプ（Post, Comment, User）
- `reportable_id`: 通報対象のID
- `reason`: 通報理由（選択肢: spam, inappropriate, harassment, other）
- `details`: 詳細説明（自由入力、任意）
- `status`: ステータス（pending, reviewing, resolved, dismissed）
- `resolved_by`: 対応した管理者のユーザーID（Phase 3）
- `resolved_at`: 解決日時（Phase 3）

**GORM Model:**
```go
type Report struct {
    BaseModel
    ReporterID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"reporter_id"`
    ReportableType string     `gorm:"type:varchar(50);not null;index:idx_reports_reportable" json:"reportable_type"`
    ReportableID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_reports_reportable" json:"reportable_id"`
    Reason         string     `gorm:"type:varchar(100);not null" json:"reason"`
    Details        string     `gorm:"type:text" json:"details,omitempty"`
    Status         string     `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
    ResolvedBy     *uuid.UUID `gorm:"type:uuid" json:"resolved_by,omitempty"`
    ResolvedAt     *time.Time `json:"resolved_at,omitempty"`

    // Relations
    Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
    Resolver *User `gorm:"foreignKey:ResolvedBy" json:"resolver,omitempty"`
}
```

---

## Phase 2 追加テーブル

### 10. retweets テーブル

リツイート（シェア）

```sql
CREATE TABLE retweets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT unique_retweet UNIQUE (user_id, post_id)
);

CREATE INDEX idx_retweets_user_id ON retweets(user_id);
CREATE INDEX idx_retweets_post_id ON retweets(post_id);
CREATE INDEX idx_retweets_created_at ON retweets(created_at DESC);
```

**GORM Model:**
```go
type Retweet struct {
    BaseModel
    UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    PostID uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`

    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
```

---

### 11. hashtags テーブル

ハッシュタグ

```sql
CREATE TABLE hashtags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT name_format CHECK (name ~ '^[a-zA-Z0-9_]+$')
);

CREATE INDEX idx_hashtags_name ON hashtags(name);
```

**GORM Model:**
```go
type Hashtag struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
    CreatedAt time.Time `gorm:"not null" json:"created_at"`
    UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

    Posts []Post `gorm:"many2many:post_hashtags" json:"posts,omitempty"`
}
```

---

### 12. post_hashtags テーブル

投稿とハッシュタグの多対多関連

```sql
CREATE TABLE post_hashtags (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    hashtag_id UUID NOT NULL REFERENCES hashtags(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (post_id, hashtag_id)
);

CREATE INDEX idx_post_hashtags_hashtag_id ON post_hashtags(hashtag_id);
```

---

### 13. blocks テーブル

ブロック

```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT no_self_block CHECK (blocker_id != blocked_id),
    CONSTRAINT unique_block UNIQUE (blocker_id, blocked_id)
);

CREATE INDEX idx_blocks_blocker_id ON blocks(blocker_id);
CREATE INDEX idx_blocks_blocked_id ON blocks(blocked_id);
```

**GORM Model:**
```go
type Block struct {
    BaseModel
    BlockerID uuid.UUID `gorm:"type:uuid;not null;index" json:"blocker_id"`
    BlockedID uuid.UUID `gorm:"type:uuid;not null;index" json:"blocked_id"`

    Blocker User `gorm:"foreignKey:BlockerID" json:"blocker,omitempty"`
    Blocked User `gorm:"foreignKey:BlockedID" json:"blocked,omitempty"`
}
```

---

## ER図

```
users
  ├─< posts (user_id)
  │   ├─< post_media (post_id)
  │   ├─< comments (post_id)
  │   ├─< likes (likeable_id)
  │   ├─< bookmarks (post_id)
  │   ├─< retweets (post_id)
  │   └─<>─ hashtags (post_hashtags)
  ├─< comments (user_id)
  │   └─< likes (likeable_id)
  ├─< likes (user_id)
  ├─< bookmarks (user_id)
  ├─< follows (follower_id / followed_id)
  ├─< notifications (user_id / actor_id)
  ├─< reports (reporter_id)
  ├─< retweets (user_id)
  └─< blocks (blocker_id / blocked_id)

comments
  └─< comments (parent_comment_id) [self-reference]
```

## マイグレーション戦略

### 初期マイグレーション (Phase 1)
1. `users` テーブル作成
2. `posts` テーブル作成
3. `post_media` テーブル作成（Phase 2で使用）
4. `comments` テーブル作成
5. `likes` テーブル作成
6. `bookmarks` テーブル作成
7. `follows` テーブル作成
8. `notifications` テーブル作成
9. `reports` テーブル作成

### Phase 2 マイグレーション
1. `retweets` テーブル作成
2. `hashtags` テーブル作成
3. `post_hashtags` テーブル作成
4. `blocks` テーブル作成
5. `users` テーブルに `avatar_url`, `header_url` カラム追加
6. `posts` テーブルに `visibility` カラム追加

### GORM AutoMigrate 使用例

```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &Post{},
        &PostMedia{},
        &Comment{},
        &Like{},
        &Bookmark{},
        &Follow{},
        &Notification{},
        &Report{},
        // Phase 2
        &Retweet{},
        &Hashtag{},
        &Block{},
    )
}
```

## パフォーマンス最適化

### 重要なインデックス
- `users.email`, `users.username`: ログイン・ユーザー検索
- `posts.user_id`, `posts.created_at`: タイムライン取得
- `comments.post_id`, `comments.parent_comment_id`: コメント取得
- `likes` の複合インデックス: いいね状態確認
- `follows.follower_id`, `follows.followed_id`: フォロー関係確認
- `notifications.user_id`, `notifications.is_read`: 未読通知確認

### N+1問題の回避
GORMの `Preload` を使用:
```go
db.Preload("User").Preload("Media").Find(&posts)
```

## 次のステップ
API仕様の詳細設計（03_API_SPECIFICATION.md）
