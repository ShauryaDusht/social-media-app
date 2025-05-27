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