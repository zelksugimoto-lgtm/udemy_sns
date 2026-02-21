FROM golang:1.23-alpine

# 必要なパッケージのインストール
RUN apk add --no-cache git

WORKDIR /app

# Airのインストール（ホットリロード用）- Go 1.23互換バージョン
RUN go install github.com/air-verse/air@v1.61.5

# swaggoのインストール（OpenAPI生成用）
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

# 依存関係は起動時にインストール（go.modがない場合があるため）
# CMD内で go mod download を実行

# Airで起動
CMD ["air", "-c", ".air.toml"]
