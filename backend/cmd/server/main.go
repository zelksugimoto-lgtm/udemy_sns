package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yourusername/sns-app/internal/config"
	"github.com/yourusername/sns-app/internal/handler"
	appMiddleware "github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/repository"
	"github.com/yourusername/sns-app/internal/service"

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
	// データベース接続
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("データベース接続失敗: %v", err)
	}

	// マイグレーション実行
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("マイグレーション失敗: %v", err)
	}

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

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	followRepo := repository.NewFollowRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, followRepo, postRepo)
	postService := service.NewPostService(postRepo, userRepo, followRepo, likeRepo, bookmarkRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, userRepo, likeRepo)
	likeService := service.NewLikeService(likeRepo, postRepo, commentRepo)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, postRepo)
	followService := service.NewFollowService(followRepo, userRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	reportService := service.NewReportService(reportRepo, postRepo, commentRepo, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(likeService)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	followHandler := handler.NewFollowHandler(followService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	reportHandler := handler.NewReportHandler(reportService)

	// API routes
	api := e.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", authHandler.GetMe, appMiddleware.AuthMiddleware())

	// User routes
	users := api.Group("/users")
	users.GET("", userHandler.SearchUsers)
	users.GET("/:username", userHandler.GetProfile)
	users.PATCH("/me", userHandler.UpdateProfile, appMiddleware.AuthMiddleware())

	// Post routes
	posts := api.Group("/posts")
	posts.POST("", postHandler.CreatePost, appMiddleware.AuthMiddleware())
	posts.GET("/:id", postHandler.GetPost)
	posts.PATCH("/:id", postHandler.UpdatePost, appMiddleware.AuthMiddleware())
	posts.DELETE("/:id", postHandler.DeletePost, appMiddleware.AuthMiddleware())

	// Timeline
	api.GET("/timeline", postHandler.GetTimeline, appMiddleware.AuthMiddleware())

	// User posts
	users.GET("/:username/posts", postHandler.GetUserPosts)

	// Comment routes
	posts.POST("/:id/comments", commentHandler.CreateComment, appMiddleware.AuthMiddleware())
	posts.GET("/:id/comments", commentHandler.GetComments)
	api.DELETE("/comments/:id", commentHandler.DeleteComment, appMiddleware.AuthMiddleware())

	// Like routes
	posts.POST("/:id/like", likeHandler.LikePost, appMiddleware.AuthMiddleware())
	posts.DELETE("/:id/like", likeHandler.UnlikePost, appMiddleware.AuthMiddleware())
	api.POST("/comments/:id/like", likeHandler.LikeComment, appMiddleware.AuthMiddleware())
	api.DELETE("/comments/:id/like", likeHandler.UnlikeComment, appMiddleware.AuthMiddleware())

	// Bookmark routes
	posts.POST("/:id/bookmark", bookmarkHandler.AddBookmark, appMiddleware.AuthMiddleware())
	posts.DELETE("/:id/bookmark", bookmarkHandler.RemoveBookmark, appMiddleware.AuthMiddleware())
	api.GET("/bookmarks", bookmarkHandler.GetBookmarks, appMiddleware.AuthMiddleware())

	// Follow routes
	users.POST("/:username/follow", followHandler.Follow, appMiddleware.AuthMiddleware())
	users.DELETE("/:username/follow", followHandler.Unfollow, appMiddleware.AuthMiddleware())
	users.GET("/:username/followers", followHandler.GetFollowers)
	users.GET("/:username/following", followHandler.GetFollowing)

	// Notification routes
	notifications := api.Group("/notifications", appMiddleware.AuthMiddleware())
	notifications.GET("", notificationHandler.GetNotifications)
	notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
	notifications.POST("/read-all", notificationHandler.MarkAllAsRead)

	// Report routes
	api.POST("/reports", reportHandler.CreateReport, appMiddleware.AuthMiddleware())

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
