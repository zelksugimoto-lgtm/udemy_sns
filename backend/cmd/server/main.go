package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yourusername/sns-app/internal/config"
	"github.com/yourusername/sns-app/internal/handler"
	appMiddleware "github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"
	appLogger "github.com/yourusername/sns-app/pkg/logger"

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
	// ロガー初期化
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	appLogger.Init(env)

	log.Info().Str("env", env).Msg("Starting SNS Application")

	// データベース接続
	db, err := config.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("データベース接続失敗")
	}
	log.Info().Msg("データベース接続成功")

	// マイグレーション実行
	if err := config.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("マイグレーション失敗")
	}
	log.Info().Msg("マイグレーション完了")

	e := echo.New()

	// Middleware（適用順序が重要）
	e.Use(middleware.Recover())                                   // パニック時のリカバリー
	e.Use(appMiddleware.RequestID())                              // リクエストID生成（最優先）
	e.Use(appMiddleware.AccessLog())                              // アクセスログ（リクエストID取得後）
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
	}))

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	followRepo := repository.NewFollowRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	reportRepo := repository.NewReportRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, refreshTokenRepo)
	userService := service.NewUserService(userRepo, followRepo, postRepo)
	postService := service.NewPostService(postRepo, userRepo, followRepo, likeRepo, bookmarkRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, userRepo, likeRepo, notificationService)
	likeService := service.NewLikeService(likeRepo, postRepo, commentRepo, notificationService)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, postRepo, likeRepo)
	followService := service.NewFollowService(followRepo, userRepo, notificationService)
	reportService := service.NewReportService(reportRepo, postRepo, commentRepo, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, postService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(likeService)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	followHandler := handler.NewFollowHandler(followService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	reportHandler := handler.NewReportHandler(reportService)

	// API routes
	api := e.Group("/api/v1")

	// Auth routes（認証系API: 5回/分のレート制限）
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register, appMiddleware.AuthRateLimitMiddleware())
	auth.POST("/login", authHandler.Login, appMiddleware.AuthRateLimitMiddleware())
	auth.POST("/refresh", authHandler.RefreshToken, appMiddleware.AuthRateLimitMiddleware())
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/revoke-all", authHandler.RevokeAllTokens, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())
	auth.GET("/me", authHandler.GetMe, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())

	// User routes（一般API: 60回/分のレート制限）
	users := api.Group("/users", appMiddleware.GeneralRateLimitMiddleware())
	users.GET("", userHandler.SearchUsers, appMiddleware.OptionalAuthMiddleware())
	users.GET("/:username", userHandler.GetProfile, appMiddleware.OptionalAuthMiddleware())
	users.PATCH("/me", userHandler.UpdateProfile, appMiddleware.AuthMiddleware())

	// Post routes（一般API: 60回/分のレート制限）
	posts := api.Group("/posts", appMiddleware.GeneralRateLimitMiddleware())
	posts.POST("", postHandler.CreatePost, appMiddleware.AuthMiddleware())
	posts.GET("/:id", postHandler.GetPost, appMiddleware.OptionalAuthMiddleware())
	posts.PATCH("/:id", postHandler.UpdatePost, appMiddleware.AuthMiddleware())
	posts.DELETE("/:id", postHandler.DeletePost, appMiddleware.AuthMiddleware())

	// Timeline（一般API: 60回/分のレート制限）
	api.GET("/timeline", postHandler.GetTimeline, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())

	// User posts（usersグループに含まれるため、既にレート制限が適用されている）
	users.GET("/:username/posts", postHandler.GetUserPosts, appMiddleware.OptionalAuthMiddleware())
	users.GET("/:username/likes", userHandler.GetUserLikedPosts, appMiddleware.OptionalAuthMiddleware())

	// Comment routes（postsグループに含まれるため、既にレート制限が適用されている）
	posts.POST("/:id/comments", commentHandler.CreateComment, appMiddleware.AuthMiddleware())
	posts.GET("/:id/comments", commentHandler.GetComments, appMiddleware.OptionalAuthMiddleware())

	// Comment like/delete routes（一般API: 60回/分のレート制限）
	api.DELETE("/comments/:id", commentHandler.DeleteComment, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())
	api.POST("/comments/:id/like", likeHandler.LikeComment, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())
	api.DELETE("/comments/:id/like", likeHandler.UnlikeComment, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())

	// Like routes（postsグループに含まれるため、既にレート制限が適用されている）
	posts.POST("/:id/like", likeHandler.LikePost, appMiddleware.AuthMiddleware())
	posts.DELETE("/:id/like", likeHandler.UnlikePost, appMiddleware.AuthMiddleware())

	// Bookmark routes（一般API: 60回/分のレート制限）
	posts.POST("/:id/bookmark", bookmarkHandler.AddBookmark, appMiddleware.AuthMiddleware())
	posts.DELETE("/:id/bookmark", bookmarkHandler.RemoveBookmark, appMiddleware.AuthMiddleware())
	api.GET("/bookmarks", bookmarkHandler.GetBookmarks, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())

	// Follow routes（usersグループに含まれるため、既にレート制限が適用されている）
	users.POST("/:username/follow", followHandler.Follow, appMiddleware.AuthMiddleware())
	users.DELETE("/:username/follow", followHandler.Unfollow, appMiddleware.AuthMiddleware())
	users.GET("/:username/followers", followHandler.GetFollowers, appMiddleware.OptionalAuthMiddleware())
	users.GET("/:username/following", followHandler.GetFollowing, appMiddleware.OptionalAuthMiddleware())

	// Notification routes（一般API: 60回/分のレート制限）
	notifications := api.Group("/notifications", appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())
	notifications.GET("", notificationHandler.GetNotifications)
	notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
	notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
	notifications.POST("/read-all", notificationHandler.MarkAllAsRead)

	// Report routes（一般API: 60回/分のレート制限）
	api.POST("/reports", reportHandler.CreateReport, appMiddleware.AuthMiddleware(), appMiddleware.GeneralRateLimitMiddleware())

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := e.Start(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
