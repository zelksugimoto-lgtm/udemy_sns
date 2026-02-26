# 通知システム仕様書

## 概要

Twitter/X風のSNSアプリケーションにおける通知システムの詳細仕様。
リアルタイム性（3分間隔ポーリング）、通知の集約、UI/UXを含む。

---

## 1. 通知の発生イベント

### 1.1 通知を発生させるアクション

| アクション | 通知タイプ | 通知対象ユーザー | 通知内容 |
|-----------|----------|---------------|---------|
| 投稿にいいね | `like` | 投稿者 | 「{ユーザー名}があなたの投稿にいいねしました」 |
| コメントにいいね | `like` | コメント投稿者 | 「{ユーザー名}があなたのコメントにいいねしました」 |
| 投稿にコメント | `comment` | 投稿者 | 「{ユーザー名}があなたの投稿にコメントしました」 |
| コメントに返信 | `reply` | 親コメント投稿者 | 「{ユーザー名}があなたのコメントに返信しました」 |
| フォロー | `follow` | フォローされたユーザー | 「{ユーザー名}があなたをフォローしました」 |

### 1.2 通知を発生させない条件（抑制ルール）

1. **自分自身のアクション**: 自分の投稿にいいね、自分へのフォローなどは通知しない
2. **ブロック中**: ブロックしているユーザーからのアクションは通知しない
3. **ミュート中**: ミュートしているユーザーからのアクションは通知しない（Phase 2で実装）
4. **削除済みコンテンツ**: 投稿やコメントが削除された場合、その通知も削除
5. **アクション取り消し**: いいね解除、フォロー解除時は通知を削除

---

## 2. 通知の集約とグルーピング

### 2.1 集約の基本方針

**「通知を発生させる単位（イベント）」と「表示を整理する単位（グルーピング）」を分離**

- **イベント単位**: 各アクションごとに通知レコードを作成（DB保存）
- **表示単位**: 同じ条件の通知を集約して表示（フロントエンド処理）

### 2.2 グルーピングの軸

**集約条件**: `タイプ × 対象（投稿/コメント） × 時間窓`

| 軸 | 詳細 |
|----|------|
| **タイプ** | `like`, `comment`, `reply`, `follow` |
| **対象ID** | 投稿ID、コメントID（対象がある場合のみ） |
| **時間窓** | 24時間以内の同一条件の通知を集約 |

### 2.3 集約表示の例

#### パターン1: 単一ユーザー
```
{アバター} {ユーザー名}があなたの投稿にいいねしました
```

#### パターン2: 複数ユーザー（2〜3人）
```
{アバター3つ} {ユーザーA}、{ユーザーB}、{ユーザーC}があなたの投稿にいいねしました
```

#### パターン3: 多数ユーザー（4人以上）
```
{アバター3つ+数字} {ユーザーA}、{ユーザーB}、他2人があなたの投稿にいいねしました
```

### 2.4 集約ロジック（フロントエンド）

```typescript
// 擬似コード
function groupNotifications(notifications: Notification[]): GroupedNotification[] {
  const groups = new Map<string, Notification[]>();

  for (const notification of notifications) {
    // グループキー: type_targetType_targetId_timeWindow
    const key = `${notification.type}_${notification.target_type}_${notification.target_id}_${getTimeWindow(notification.created_at)}`;

    if (!groups.has(key)) {
      groups.set(key, []);
    }
    groups.get(key)!.push(notification);
  }

  return Array.from(groups.values()).map(group => ({
    type: group[0].type,
    actors: group.map(n => n.actor),
    target: group[0].target,
    count: group.length,
    latestCreatedAt: group[0].created_at,
    isRead: group.every(n => n.is_read),
  }));
}

function getTimeWindow(createdAt: string): string {
  const now = new Date();
  const date = new Date(createdAt);
  const diffHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60);

  if (diffHours <= 24) return 'today';
  if (diffHours <= 48) return 'yesterday';
  return 'older';
}
```

---

## 3. UI/UX仕様

### 3.1 日付セパレータ

通知一覧を時系列で表示し、以下のセパレータで区切る：

