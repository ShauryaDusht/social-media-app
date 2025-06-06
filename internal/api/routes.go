package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"social-media-app/internal/api/handlers"
	"social-media-app/internal/api/middleware"
	"social-media-app/internal/config"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Social Media API is running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Please use POST method for login",
				})
			})
			auth.POST("/logout", handlers.Logout)
		}

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", handlers.GetProfile)
				users.PUT("/profile", handlers.UpdateProfile)
				users.GET("/:id", handlers.GetUserByID)
				users.GET("/search", handlers.SearchUsers)
			}

			// Post routes
			posts := protected.Group("/posts")
			{
				posts.GET("/", handlers.GetPosts)
				posts.POST("/", handlers.CreatePost)
				posts.GET("/:id", handlers.GetPostByID)
				posts.PUT("/:id", handlers.UpdatePost)
				posts.DELETE("/:id", handlers.DeletePost)
				posts.GET("/user/:user_id", handlers.GetUserPosts)
			}

			// Like routes
			likes := protected.Group("/likes")
			{
				likes.POST("/", handlers.LikePost)
				likes.DELETE("/:post_id", handlers.UnlikePost)
			}

			// Follow routes
			follows := protected.Group("/follows")
			{
				follows.POST("/", handlers.FollowUser)
				follows.DELETE("/:user_id", handlers.UnfollowUser)
				follows.GET("/followers/:user_id", handlers.GetFollowers)
				follows.GET("/following/:user_id", handlers.GetFollowing)
			}

			// Timeline route
			timeline := protected.Group("/timeline")
			{
				timeline.GET("/", handlers.GetTimeline)
			}
		}
	}

	// Serve static files (for uploaded images, CSS, JS)
	router.Static("/static", "./web/static")

	// Serve HTML files at root level
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/index.html", "./web/index.html")
	router.StaticFile("/login", "./web/login.html")
	router.StaticFile("/login.html", "./web/login.html")
	router.StaticFile("/signup", "./web/signup.html")
	router.StaticFile("/signup.html", "./web/signup.html")
	router.StaticFile("/posts", "./web/posts.html")
	router.StaticFile("/posts.html", "./web/posts.html")
	router.StaticFile("/profile", "./web/profile.html")
	router.StaticFile("/profile.html", "./web/profile.html")

	// Maintain backward compatibility for /web/ paths
	router.StaticFile("/web/login.html", "./web/login.html")
	router.StaticFile("/web/signup.html", "./web/signup.html")
	router.StaticFile("/web/posts.html", "./web/posts.html")
	router.StaticFile("/web/profile.html", "./web/profile.html")
	router.StaticFile("/web/index.html", "./web/index.html")

	return router
}
