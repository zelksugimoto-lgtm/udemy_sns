package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/pkg/errors"
)

// AuthMiddleware JWT認証ミドルウェア
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
			}

			// "Bearer "プレフィックスを確認
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, errors.Unauthorized("無効な認証形式です"))
			}

			// トークンを抽出
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// JWTシークレットを取得
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-secret-key" // 開発環境用のデフォルト値
			}

			// トークンをパース
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 署名方式の検証
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(secret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, errors.Unauthorized("無効なトークンです"))
			}

			// クレームを取得
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// ユーザーIDを取得
				userIDStr, ok := claims["user_id"].(string)
				if !ok {
					return c.JSON(http.StatusUnauthorized, errors.Unauthorized("無効なトークンです"))
				}

				// UUIDに変換
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, errors.Unauthorized("無効なユーザーIDです"))
				}

				// コンテキストにユーザーIDを設定
				c.Set("user_id", userID)

				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, errors.Unauthorized("無効なトークンです"))
		}
	}
}

// GetUserID コンテキストからユーザーIDを取得
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, http.ErrAbortHandler
	}
	return userID, nil
}
