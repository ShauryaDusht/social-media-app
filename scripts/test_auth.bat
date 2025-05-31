@echo off
echo Testing Authentication Endpoints
echo ==============================

set API_URL=http://localhost:8080/api
set CURL_OPTS=-s -w "\n"

echo.
echo 1. Testing User Registration
echo --------------------------
curl %CURL_OPTS% -X POST %API_URL%/auth/register -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"email\":\"testuser@example.com\",\"password\":\"password123\",\"first_name\":\"Test\",\"last_name\":\"User\"}"

echo.
echo 2. Testing User Login
echo --------------------
for /f "tokens=*" %%a in ('curl %CURL_OPTS% -X POST %API_URL%/auth/login -H "Content-Type: application/json" -d "{\"email\":\"testuser@example.com\",\"password\":\"password123\"}"') do (
    echo Response: %%a
    set RESPONSE=%%a
)

echo.
echo 3. Extracting Token (simulated - in a real script you would parse the JSON)
echo ---------------------------------------------------------------------
echo Simulating token extraction from response...
set "TOKEN=%RESPONSE:~11,24%"
echo Token: %TOKEN%

echo.
echo 4. Testing Protected Endpoint (Get Profile)
echo -----------------------------------------
curl %CURL_OPTS% -X GET %API_URL%/users/profile -H "Authorization: Bearer %TOKEN%"

echo.
echo 5. Testing Logout
echo ---------------
curl %CURL_OPTS% -X POST %API_URL%/auth/logout

echo.
echo 6. Testing Protected Endpoint After Logout (should fail)
echo -----------------------------------------------------
curl %CURL_OPTS% -X GET %API_URL%/users/profile -H "Authorization: Bearer %TOKEN%"

echo.
echo Test completed!
echo ===============

pause
