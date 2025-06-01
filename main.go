package main

import (
	"log"
	"os"
	"path/filepath"

	"social-media-app/internal/api"
	"social-media-app/internal/config"
	"social-media-app/internal/database"
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

	// Load configuration
	cfg := config.Load()

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

	// Setup routes
	router := api.SetupRoutes(cfg)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
