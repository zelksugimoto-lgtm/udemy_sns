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

const (
	// AccessTokenCookie アクセストークンのCookie名
	AccessTokenCookie = "access_token"
)

// AuthMiddleware JWT認証ミドルウェア（Cookie優先、AuthorizationヘッダーもサポートAuthMiddleware）
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// トークンを取得（Cookie優先）
			tokenString, err := getTokenFromRequest(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized(err.Error()))
			}

			// トークンを検証してユーザーIDを取得
			userID, err := validateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, pkgErrors.Unauthorized(err.Error()))
			}

			// コンテキストにユーザーIDを設定
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware オプションのJWT認証ミドルウェア（認証失敗でもエラーを返さない）
func OptionalAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// トークンを取得（Cookie優先）
			tokenString, err := getTokenFromRequest(c)
			if err != nil {
				// 認証情報がない場合はそのまま次へ
				return next(c)
			}

			// トークンを検証してユーザーIDを取得
			userID, err := validateToken(tokenString)
			if err != nil {
				// 検証失敗でもそのまま次へ
				return next(c)
			}

			// コンテキストにユーザーIDを設定
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// getTokenFromRequest リクエストからトークンを取得（Cookie優先、Authorizationヘッダーもサポート）
func getTokenFromRequest(c echo.Context) (string, error) {
	// 1. Cookieからトークンを取得（優先）
	cookie, err := c.Cookie(AccessTokenCookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	// 2. Authorization ヘッダーからトークンを取得（フォールバック）
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("認証が必要です")
	}

	// "Bearer "プレフィックスを確認
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("無効な認証形式です")
	}

	// トークンを抽出
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return "", errors.New("トークンが空です")
	}

	return tokenString, nil
}

// validateToken トークンを検証してユーザーIDを返す
func validateToken(tokenString string) (uuid.UUID, error) {
	// JWTシークレットを取得
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // 開発環境用のデフォルト値
	}

	// トークンをパース
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名方式の検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無効な署名アルゴリズムです")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return uuid.Nil, errors.New("無効なトークンです")
	}

	if !token.Valid {
		return uuid.Nil, errors.New("トークンの有効期限が切れています")
	}

	// クレームを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("無効なトークンです")
	}

	// ユーザーIDを取得
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("トークンにuser_idが含まれていません")
	}

	// UUIDに変換
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("無効なユーザーIDです")
	}

	return userID, nil
}

// GetUserID コンテキストからユーザーIDを取得
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("認証が必要です")
	}
	return userID, nil
}
