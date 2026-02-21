# ファイルアップロード仕様書 (Phase 2)

## 概要

Firebase Storageを使用した画像・動画のアップロード機能

## サポートするファイル形式

### 画像
- **形式**: JPEG, PNG, GIF, WebP
- **最大サイズ**: 5MB（環境変数で変更可能）
- **最大枚数**: 4枚/投稿

### 動画
- **形式**: MP4, MOV, AVI, WebM
- **最大サイズ**: 20MB（環境変数で変更可能、デフォルト）
- **最大枚数**: 4本/投稿（画像と合わせて最大4つ）

## アップロードフロー

### 投稿時のメディアアップロード

```
[クライアント]
    ↓ 1. ファイル選択
    <input type="file" accept="image/*,video/*" multiple />
    ↓ 2. クライアント側バリデーション
    - ファイル形式チェック
    - ファイルサイズチェック
    - 枚数チェック（最大4つ）
    ↓ 3. Firebase Storageへ直接アップロード
    const storageRef = ref(storage, `posts/${userId}/${uuid()}.${ext}`);
    await uploadBytes(storageRef, file);
    const url = await getDownloadURL(storageRef);
    ↓ 4. 投稿作成APIリクエスト
    POST /api/v1/posts
    {
      content: "投稿内容",
      media: [
        { media_type: "image", media_url: "https://..." },
        { media_type: "video", media_url: "https://..." }
      ]
    }
    ↓
[バックエンド]
    ↓ 5. URL検証
    - Firebase Storage URLか確認
    - メディアタイプ検証
    ↓ 6. 投稿 & メディア保存
    - posts テーブルに投稿保存
    - post_media テーブルにメディア情報保存
    ↓ 7. レスポンス返却
```

### プロフィール画像アップロード

```
[クライアント]
    ↓ 1. アイコン画像選択
    ↓ 2. Firebase Storageへアップロード
    const storageRef = ref(storage, `avatars/${userId}/${uuid()}.${ext}`);
    await uploadBytes(storageRef, file);
    const avatarUrl = await getDownloadURL(storageRef);
    ↓ 3. プロフィール更新APIリクエスト
    PATCH /api/v1/users/me
    {
      avatar_url: "https://..."
    }
    ↓
[バックエンド]
    ↓ 4. URL検証 & ユーザー更新
```

---

## Firebase Storage ディレクトリ構造

```
your-bucket/
├── posts/
│   ├── {user_id}/
│   │   ├── {uuid}.jpg
│   │   ├── {uuid}.png
│   │   └── {uuid}.mp4
├── avatars/
│   └── {user_id}/
│       └── {uuid}.jpg
└── headers/
    └── {user_id}/
        └── {uuid}.jpg
```

---

## バックエンド実装

### 環境変数

```bash
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_STORAGE_BUCKET=your-bucket.appspot.com
MAX_VIDEO_SIZE_MB=20
MAX_IMAGE_SIZE_MB=5
MAX_UPLOAD_FILES=4
```

### Firebase Admin SDK 初期化 (Go)

```go
package config

import (
    "context"
    firebase "firebase.google.com/go"
    "firebase.google.com/go/storage"
    "google.golang.org/api/option"
)

func InitFirebase() (*storage.Client, error) {
    ctx := context.Background()

    config := &firebase.Config{
        ProjectID:     os.Getenv("FIREBASE_PROJECT_ID"),
        StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET"),
    }

    app, err := firebase.NewApp(ctx, config, option.WithCredentialsFile("serviceAccountKey.json"))
    if err != nil {
        return nil, err
    }

    client, err := app.Storage(ctx)
    if err != nil {
        return nil, err
    }

    return client, nil
}
```

### URL検証

```go
package service

import (
    "fmt"
    "strings"
)

func ValidateFirebaseStorageURL(url, bucket string) error {
    expectedPrefix := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s", bucket)

    if !strings.HasPrefix(url, expectedPrefix) {
        return errors.New("invalid firebase storage URL")
    }

    return nil
}

func ValidateMediaType(mediaType string) error {
    validTypes := []string{"image", "video"}

    for _, t := range validTypes {
        if mediaType == t {
            return nil
        }
    }

    return errors.New("invalid media type")
}
```

---

## フロントエンド実装

### Firebase初期化

```typescript
// config/firebase.ts
import { initializeApp } from 'firebase/app';
import { getStorage } from 'firebase/storage';

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
};

const app = initializeApp(firebaseConfig);
export const storage = getStorage(app);
```

### ファイルアップロードユーティリティ

