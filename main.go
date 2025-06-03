package main

import (
	"log"
	"os"
	"path/filepath"

	"social-media-app/internal/api"
	"social-media-app/internal/api/handlers"
	"social-media-app/internal/config"
	"social-media-app/internal/database"
	"social-media-app/internal/repository"
	"social-media-app/internal/services"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory:", err)
	} else {
		log.Println("Current working directory:", cwd)

		envPath := filepath.Join(cwd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			log.Println(".env file found at:", envPath)
		} else {
			log.Println(".env file not found at:", envPath, "Error:", err)
		}
	}

	// load config
	cfg := config.Load()

	// Connect to database first
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Add custom constraints
	if err := database.AddConstraints(); err != nil {
		log.Printf("Warning: Failed to add constraints: %v", err)
	}

	// Get the database instance after connection
	db := database.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	followRepo := repository.NewFollowRepository(db)
	cacheRepo := repository.NewCacheRepository(cfg)

	// Initialize services
	postService := services.NewPostService(postRepo, likeRepo, cacheRepo)
	likeService := services.NewLikeService(likeRepo, postRepo)
	followService := services.NewFollowService(followRepo, userRepo, cacheRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, cfg)

	// Initialize handlers
	handlers.InitPostHandler(postService)
	handlers.InitLikeHandler(likeService)
	handlers.InitFollowHandler(followService)
	handlers.InitUserHandler(userService)
	handlers.InitAuthHandler(authService)

	// setup routes
	router := api.SetupRoutes(cfg)

	// start server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
