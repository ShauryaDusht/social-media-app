# Social Media App

A modern social media platform built with Go, PostgreSQL, and Redis.

## Quick Start

### Prerequisites
- Go 1.21+
- Docker Desktop
- Git

### Setup
1. Clone the repository
2. Run setup script: `setup.bat`
3. Update `.env` file with your configuration
4. Start the application: `start.bat`

### Development
- API will be available at: `http://localhost:8080`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

### API Endpoints
- `GET /health` - Health check
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/posts` - Get posts
- `POST /api/posts` - Create post

### Testing
Run tests: `go test ./...`

### Docker Commands
- Start services: `docker-compose up -d`
- Stop services: `docker-compose down`
- View logs: `docker-compose logs -f`
- Rebuild: `docker-compose up --build`

---

# Micro-Blogging Website Project Status

This section outlines the progress and current status of the Micro-Blogging Website Project.

## Phase 1: Environment Setup - [DONE]
- Set up environment
- Docker
- Project directory setup

## Phase 2: Models, Database, and API Routes - [DONE]
- Add models
- Set up PostgreSQL Redis locally
- Define REST API routes
- Make migrations

## Phase 3: Authentication and UI - [DONE]
- Add login, signup, logout
- JWT based only
- Add basic UI for login, signup, and posts

## Phase 4: Post APIs - [IN PROGRESS]
- CRUD APIs for posts - [DONE]
- Like/unlike posts - [BUG] : likes not registering
- User timeline logic - [DONE]

## Phase 5: Follow System and Caching - [TODO]
- Follow System - [BUG] : refollow not working
- Caching Timeline - [TODO]

## Phase 6: Rate Limiting - [TODO]
- Rate limiting for APIs
- Use token bucket or fixed window (via Redis)
- Per user or IP — apply on post creation, likes, follow, etc

## Phase 7: Deployment and CI/CD - [TODO]
- Docker
- Nginx
- CI/CD

## Phase 8: Monitoring - [TODO]
- Prometheus metrics
- Grafana dashboard

## Phase 9: Testing - [TODO]
- Testing using python/go scripts

## Additional Features - [TODO]
- Real-time notifications via WebSockets for likes and follows
- Search functionality for users, posts, and hashtags

---

# Directory Structure


```
social-media-app/
├── deployments/
│   └── docker/
│       └── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── follows.go
│   │   │   ├── likes.go
│   │   │   ├── posts.go
│   │   │   └── users.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   └── cors.go
│   │   └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migrations/
│   │       ├── 001_create_users.sql
│   │       ├── 002_create_post.sql
│   │       ├── 003_create_likes.sql
│   │       └── 004_create_follow.sql
│   ├── models/
│   │   ├── follow.go
│   │   ├── like.go
│   │   ├── post.go
│   │   ├── update_profile.go
│   │   └── user.go
│   ├── repository/
│   │   ├── cache_repo.go
│   │   ├── follow_repo.go
│   │   ├── interfaces.go
│   │   ├── like_repo.go
│   │   ├── post_repo.go
│   │   └── user_repo.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── follow_service.go
│   │   ├── like_service.go
│   │   ├── post_service.go
│   │   └── user_service.go
│   └── utils/
│       ├── hash.go
│       ├── jwt.go
│       └── response.go
├── LICENSE
├── main.go
├── README.md
├── scripts/
│   ├── init.sql
│   ├── migrate.bat
│   ├── test_auth.bat
│   └── test_auth_with_token.bat
├── tree.py
└── web/
    ├── index.html
    ├── login.html
    ├── posts.html
    ├── profile.html
    ├── signup.html
    └── static/
        ├── css/
        │   └── styles.css
        ├── img/
        │   └── default-avatar.png
        └── js/
            ├── auth.js
            ├── config.js
            ├── follows.js
            ├── index.js
            ├── posts.js
            ├── profile.js
            ├── search.js
            └── timeline.js
```