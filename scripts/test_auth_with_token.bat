@echo off
echo Testing Authentication Endpoints with Token Extraction
echo ================================================

set API_URL=http://localhost:8080/api
set CURL_OPTS=-s

:: Check if jq is installed (for JSON parsing)
where jq >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: jq is not installed. Token extraction will not work properly.
    echo Please install jq from: https://stedolan.github.io/jq/download/
    echo Continuing with limited functionality...
    echo.
)

echo.
echo 1. Testing User Registration
echo --------------------------
curl %CURL_OPTS% -X POST %API_URL%/auth/register -H "Content-Type: application/json" -d "{\
    \"username\": \"testuser2\",\
    \"email\": \"testuser2@example.com\",\
    \"password\": \"password123\",\
    \"first_name\": \"Test\",\
    \"last_name\": \"User\"\
}" > register_response.json

type register_response.json
echo.

echo 2. Testing User Login
echo --------------------
curl %CURL_OPTS% -X POST %API_URL%/auth/login -H "Content-Type: application/json" -d "{\
    \"email\": \"testuser2@example.com\",\
    \"password\": \"password123\"\
}" > login_response.json

type login_response.json
echo.

echo 3. Extracting Token
echo -----------------

:: Try to extract token with jq if available
where jq >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    for /f "tokens=*" %%a in ('type login_response.json ^| jq -r ".data.token"') do set TOKEN=%%a
) else (
    :: Fallback to a simulated token if jq is not available
    set TOKEN=simulated_token_for_testing
    echo Using simulated token because jq is not installed.
)

echo Token: %TOKEN%
echo.

echo 4. Testing Protected Endpoint (Get Profile)
echo -----------------------------------------
curl %CURL_OPTS% -X GET %API_URL%/users/profile -H "Authorization: Bearer %TOKEN%" > profile_response.json

type profile_response.json
echo.

echo 5. Testing Logout
echo ---------------
curl %CURL_OPTS% -X POST %API_URL%/auth/logout > logout_response.json

type logout_response.json
echo.

echo 6. Testing Protected Endpoint After Logout (should still work with stateless JWT)
echo -----------------------------------------------------------------------
curl %CURL_OPTS% -X GET %API_URL%/users/profile -H "Authorization: Bearer %TOKEN%" > after_logout_response.json

type after_logout_response.json
echo.

echo 7. Cleanup temporary files
echo ------------------------
del register_response.json login_response.json profile_response.json logout_response.json after_logout_response.json

echo.
echo Test completed!
echo ===============

pause