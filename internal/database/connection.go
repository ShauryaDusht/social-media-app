package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"social-media-app/internal/config"
	"social-media-app/internal/models"
)

var DB *gorm.DB

// Connect establishes database connection
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Like{},
		&models.Follow{},
	)
}

// AddConstraints adds custom database constraints
func AddConstraints() error {
	// Add unique constraint for likes (user can only like a post once)
	if err := DB.Exec(`
		ALTER TABLE likes 
		ADD CONSTRAINT unique_user_post_like 
		UNIQUE (user_id, post_id)
	`).Error; err != nil {
		log.Printf("Constraint already exists or error: %v", err)
	}

	// Add unique constraint for follows (user can only follow another user once)
	if err := DB.Exec(`
		ALTER TABLE follows 
		ADD CONSTRAINT unique_follower_following 
		UNIQUE (follower_id, following_id)
	`).Error; err != nil {
		log.Printf("Constraint already exists or error: %v", err)
	}

	// Add check constraint to prevent self-following
	if err := DB.Exec(`
		ALTER TABLE follows 
		ADD CONSTRAINT no_self_follow 
		CHECK (follower_id != following_id)
	`).Error; err != nil {
		log.Printf("Constraint already exists or error: %v", err)
	}

	return nil
}

// GetDB returns database instance
func GetDB() *gorm.DB {
	return DB
}
