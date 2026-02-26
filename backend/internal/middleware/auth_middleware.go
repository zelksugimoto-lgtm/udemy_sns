package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/yourusername/sns-app/pkg/errors"
)

// AuthMiddleware JWT認証ミドルウェア
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("認証が必要です"))
			}

			// "Bearer "プレフィックスを確認
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("無効な認証形式です"))
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
				return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("無効なトークンです"))
			}

			// クレームを取得
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// ユーザーIDを取得
				userIDStr, ok := claims["user_id"].(string)
				if !ok {
					return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("無効なトークンです"))
				}

				// UUIDに変換
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("無効なユーザーIDです"))
				}

				// コンテキストにユーザーIDを設定
				c.Set("user_id", userID)

				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized("無効なトークンです"))
		}
	}
}

// OptionalAuthMiddleware オプションのJWT認証ミドルウェア（認証失敗でもエラーを返さない）
func OptionalAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// 認証ヘッダーがない場合はそのまま次へ
				return next(c)
			}

			// "Bearer "プレフィックスを確認
			if !strings.HasPrefix(authHeader, "Bearer ") {
				// 無効な形式でもそのまま次へ
				return next(c)
			}

			// トークンを抽出
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// JWTシークレットを取得
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-secret-key"
			}

			// トークンをパース
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(secret), nil
			})

			if err != nil {
				// パースエラーでもそのまま次へ
				return next(c)
			}

			// クレームを取得
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userIDStr, ok := claims["user_id"].(string)
				if !ok {
					return next(c)
				}

				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					return next(c)
				}

				// コンテキストにユーザーIDを設定
				c.Set("user_id", userID)
			}

			return next(c)
		}
	}
}

// GetUserID コンテキストからユーザーIDを取得
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("認証が必要です")
	}
	return userID, nil
}