```typescript
// utils/fileUpload.ts
import { ref, uploadBytes, getDownloadURL } from 'firebase/storage';
import { storage } from '@/config/firebase';
import { v4 as uuidv4 } from 'uuid';

const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB
const MAX_VIDEO_SIZE = 20 * 1024 * 1024; // 20MB

const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
const ALLOWED_VIDEO_TYPES = ['video/mp4', 'video/quicktime', 'video/x-msvideo', 'video/webm'];

export type MediaType = 'image' | 'video';

export const validateFile = (file: File): { valid: boolean; error?: string } => {
  const isImage = ALLOWED_IMAGE_TYPES.includes(file.type);
  const isVideo = ALLOWED_VIDEO_TYPES.includes(file.type);

  if (!isImage && !isVideo) {
    return { valid: false, error: 'サポートされていないファイル形式です' };
  }

  if (isImage && file.size > MAX_IMAGE_SIZE) {
    return { valid: false, error: '画像サイズは5MB以下にしてください' };
  }

  if (isVideo && file.size > MAX_VIDEO_SIZE) {
    return { valid: false, error: '動画サイズは20MB以下にしてください' };
  }

  return { valid: true };
};

export const getMediaType = (file: File): MediaType => {
  return ALLOWED_IMAGE_TYPES.includes(file.type) ? 'image' : 'video';
};

export const uploadFile = async (
  file: File,
  userId: string,
  directory: 'posts' | 'avatars' | 'headers'
): Promise<{ url: string; mediaType: MediaType }> => {
  const validation = validateFile(file);
  if (!validation.valid) {
    throw new Error(validation.error);
  }

  const ext = file.name.split('.').pop();
  const filename = `${uuidv4()}.${ext}`;
  const storageRef = ref(storage, `${directory}/${userId}/${filename}`);

  await uploadBytes(storageRef, file);
  const url = await getDownloadURL(storageRef);

  return {
    url,
    mediaType: getMediaType(file),
  };
};
```

### メディアアップロードコンポーネント

```typescript
// components/post/MediaUpload.tsx
import { useState } from 'react';
import { uploadFile } from '@/utils/fileUpload';
import { useAuth } from '@/contexts/AuthContext';

type Media = {
  media_type: 'image' | 'video';
  media_url: string;
  preview?: string;
};

export const MediaUpload = ({ onMediaChange }) => {
  const { user } = useAuth();
  const [media, setMedia] = useState<Media[]>([]);
  const [uploading, setUploading] = useState(false);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);

    if (media.length + files.length > 4) {
      alert('メディアは最大4つまでアップロードできます');
      return;
    }

    setUploading(true);

    try {
      const uploadPromises = files.map(file =>
        uploadFile(file, user!.id, 'posts')
      );

      const results = await Promise.all(uploadPromises);

      const newMedia: Media[] = results.map(result => ({
        media_type: result.mediaType,
        media_url: result.url,
      }));

      const updatedMedia = [...media, ...newMedia];
      setMedia(updatedMedia);
      onMediaChange(updatedMedia);
    } catch (error) {
      alert(error.message);
    } finally {
      setUploading(false);
    }
  };

  const handleRemove = (index: number) => {
    const updatedMedia = media.filter((_, i) => i !== index);
    setMedia(updatedMedia);
    onMediaChange(updatedMedia);
  };

  return (
    <div>
      <input
        type="file"
        accept="image/*,video/*"
        multiple
        onChange={handleFileSelect}
        disabled={uploading || media.length >= 4}
      />

      {uploading && <CircularProgress />}

      <div>
        {media.map((item, index) => (
          <div key={index}>
            {item.media_type === 'image' ? (
              <img src={item.media_url} alt="" />
            ) : (
              <video src={item.media_url} controls />
            )}
            <IconButton onClick={() => handleRemove(index)}>
              <DeleteIcon />
            </IconButton>
          </div>
        ))}
      </div>
    </div>
  );
};
```

---

## セキュリティルール (Firebase Storage)

```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    // 認証済みユーザーのみアップロード可能
    match /posts/{userId}/{filename} {
      allow read: if true;  // 誰でも読み取り可能
      allow write: if request.auth != null
                   && request.auth.uid == userId
                   && request.resource.size < 20 * 1024 * 1024;  // 20MB
    }

    match /avatars/{userId}/{filename} {
      allow read: if true;
      allow write: if request.auth != null
                   && request.auth.uid == userId
                   && request.resource.size < 5 * 1024 * 1024;  // 5MB
    }

    match /headers/{userId}/{filename} {
      allow read: if true;
      allow write: if request.auth != null
                   && request.auth.uid == userId
                   && request.resource.size < 5 * 1024 * 1024;
    }
  }
}
```

---

## パフォーマンス最適化

### クライアント側
1. **プログレッシブアップロード**: 複数ファイルを並列アップロード
2. **プレビュー生成**: アップロード前にローカルでプレビュー表示
3. **圧縮**: 画像を圧縮してからアップロード（browser-image-compression等）

### サーバー側
1. **CDN**: Firebase Storage は自動的にCDN配信
2. **キャッシュ**: `Cache-Control` ヘッダー設定

---

## エラーハンドリング

### クライアント側エラー
- ファイル形式が不正
- ファイルサイズ超過
- 最大枚数超過
- ネットワークエラー

### サーバー側エラー
- 無効なFirebase Storage URL
- 不正なメディアタイプ

---

## 将来の拡張（Phase 4）

### サムネイル自動生成
- Cloud Functionsを使用
- 動画: 最初のフレームをサムネイル化
- 画像: リサイズしてサムネイル生成

```javascript
// Cloud Functions (例)
const functions = require('firebase-functions');
const admin = require('firebase-admin');
const spawn = require('child-process-promise').spawn;

exports.generateThumbnail = functions.storage.object().onFinalize(async (object) => {
  const filePath = object.name;

  if (filePath.startsWith('posts/') && object.contentType.startsWith('video/')) {
    // ffmpegで動画からサムネイル生成
    // ...
  }
});
```