```
━━━━━━━━━━━━━━━━
  今日
━━━━━━━━━━━━━━━━
[通知1]
[通知2]
━━━━━━━━━━━━━━━━
  昨日
━━━━━━━━━━━━━━━━
[通知3]
[通知4]
━━━━━━━━━━━━━━━━
  それ以前
━━━━━━━━━━━━━━━━
[通知5]
[通知6]
```

### 3.2 未読/既読の表示

| 状態 | 表示方法 |
|------|---------|
| **未読** | 背景色: `action.hover`、フォント: 太字、青い点インジケーター |
| **既読** | 背景色: 透明、フォント: 通常 |

### 3.3 通知アイテムのレイアウト

```
┌─────────────────────────────────────┐
│ [アバター] {テキスト}                  │
│   [アイコン]  {時刻}                   │
│   ⚫ (未読インジケーター)               │
└─────────────────────────────────────┘
```

### 3.4 ヘッダーの未読バッジ

```
🔔 5   ← 未読通知が5件
```

- 未読数が0の場合: バッジ非表示
- 未読数が1〜99の場合: 数字表示
- 未読数が100以上の場合: "99+" 表示

---

## 4. リアルタイム性（ポーリング）

### 4.1 ポーリング仕様

| 項目 | 仕様 |
|------|------|
| **更新間隔** | 3分（180秒） |
| **対象API** | `GET /api/v1/notifications?limit=50&offset=0` |
| **条件** | ログイン中のユーザーのみ |
| **停止条件** | ログアウト、ページを離れる |

### 4.2 新着通知の挙動

#### ケース1: 通知ページ閲覧中
1. ポーリングで新着通知を検知
2. リストの先頭に新しい通知を追加（アニメーション付き）
3. 未読バッジの数を更新
4. **既読化はしない**（ユーザーがクリックするまで未読維持）

#### ケース2: 他のページ閲覧中
1. ポーリングで新着通知を検知
2. ヘッダーの未読バッジを更新
3. 通知アイコンを軽くアニメーション（注意を引く）

### 4.3 既読化のタイミング

| アクション | 既読化 |
|-----------|--------|
| 通知ページを開く | **しない** |
| 通知アイテムをクリック | その通知のみ既読化 |
| 「すべて既読にする」ボタン | 全て既読化 |

**変更理由**: 現在の実装では通知ページを開くと自動で全既読になるが、これはUX的に好ましくない。ユーザーが明示的にクリックした通知のみ既読化する。

---

## 5. データベース設計

### 5.1 テーブル構造（現行）

```sql
CREATE TABLE notifications (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,           -- 通知を受け取るユーザー
  actor_id UUID,                    -- アクションを起こしたユーザー
  type VARCHAR(50) NOT NULL,        -- like, comment, reply, follow
  target_type VARCHAR(50),          -- Post, Comment
  target_id UUID,                   -- 投稿ID、コメントID
  message TEXT NOT NULL,            -- 通知メッセージ
  is_read BOOLEAN DEFAULT FALSE,    -- 既読フラグ
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,
  INDEX idx_user_created (user_id, created_at DESC),
  INDEX idx_user_read (user_id, is_read),
  INDEX idx_target (target_type, target_id)
);
```

### 5.2 通知削除ポリシー

1. **アクション取り消し時**: 該当する通知を論理削除
2. **対象コンテンツ削除時**: 関連する通知を論理削除
3. **ユーザー削除時**: 該当ユーザーが関与する通知を論理削除
4. **保存期間**: 90日以上経過した既読通知は物理削除（バッチ処理）

---

## 6. API仕様

### 6.1 通知一覧取得

```
GET /api/v1/notifications
```

**クエリパラメータ**:
- `limit`: 取得件数（デフォルト: 50）
- `offset`: オフセット（デフォルト: 0）

**レスポンス**:
```json
{
  "notifications": [
    {
      "id": "uuid",
      "type": "like",
      "message": "ユーザーAがあなたの投稿にいいねしました",
      "is_read": false,
      "actor_user": {
        "id": "uuid",
        "username": "userA",
        "display_name": "ユーザーA",
        "avatar_url": "https://..."
      },
      "target_type": "Post",
      "target_id": "uuid",
      "created_at": "2026-02-26T12:00:00Z"
    }
  ],
  "unread_count": 5,
  "pagination": {
    "total": 100,
    "limit": 50,
    "offset": 0,
    "has_more": true
  }
}
```

