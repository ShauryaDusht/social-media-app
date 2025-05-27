@echo off
echo Stopping Social Media App...

echo Stopping Docker containers...
docker-compose down

echo All services stopped.
pause