# Postly

A Twitter like Micro-Blogging website built with Go, PostgreSQL, and Redis.

## Quick Start

### Prerequisites
- Go 1.21+
- Docker
- PostgreSQL
- Redis

### Setup using Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/ShauryaDusht/go-social-media-app
   ```

2. Navigate to the project directory:
   ```bash
   cd go-social-media-app
   ```
3.  Build the Docker images:
   ```bash
   docker-compose build
   ```

4. Start the services:
   ```bash
   docker-compose up
   ```

5. Access the application at `http://localhost:8081`

6. Access Grafana for monitoring at `http://localhost:3000`

### Docker Commands
- Start services: `docker-compose up -d`
- Stop services: `docker-compose down`
- View logs: `docker-compose logs -f`
- Rebuild: `docker-compose up --build`

# Screenshots
![Profile Page](images/profile.png)
![Timeline Page](images/timeline.png)
![Rate Limiting](images/rateLimiting.png)
![Grafana Dashboard](images/grafana-dashboard.png)

---

# Project Status

This section outlines the progress and current status of the Micro-Blogging Project.

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
- Like/unlike posts - [BUG] : re-liking a post gives error
- User timeline logic - [DONE]

## Phase 5: Follow System and Caching - [DONE]
- Follow System - [DONE]
- Caching Timeline - [DONE]
- User profile search functionality - [DONE]

## Phase 6: Rate Limiting - [DONE]
- Rate limiting for APIs - [DONE]
- Use token bucket or fixed window (via Redis) - [DONE]
- Per user or IP — apply on post creation, likes, follow, etc - [DONE]

## Phase 7: Deployment and CI/CD - [TODO]
- Docker
- Nginx (for reverse proxy)

## Phase 8: Monitoring - [DONE]
- Prometheus metrics - [DONE]
- Grafana dashboard - [DONE]

## Phase 9: Testing - [TODO]
- Testing using go scripts

## Additional Features - [TODO]
- Add pagination to posts
- Improve search functionality by adding some fuzzyness
- Add feature to see followers and following list of a user
- Add comments to posts

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
├── images/
│   ├── profile.png
│   ├── rateLimiting.png
│   └── timeline.png
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
│   │   │   ├── cors.go
│   │   │   ├── metrics.go
│   │   │   └── rate_limit.go
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
├── monitoring/
│   ├── grafana/
│   │   └── dashboards/
│   │       └── api-metrics-dashboard.json
│   └── prometheus/
│       └── prometheus.yml
├── README.md
├── scripts/
│   ├── init.sql
│   ├── migrate.bat
│   ├── test_auth.bat
│   └── test_auth_with_token.bat
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
            ├── auth-check.js
            ├── auth.js
            ├── config.js
            ├── follows.js
            ├── index.js
            ├── posts.js
            ├── profile.js
            ├── search.js
            └── timeline.js
```