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
	// connect tot db
	db := database.GetDB()
	// load config
	cfg := config.Load()

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	followRepo := repository.NewFollowRepository(db)
	cacheRepo := repository.NewCacheRepository(cfg)

	postService := services.NewPostService(postRepo, likeRepo, cacheRepo)
	likeService := services.NewLikeService(likeRepo, postRepo)
	followService := services.NewFollowService(followRepo, userRepo, cacheRepo)

	handlers.InitPostHandler(postService)
	handlers.LikePost(likeService)
	handlers.FollowUser(followService)

	// Connect to database
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

	// setup routes
	router := api.SetupRoutes(cfg)

	// start server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
