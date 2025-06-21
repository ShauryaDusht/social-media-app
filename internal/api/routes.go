package api

import (
	"net/http"
	"time"

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

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg)

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
		// Authentication routes (rate limited by IP)
		auth := api.Group("/auth")
		{
			// More restrictive rate limiting for auth endpoints
			authRateLimit := middleware.CustomRateLimitConfig{
				Requests: 10, // 10 requests per hour for auth
				Window:   time.Hour,
			}

			auth.POST("/register",
				rateLimiter.CustomRateLimit("register", authRateLimit),
				handlers.Register,
			)
			auth.POST("/login",
				rateLimiter.CustomRateLimit("login", authRateLimit),
				handlers.Login,
			)
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
			// User routes (light rate limiting)
			users := protected.Group("/users")
			{
				users.GET("/profile", handlers.GetProfile)
				users.PUT("/profile",
					rateLimiter.RateLimitByUser("update_profile"),
					handlers.UpdateProfile,
				)
				users.GET("/:id", handlers.GetUserByID)
				users.GET("/search",
					rateLimiter.RateLimitByUser("search"),
					handlers.SearchUsers,
				)
			}

			// Post routes (moderate rate limiting)
			posts := protected.Group("/posts")
			{
				posts.GET("/", handlers.GetPosts)

				// Stricter rate limiting for post creation
				postCreateRateLimit := middleware.CustomRateLimitConfig{
					Requests: 20, // 20 posts per hour
					Window:   time.Hour,
				}
				posts.POST("/",
					rateLimiter.CustomRateLimit("create_post", postCreateRateLimit),
					handlers.CreatePost,
				)

				posts.GET("/:id", handlers.GetPostByID)
				posts.PUT("/:id",
					rateLimiter.RateLimitByUser("update_post"),
					handlers.UpdatePost,
				)
				posts.DELETE("/:id",
					rateLimiter.RateLimitByUser("delete_post"),
					handlers.DeletePost,
				)
				posts.GET("/user/:user_id", handlers.GetUserPosts)
			}

			// Like routes (stricter rate limiting to prevent spam)
			likes := protected.Group("/likes")
			{
				likeRateLimit := middleware.CustomRateLimitConfig{
					Requests: 60, // 60 likes per hour
					Window:   time.Hour,
				}

				likes.POST("/",
					rateLimiter.CustomRateLimit("like_post", likeRateLimit),
					handlers.LikePost,
				)
				likes.DELETE("/:post_id",
					rateLimiter.CustomRateLimit("unlike_post", likeRateLimit),
					handlers.UnlikePost,
				)
			}

			// Follow routes (moderate rate limiting)
			follows := protected.Group("/follows")
			{
				followRateLimit := middleware.CustomRateLimitConfig{
					Requests: 50, // 50 follows per hour
					Window:   time.Hour,
				}

				follows.POST("/",
					rateLimiter.CustomRateLimit("follow_user", followRateLimit),
					handlers.FollowUser,
				)
				follows.DELETE("/:user_id",
					rateLimiter.CustomRateLimit("unfollow_user", followRateLimit),
					handlers.UnfollowUser,
				)
				follows.GET("/followers/:user_id", handlers.GetFollowers)
				follows.GET("/following/:user_id", handlers.GetFollowing)
			}

			// Timeline route (light rate limiting)
			timeline := protected.Group("/timeline")
			{
				timeline.GET("/",
					rateLimiter.RateLimitByUser("timeline"),
					handlers.GetTimeline,
				)
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
