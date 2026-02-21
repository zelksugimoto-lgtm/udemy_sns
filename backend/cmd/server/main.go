package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yourusername/sns-app/internal/handler"

	_ "github.com/yourusername/sns-app/docs"
)

// @title           SNS Application API
// @version         1.0
// @description     Twitter風のSNSアプリケーションのREST API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Handlers
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	postHandler := handler.NewPostHandler()
	commentHandler := handler.NewCommentHandler()
	likeHandler := handler.NewLikeHandler()
	bookmarkHandler := handler.NewBookmarkHandler()
	followHandler := handler.NewFollowHandler()
	notificationHandler := handler.NewNotificationHandler()
	reportHandler := handler.NewReportHandler()

	// API routes
	api := e.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", authHandler.GetMe) // 要認証

	// User routes
	users := api.Group("/users")
	users.GET("", userHandler.SearchUsers)          // ユーザー検索
	users.GET("/:username", userHandler.GetProfile) // プロフィール取得
	users.PATCH("/me", userHandler.UpdateProfile)   // プロフィール更新（要認証）

	// Post routes
	posts := api.Group("/posts")
	posts.POST("", postHandler.CreatePost)       // 投稿作成（要認証）
	posts.GET("/:id", postHandler.GetPost)       // 投稿取得
	posts.PATCH("/:id", postHandler.UpdatePost)  // 投稿更新（要認証）
	posts.DELETE("/:id", postHandler.DeletePost) // 投稿削除（要認証）

	// Timeline
	api.GET("/timeline", postHandler.GetTimeline) // タイムライン（要認証）

	// User posts
	users.GET("/:username/posts", postHandler.GetUserPosts)

	// Comment routes
	posts.POST("/:id/comments", commentHandler.CreateComment)   // コメント作成（要認証）
	posts.GET("/:id/comments", commentHandler.GetComments)      // コメント一覧取得
	api.DELETE("/comments/:id", commentHandler.DeleteComment)   // コメント削除（要認証）

	// Like routes
	posts.POST("/:id/like", likeHandler.LikePost)             // 投稿いいね（要認証）
	posts.DELETE("/:id/like", likeHandler.UnlikePost)         // 投稿いいね解除（要認証）
	api.POST("/comments/:id/like", likeHandler.LikeComment)   // コメントいいね（要認証）
	api.DELETE("/comments/:id/like", likeHandler.UnlikeComment) // コメントいいね解除（要認証）

	// Bookmark routes
	posts.POST("/:id/bookmark", bookmarkHandler.AddBookmark)    // ブックマーク追加（要認証）
	posts.DELETE("/:id/bookmark", bookmarkHandler.RemoveBookmark) // ブックマーク削除（要認証）
	api.GET("/bookmarks", bookmarkHandler.GetBookmarks)         // ブックマーク一覧（要認証）

	// Follow routes
	users.POST("/:username/follow", followHandler.Follow)         // フォロー（要認証）
	users.DELETE("/:username/follow", followHandler.Unfollow)     // フォロー解除（要認証）
	users.GET("/:username/followers", followHandler.GetFollowers) // フォロワー一覧
	users.GET("/:username/following", followHandler.GetFollowing) // フォロー中一覧

	// Notification routes
	notifications := api.Group("/notifications")
	notifications.GET("", notificationHandler.GetNotifications)           // 通知一覧（要認証）
	notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)     // 通知既読（要認証）
	notifications.POST("/read-all", notificationHandler.MarkAllAsRead)   // 全通知既読（要認証）

	// Report routes
	api.POST("/reports", reportHandler.CreateReport) // 通報（要認証）

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
