package main

import (
	"log"

	"social-media-app/internal/api"
	"social-media-app/internal/config"
	"social-media-app/internal/database"
)

func main() {
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
