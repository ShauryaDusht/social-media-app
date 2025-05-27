# Phase 2: Models, Database Setup, and API Routes

## File 1: internal/config/config.go
```go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
	RateLimit RateLimitConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type ServerConfig struct {
	Host string
	Port string
	Env  string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse JWT expiry
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	// Parse rate limit window
	rateLimitWindow, err := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1h"))
	if err != nil {
		rateLimitWindow = time.Hour
	}

	// Parse rate limit requests
	rateLimitRequests, err := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	if err != nil {
		rateLimitRequests = 100
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "password123"),
			DBName:   getEnv("DB_NAME", "social_media"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			Expiry: jwtExpiry,
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENVIRONMENT", "development"),
		},
		RateLimit: RateLimitConfig{
			Requests: rateLimitRequests,
			Window:   rateLimitWindow,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

## File 2: internal/models/user.go
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"size:50"`
	LastName  string         `json:"last_name" gorm:"size:50"`
	Bio       string         `json:"bio" gorm:"size:500"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts     []Post   `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Likes     []Like   `json:"likes,omitempty" gorm:"foreignKey:UserID"`
	Followers []Follow `json:"followers,omitempty" gorm:"foreignKey:FollowingID"`
	Following []Follow `json:"following,omitempty" gorm:"foreignKey:FollowerID"`
}

// UserResponse is used for API responses (excludes sensitive data)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Bio:       u.Bio,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
	}
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
}

// LoginRequest represents user login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
```

## File 3: internal/models/post.go
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Content   string         `json:"content" gorm:"not null;size:1000"`
	ImageURL  string         `json:"image_url" gorm:"size:255"`
	LikeCount int            `json:"like_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User  User   `json:"user" gorm:"foreignKey:UserID"`
	Likes []Like `json:"likes,omitempty" gorm:"foreignKey:PostID"`
}

// PostResponse represents a post with user info for API responses
type PostResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	ImageURL  string       `json:"image_url"`
	LikeCount int          `json:"like_count"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserResponse `json:"user"`
	IsLiked   bool         `json:"is_liked"` // Whether current user liked this post
}

// CreatePostRequest represents the request to create a new post
type CreatePostRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ImageURL string `json:"image_url"`
}

// UpdatePostRequest represents the request to update a post
type UpdatePostRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ImageURL string `json:"image_url"`
}

// ToResponse converts Post to PostResponse
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        p.ID,
		Content:   p.Content,
		ImageURL:  p.ImageURL,
		LikeCount: p.LikeCount,
		CreatedAt: p.CreatedAt,
		User:      p.User.ToResponse(),
		IsLiked:   false, // This will be set by the service layer
	}
}
```

## File 4: internal/models/like.go
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
	Post Post `json:"post" gorm:"foreignKey:PostID"`
}

// Ensure unique constraint on user_id and post_id combination
func (Like) TableName() string {
	return "likes"
}

// LikeRequest represents a like/unlike request
type LikeRequest struct {
	PostID uint `json:"post_id" binding:"required"`
}
```

## File 5: internal/models/follow.go
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type Follow struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FollowerID  uint           `json:"follower_id" gorm:"not null;index"`  // User who follows
	FollowingID uint           `json:"following_id" gorm:"not null;index"` // User being followed
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Follower  User `json:"follower" gorm:"foreignKey:FollowerID"`
	Following User `json:"following" gorm:"foreignKey:FollowingID"`
}

// Ensure unique constraint on follower_id and following_id combination
func (Follow) TableName() string {
	return "follows"
}

// FollowRequest represents a follow/unfollow request
type FollowRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// FollowResponse represents follow relationship info
type FollowResponse struct {
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}
```

## File 6: internal/database/connection.go
```go
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
```

## File 7: internal/database/migrations/001_create_users.sql
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    bio VARCHAR(500),
    avatar VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
```

## File 8: internal/database/migrations/002_create_posts.sql
```sql
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL CHECK (LENGTH(content) <= 1000),
    image_url VARCHAR(255),
    like_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);
```

## File 9: internal/database/migrations/003_create_likes.sql
```sql
CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id, post_id)
);

CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id);
CREATE INDEX IF NOT EXISTS idx_likes_deleted_at ON likes(deleted_at);
```

## File 10: internal/database/migrations/004_create_follows.sql
```sql
CREATE TABLE IF NOT EXISTS follows (
    id SERIAL PRIMARY KEY,
    follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id);
CREATE INDEX IF NOT EXISTS idx_follows_deleted_at ON follows(deleted_at);
```

## File 11: internal/api/routes.go
```go
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
			}

			// Post routes
			posts := protected.Group("/posts")
			{
				posts.GET("/", handlers.GetPosts)              // Get timeline/all posts
				posts.POST("/", handlers.CreatePost)           // Create new post
				posts.GET("/:id", handlers.GetPostByID)        // Get specific post
				posts.PUT("/:id", handlers.UpdatePost)         // Update post
				posts.DELETE("/:id", handlers.DeletePost)      // Delete post
				posts.GET("/user/:user_id", handlers.GetUserPosts) // Get posts by user
			}

			// Like routes
			likes := protected.Group("/likes")
			{
				likes.POST("/", handlers.LikePost)      // Like a post
				likes.DELETE("/:post_id", handlers.UnlikePost) // Unlike a post
			}

			// Follow routes
			follows := protected.Group("/follows")
			{
				follows.POST("/", handlers.FollowUser)               // Follow a user
				follows.DELETE("/:user_id", handlers.UnfollowUser)   // Unfollow a user
				follows.GET("/followers/:user_id", handlers.GetFollowers) // Get user's followers
				follows.GET("/following/:user_id", handlers.GetFollowing) // Get who user follows
			}

			// Timeline route
			timeline := protected.Group("/timeline")
			{
				timeline.GET("/", handlers.GetTimeline) // Get personalized timeline
			}
		}
	}

	// Serve static files (for uploaded images, CSS, JS)
	router.Static("/static", "./web/static")

	return router
}
```

## File 12: scripts/migrate.bat
```batch
@echo off
echo Running database migrations...

