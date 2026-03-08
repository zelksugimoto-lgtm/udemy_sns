package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	urlpkg "net/url"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/config"
)

// StorageService Firebase Storageサービス
type StorageService struct {
	storageBucket string
}

// NewStorageService StorageServiceのコンストラクタ
func NewStorageService(storageBucket string) *StorageService {
	return &StorageService{
		storageBucket: storageBucket,
	}
}

// ValidateMediaURL メディアURLがFirebase Storage URLかどうかを検証
// セキュリティ: ドメイン混同攻撃を防ぐため、完全一致でホストを検証
func (s *StorageService) ValidateMediaURL(url string) error {
	if url == "" {
		return fmt.Errorf("メディアURLが空です")
	}

	// URLをパース
	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return fmt.Errorf("無効なURLです")
	}

	// スキームを検証（HTTPSのみ許可）
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("HTTPSのURLのみ許可されています")
	}

	// ホストを完全一致で検証（前方一致ではない）
	// これによりサブドメイン混同攻撃を防止
	validHosts := []string{
		"firebasestorage.googleapis.com",
		"storage.googleapis.com",
	}

	hostValid := false
	for _, validHost := range validHosts {
		if parsedURL.Host == validHost {
			hostValid = true
			break
		}
	}

	if !hostValid {
		return fmt.Errorf("無効なホストです")
	}

	// パス内のバケット名を検証
	if parsedURL.Host == "firebasestorage.googleapis.com" {
		// 形式: /v0/b/{bucket}/o/...
		expectedPathPrefix := fmt.Sprintf("/v0/b/%s/o/", s.storageBucket)
		if !strings.HasPrefix(parsedURL.Path, expectedPathPrefix) {
			return fmt.Errorf("無効なFirebase Storage URLです")
		}
	}

	if parsedURL.Host == "storage.googleapis.com" {
		// 形式: /{bucket}/...
		expectedPathPrefix := fmt.Sprintf("/%s/", s.storageBucket)
		if !strings.HasPrefix(parsedURL.Path, expectedPathPrefix) {
			return fmt.Errorf("無効なFirebase Storage URLです")
		}
	}

	return nil
}

// ValidateMediaType メディアタイプを検証
func (s *StorageService) ValidateMediaType(mediaType string) error {
	validTypes := map[string]bool{
		"image": true,
		"video": true,
	}

	if !validTypes[mediaType] {
		return fmt.Errorf("無効なメディアタイプです: %s（有効な値: image, video）", mediaType)
	}

	return nil
}

// ValidateMediaOwnership メディアURLが指定ユーザーに属することを検証
// セキュリティ: IDOR攻撃を防ぐため、URLパス内のユーザーIDを検証
func (s *StorageService) ValidateMediaOwnership(url string, userID uuid.UUID) error {
	// URLをパース
	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return fmt.Errorf("無効なURLです")
	}

	// 期待されるパスプレフィックス: posts/{userID}/...
	expectedPathPrefix := fmt.Sprintf("posts/%s/", userID.String())

	// firebasestorage.googleapis.com の場合
	if parsedURL.Host == "firebasestorage.googleapis.com" {
		// パス形式: /v0/b/{bucket}/o/posts/{userID}/...
		// /o/ 以降のオブジェクトパスを抽出
		parts := strings.Split(parsedURL.Path, "/o/")
		if len(parts) != 2 {
			return fmt.Errorf("無効なURLパスです")
		}

		// URLエンコードされている可能性があるのでデコード
		objectPath, err := urlpkg.QueryUnescape(parts[1])
		if err != nil {
			return fmt.Errorf("URLパスのデコードに失敗しました")
		}

		if !strings.HasPrefix(objectPath, expectedPathPrefix) {
			return fmt.Errorf("このメディアを使用する権限がありません")
		}
	}

	// storage.googleapis.com の場合
	if parsedURL.Host == "storage.googleapis.com" {
		// パス形式: /{bucket}/posts/{userID}/...
		// バケット名以降のパスを抽出
		pathParts := strings.SplitN(parsedURL.Path, "/", 3) // ["", "{bucket}", "posts/{userID}/..."]
		if len(pathParts) < 3 {
			return fmt.Errorf("無効なURLパスです")
		}

		objectPath := pathParts[2]
		if !strings.HasPrefix(objectPath, expectedPathPrefix) {
			return fmt.Errorf("このメディアを使用する権限がありません")
		}
	}

	return nil
}

// UploadFile Firebase Storageにファイルをアップロード
func (s *StorageService) UploadFile(ctx context.Context, path string, file io.Reader, contentType string) (string, error) {
	// Firebase Storage clientを初期化
	storageClient, err := config.InitFirebase()
	if err != nil {
		return "", fmt.Errorf("Firebase Storageの初期化に失敗しました: %w", err)
	}

	// バケットを取得
	bucket, err := storageClient.Bucket(s.storageBucket)
	if err != nil {
		return "", fmt.Errorf("バケットの取得に失敗しました: %w", err)
	}

	// オブジェクトを作成
	obj := bucket.Object(path)

	// ライターを作成
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	// ファイルをコピー
	if _, err := io.Copy(wc, file); err != nil {
		wc.Close()
		return "", fmt.Errorf("ファイルのアップロードに失敗しました: %w", err)
	}

	// ライターを閉じる
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("ファイルのアップロードに失敗しました: %w", err)
	}

	// 署名付きURL（Signed URL）を生成
	// サービスアカウントキーから認証情報を取得
	serviceAccountKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	if serviceAccountKey == "" {
		return "", fmt.Errorf("FIREBASE_SERVICE_ACCOUNT_KEY環境変数が設定されていません")
	}

	var creds map[string]interface{}
	if err := json.Unmarshal([]byte(serviceAccountKey), &creds); err != nil {
		return "", fmt.Errorf("サービスアカウントキーの解析に失敗しました: %w", err)
	}

	// 秘密鍵を取得
	privateKeyPEM := creds["private_key"].(string)
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("秘密鍵のデコードに失敗しました")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("秘密鍵のパースに失敗しました: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("RSA秘密鍵の変換に失敗しました")
	}

	// 秘密鍵をPEM形式のバイト配列に変換
	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKey),
	})

	// 署名付きURLを生成（有効期限: 10年）
	opts := &storage.SignedURLOptions{
		GoogleAccessID: creds["client_email"].(string),
		PrivateKey:     privateKeyBytes,
		Expires:        time.Now().Add(10 * 365 * 24 * time.Hour), // 10年
		Method:         "GET",
	}

	signedURL, err := storage.SignedURL(s.storageBucket, path, opts)
	if err != nil {
		return "", fmt.Errorf("署名付きURLの生成に失敗しました: %w", err)
	}

	return signedURL, nil
}
