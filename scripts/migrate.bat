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