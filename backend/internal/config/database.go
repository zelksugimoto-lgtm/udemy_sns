package config

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB データベース接続を初期化
func InitDB() (*gorm.DB, error) {
	// 環境変数から接続情報を取得
	host := getEnv("DB_HOST", "db")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "social_media_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// DSN（Data Source Name）を構築
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// GORM設定
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// データベース接続
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	log.Println("データベース接続成功")

	return db, nil
}

// AutoMigrate データベースマイグレーションを実行
func AutoMigrate(db *gorm.DB) error {
	log.Println("データベースマイグレーション開始...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
		&model.Like{},
		&model.Bookmark{},
		&model.Follow{},
		&model.Notification{},
		&model.Report{},
	)

	if err != nil {
		return fmt.Errorf("マイグレーションエラー: %w", err)
	}

	log.Println("データベースマイグレーション完了")

	return nil
}

// getEnv 環境変数を取得（デフォルト値あり）
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