### 6.2 未読通知数取得（新規）

```
GET /api/v1/notifications/unread-count
```

**レスポンス**:
```json
{
  "count": 5
}
```

### 6.3 個別既読化

```
PATCH /api/v1/notifications/:id/read
```

### 6.4 全既読化

```
POST /api/v1/notifications/read-all
```

---

## 7. 実装フェーズ

### Phase 1（今回実装）

- ✅ 通知の自動生成（いいね、コメント、フォロー）
- ✅ 通知一覧の集約表示
- ✅ 日付セパレータ
- ✅ 未読バッジ表示
- ✅ 3分間隔ポーリング
- ✅ クリック時の既読化（ページ表示時の自動既読を廃止）
- ✅ 自分自身のアクション除外
- ✅ 未読通知数取得API

### Phase 2（将来実装）

- ⬜ ミュート機能
- ⬜ 通知設定（通知タイプごとのON/OFF）
- ⬜ WebSocket/SSEによるリアルタイム通知
- ⬜ プッシュ通知（PWA）
- ⬜ メール通知
- ⬜ 通知の詳細フィルタリング

---

## 8. パフォーマンス考慮事項

### 8.1 データベースインデックス

```sql
-- 通知一覧取得の高速化
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);

-- 未読数カウントの高速化
CREATE INDEX idx_notifications_user_read ON notifications(user_id, is_read);

-- 通知削除時の検索高速化
CREATE INDEX idx_notifications_target ON notifications(target_type, target_id);
CREATE INDEX idx_notifications_actor ON notifications(actor_id, created_at DESC);
```

### 8.2 キャッシング戦略

- 未読通知数: Redis（TTL: 30秒）
- 通知一覧: クライアント側React Queryキャッシュ（staleTime: 180秒 = 3分）

---

## 9. テストケース

### 9.1 通知生成テスト

| テストケース | 期待結果 |
|------------|---------|
| ユーザーAがユーザーBの投稿にいいね | ユーザーBに通知が届く |
| ユーザーAが自分の投稿にいいね | 通知が発生しない |
| ユーザーAがユーザーBをフォロー | ユーザーBに通知が届く |
| ユーザーAがいいねを取り消し | 通知が削除される |
| 投稿が削除される | 関連する全通知が削除される |

### 9.2 集約表示テスト

| テストケース | 期待結果 |
|------------|---------|
| 同じ投稿に3人がいいね（24時間以内） | 1つのグループとして表示 |
| 異なる投稿にいいね | 別々のグループとして表示 |
| 25時間前のいいねと新しいいいね | 別々のグループとして表示 |

### 9.3 ポーリングテスト

| テストケース | 期待結果 |
|------------|---------|
| 3分経過後 | 新しい通知が取得される |
| ログアウト | ポーリング停止 |
| ページ遷移 | ポーリング継続（グローバル） |

---

## 10. セキュリティ考慮事項

1. **認証**: すべての通知APIは認証必須
2. **認可**: 自分宛ての通知のみ取得可能
3. **レート制限**: 未読数取得APIは1分間に60回まで
4. **SQL注入対策**: GORM使用により自動対策
5. **XSS対策**: 通知メッセージのサニタイズ

---

## 11. エラーハンドリング

| エラーケース | HTTPステータス | 対応 |
|-------------|--------------|------|
| 未認証 | 401 | ログインページへリダイレクト |
| 他人の通知を既読化 | 403 | エラーメッセージ表示 |
| 通知が見つからない | 404 | エラーメッセージ表示 |
| サーバーエラー | 500 | リトライ、エラー通知 |

---

## 付録: 通知メッセージテンプレート

```typescript
const NOTIFICATION_TEMPLATES = {
  like: {
    post: (actor: string) => `${actor}があなたの投稿にいいねしました`,
    comment: (actor: string) => `${actor}があなたのコメントにいいねしました`,
  },
  comment: (actor: string) => `${actor}があなたの投稿にコメントしました`,
  reply: (actor: string) => `${actor}があなたのコメントに返信しました`,
  follow: (actor: string) => `${actor}があなたをフォローしました`,
};
```
