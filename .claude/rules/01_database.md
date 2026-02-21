# データベース操作ルール

## 基本原則

### 1. 論理削除（Soft Delete）の使用
- **必須**: すべてのテーブルで論理削除を使用
- `deleted_at` フィールドで削除状態を管理
- 物理削除は原則禁止（管理者の明示的な指示がある場合のみ）

### 2. 共通フィールド（BaseModel）
すべてのテーブルに以下のフィールドを含める：
```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CreatedAt time.Time      `gorm:"not null"`
    UpdatedAt time.Time      `gorm:"not null"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 3. 主キー
- **UUID** を主キーとして使用
- PostgreSQLの `gen_random_uuid()` で自動生成

---

## マイグレーション

### GORM AutoMigrate 使用
```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &Post{},
        &Comment{},
        // ...
    )
}
```

### マイグレーション実行タイミング
- アプリケーション起動時に自動実行
- 新しいモデル追加時
- 既存モデルのフィールド追加時

### 注意事項
- **本番環境**: マイグレーション実行前に必ずバックアップ
- **カラム削除**: AutoMigrateでは削除されないため、手動SQL実行が必要

---

## インデックス

### インデックスを設定すべきカラム
1. **外部キー**: すべての外部キーにインデックス
2. **検索条件**: WHERE句で使用されるカラム
3. **ソート条件**: ORDER BY句で使用されるカラム
4. **ユニーク制約**: メールアドレス、ユーザー名等

### GORM インデックス設定例
```go
type User struct {
    Email    string `gorm:"uniqueIndex:idx_users_email,where:deleted_at IS NULL"`
    Username string `gorm:"uniqueIndex:idx_users_username,where:deleted_at IS NULL"`
}
```

### 論理削除を考慮したインデックス
- ユニークインデックスには `where:deleted_at IS NULL` を追加
- 削除済みレコードを除外してユニーク性を保証

---

## クエリのベストプラクティス

### 1. N+1問題の回避
```go
// 悪い例
posts := []Post{}
db.Find(&posts)
for _, post := range posts {
    db.Model(&post).Association("User").Find(&post.User)  // N+1発生
}

// 良い例
posts := []Post{}
db.Preload("User").Find(&posts)  // 一度にロード
```

### 2. ページネーション
```go
func GetPosts(db *gorm.DB, limit, offset int) ([]Post, error) {
    var posts []Post
    err := db.
        Preload("User").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&posts).
        Error
    return posts, err
}
```

### 3. 論理削除されたレコードの除外
- GORMは自動的に `deleted_at IS NULL` を追加
- 削除済みを含める場合: `db.Unscoped().Find(&posts)`

---

## トランザクション

### 使用すべき場面
- 複数テーブルへの書き込み
- 整合性が重要な操作
- エラー時のロールバックが必要な場合

### 実装例
```go
func CreatePostWithMedia(db *gorm.DB, post *Post, media []PostMedia) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 投稿作成
        if err := tx.Create(post).Error; err != nil {
            return err
        }

        // メディア作成
        for i := range media {
            media[i].PostID = post.ID
            if err := tx.Create(&media[i]).Error; err != nil {
                return err
            }
        }

        return nil
    })
}
```

---

## データ整合性

### 外部キー制約
```go
type Post struct {
    UserID uuid.UUID `gorm:"type:uuid;not null"`
    User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
```

### CHECK制約
```sql
-- 投稿の文字数制限
ALTER TABLE posts ADD CONSTRAINT content_length CHECK (LENGTH(content) <= 250);

-- 自分自身をフォローできない
ALTER TABLE follows ADD CONSTRAINT no_self_follow CHECK (follower_id != followed_id);
```

---

## セキュリティ

### SQLインジェクション対策
- **必須**: GORMのメソッドを使用（Rawクエリは避ける）
- パラメータは必ずプレースホルダーを使用

```go
// 悪い例
db.Raw("SELECT * FROM users WHERE email = '" + email + "'")  // SQLインジェクションの危険

// 良い例
db.Where("email = ?", email).First(&user)
```

### パスワードハッシュ化
```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

## パフォーマンス最適化

### 1. 必要なカラムのみ選択
```go
// 悪い例
db.Find(&users)  // すべてのカラムを取得

// 良い例
db.Select("id", "username", "display_name").Find(&users)
```

### 2. カウントクエリの最適化
```go
// 悪い例
var posts []Post
db.Find(&posts)
count := len(posts)  // 全レコードを取得

// 良い例
var count int64
db.Model(&Post{}).Count(&count)
```

### 3. バッチ処理
```go
// 大量のレコード挿入
db.CreateInBatches(posts, 100)  // 100件ずつバッチ処理
```

---

## エラーハンドリング

### GORMエラーの確認
```go
if err := db.First(&user, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New("user not found")
    }
    return nil, err
}
```

### エラーログ
- データベースエラーは必ずログに記録
- 本番環境では詳細なエラーをクライアントに返さない

---

## Docker環境での注意事項

### データ永続化
```yaml
# docker-compose.yml
services:
  db:
    volumes:
      - postgres_data:/var/lib/postgresql/data  # 永続化

volumes:
  postgres_data:
```

### 絶対に禁止
- **`docker compose down -v`**: ボリュームが削除され、データが消失
- **`docker volume rm postgres_data`**: データベースが完全削除

### 安全な操作
- 再起動: `docker compose restart db`
- 停止・起動: `docker compose down && docker compose up -d`

---

## トラブルシューティング

### マイグレーションエラー
1. データベース接続確認
2. モデル定義の確認（タグの誤り等）
3. 既存データとの整合性確認

### 接続エラー
1. 環境変数確認（`.env`）
2. Dockerコンテナ起動確認
3. ポート確認（5432が使用中でないか）

### パフォーマンス問題
1. スロークエリログ確認
2. インデックスの有無確認
3. N+1問題の確認

---

## 参照
- GORM公式ドキュメント: https://gorm.io/
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
