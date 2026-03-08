package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

// InitFirebase Firebase Admin SDKを初期化
// 環境変数 FIREBASE_SERVICE_ACCOUNT_KEY からサービスアカウントキーのJSON文字列を読み込む
func InitFirebase() (*storage.Client, error) {
	ctx := context.Background()

	// 環境変数からサービスアカウントキーを取得
	serviceAccountKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	if serviceAccountKey == "" {
		// 環境変数が設定されていない場合、Firebase機能は利用不可として警告を出すが、エラーにはしない
		// これにより、Firebase機能を使わない場合でもアプリが起動できる
		return nil, fmt.Errorf("FIREBASE_SERVICE_ACCOUNT_KEY環境変数が設定されていません。画像アップロード機能は利用できません")
	}

	// JSON文字列をバイト配列に変換して検証
	var credJSON map[string]interface{}
	if err := json.Unmarshal([]byte(serviceAccountKey), &credJSON); err != nil {
		return nil, fmt.Errorf("FIREBASE_SERVICE_ACCOUNT_KEYの形式が不正です: %w", err)
	}

	// Firebase設定
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	storageBucket := os.Getenv("FIREBASE_STORAGE_BUCKET")

	if projectID == "" || storageBucket == "" {
		return nil, fmt.Errorf("FIREBASE_PROJECT_ID または FIREBASE_STORAGE_BUCKET が設定されていません")
	}

	config := &firebase.Config{
		ProjectID:     projectID,
		StorageBucket: storageBucket,
	}

	// Firebase Admin SDK初期化
	app, err := firebase.NewApp(ctx, config, option.WithCredentialsJSON([]byte(serviceAccountKey)))
	if err != nil {
		return nil, fmt.Errorf("Firebase初期化エラー: %w", err)
	}

	// Storage Client取得
	client, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("Firebase Storage Clientの取得エラー: %w", err)
	}

	return client, nil
}
