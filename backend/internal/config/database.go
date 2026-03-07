package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/yourusername/sns-app/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// クエリカウンター（リクエスト単位でリセット可能）
var queryCounter int32

// カスタムロガー（クエリ番号を表示）
type CustomLogger struct {
	logger.Interface
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// クエリカウンターをインクリメント
	count := atomic.AddInt32(&queryCounter, 1)

	sql, rows := fc()
	elapsed := time.Since(begin)

	// クエリ番号と実行時間を表示
	log.Printf("[Query #%d] [%.2fms] Rows: %d | SQL: %s", count, float64(elapsed.Microseconds())/1000, rows, sql)

	// 元のロガーにも渡す
	l.Interface.Trace(ctx, begin, fc, err)
}

// ResetQueryCounter クエリカウンターをリセット
func ResetQueryCounter() {
	atomic.StoreInt32(&queryCounter, 0)
}

// GetQueryCount 現在のクエリ数を取得
func GetQueryCount() int32 {
	return atomic.LoadInt32(&queryCounter)
}

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

	// GORM設定（カスタムロガーを使用）
	config := &gorm.Config{
		Logger: &CustomLogger{
			Interface: logger.Default.LogMode(logger.Info),
		},
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
