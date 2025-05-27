@echo off
echo Setting up Social Media App...

echo Creating .env file from template...
if not exist .env (
    copy .env.example .env
    echo .env file created. Please update with your settings.
) else (
    echo .env file already exists.
)

echo Installing Go dependencies...
go mod tidy

echo Starting Docker containers...
docker-compose up -d postgres redis

echo Waiting for database to be ready...
timeout /t 10 /nobreak

echo Setup completed!
echo.
echo Next steps:
echo 1. Update .env file with your configuration
echo 2. Run: docker-compose up -d
echo 3. Run: go run cmd/server/main.go
pause