echo Starting PostgreSQL container if not running...
docker-compose up -d postgres

echo Waiting for PostgreSQL to be ready...
timeout /t 10 /nobreak

echo Running migrations...
docker exec -i social_postgres psql -U admin -d social_media < internal/database/migrations/001_create_users.sql
docker exec -i social_postgres psql -U admin -d social_media < internal/database/migrations/002_create_posts.sql
docker exec -i social_postgres psql -U admin -d social_media < internal/database/migrations/003_create_likes.sql
docker exec -i social_postgres psql -U admin -d social_media < internal/database/migrations/004_create_follows.sql

echo Migrations completed!
pause
```

## File 13: cmd/server/main.go
```go
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
```

## File 14: internal/utils/response.go
```go
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "Validation failed",
		Error:   "Validation errors: " + joinErrors(errors),
	})
}

// joinErrors joins multiple error messages
func joinErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}
	
	result := errors[0]
	for i := 1; i < len(errors); i++ {
		result += ", " + errors[i]
	}
	return result
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *gin.Context, data interface{}, page, limit int, total int64) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Data retrieved successfully",
		Data: PaginationResponse{
			Data:       data,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
```

## File 15: internal/utils/jwt.go
```go
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the user
func GenerateToken(userID uint, username, email, secret string, expiry time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
```

## File 16: internal/utils/hash.go
```go
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

## File 17: internal/api/middleware/cors.go
```go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
```

## File 18: internal/api/middleware/auth.go
```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"social-media-app/internal/config"
	"social-media-app/internal/utils"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Load config to get JWT secret
		cfg := config.Load()

		// Validate token
		claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
```

## File 19: internal/api/handlers/auth.go (placeholder structure)
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-media-app/internal/utils"
)

// Register handles user registration
func Register(c *gin.Context) {
	// TODO: Implement user registration logic
	utils.ErrorResponse(c, http.StatusNotImplemented, "Registration endpoint not implemented yet")
}

// Login handles user login
func Login(c *gin.Context) {
	// TODO: Implement user login logic
	utils.ErrorResponse(c, http.StatusNotImplemented, "Login endpoint not implemented yet")
}
```

## File 20: internal/api/handlers/users.go (placeholder structure)
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-media-app/internal/utils"
)

// GetProfile gets current user's profile
func GetProfile(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get profile endpoint not implemented yet")
}

// UpdateProfile updates current user's profile
func UpdateProfile(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Update profile endpoint not implemented yet")
}

// GetUserByID gets user by ID
func GetUserByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get user by ID endpoint not implemented yet")
}
```

## File 21: internal/api/handlers/posts.go (placeholder structure)
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-media-app/internal/utils"
)

// GetPosts gets all posts/timeline
func GetPosts(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get posts endpoint not implemented yet")
}

// CreatePost creates a new post
func CreatePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Create post endpoint not implemented yet")
}

// GetPostByID gets a specific post
func GetPostByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get post by ID endpoint not implemented yet")
}

// UpdatePost updates a post
func UpdatePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Update post endpoint not implemented yet")
}

// DeletePost deletes a post
func DeletePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Delete post endpoint not implemented yet")
}

// GetUserPosts gets posts by a specific user
func GetUserPosts(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get user posts endpoint not implemented yet")
}

// GetTimeline gets personalized timeline
func GetTimeline(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get timeline endpoint not implemented yet")
}
```

## File 22: internal/api/handlers/likes.go (placeholder structure)
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-media-app/internal/utils"
)

// LikePost likes a post
func LikePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Like post endpoint not implemented yet")
}

// UnlikePost unlikes a post
func UnlikePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Unlike post endpoint not implemented yet")
}
```

## File 23: internal/api/handlers/follows.go (placeholder structure)
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-media-app/internal/utils"
)

// FollowUser follows a user
func FollowUser(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Follow user endpoint not implemented yet")
}

// UnfollowUser unfollows a user
func UnfollowUser(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Unfollow user endpoint not implemented yet")
}

// GetFollowers gets user's followers
func GetFollowers(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get followers endpoint not implemented yet")
}

// GetFollowing gets who user follows
func GetFollowing(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get following endpoint not implemented yet")
}
```

## File 24: scripts/init.sql
```sql
-- Initialize database with some basic setup
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$ language 'plpgsql';
```

## File 25: test-api.bat (for testing the API)
```batch
@echo off
echo Testing Social Media API...

echo.
echo Testing health endpoint...
curl -X GET http://localhost:8080/health

echo.
echo.
echo Testing registration endpoint...
curl -X POST http://localhost:8080/api/auth/register

echo.
echo.
echo API test completed!
pause
```