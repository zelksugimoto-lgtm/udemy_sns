package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init ロガーを初期化
func Init(env string) {
	// ログレベルを環境に応じて設定
	var logLevel zerolog.Level
	switch env {
	case "production":
		logLevel = zerolog.InfoLevel
	case "test":
		logLevel = zerolog.WarnLevel // テストではWARN以上のみ
	default: // development
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	// 開発環境ではコンソール出力を見やすく、本番環境ではJSON形式
	var output io.Writer
	if env == "development" {
		// 開発環境: 色付きのコンソール出力
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		// 本番/テスト環境: JSON出力
		output = os.Stdout
	}

	// グローバルロガーを設定
	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller(). // ファイル名と行番号を含める
		Logger()
}

// GetLogger コンテキストに対応したロガーを取得（将来的にコンテキストから取得できるように）
func GetLogger() *zerolog.Logger {
	return &log.Logger
}
