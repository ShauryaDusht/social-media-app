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
- API will be available at: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

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