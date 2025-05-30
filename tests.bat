@echo off
echo Running tests...

:: Run the tests and save output to a temporary file
go test ./... -v > test_output.txt

:: Start the server in background
start /B go run main.go

:: Wait for server to start
timeout /t 5 /nobreak

:: Make a request to health endpoint and check response
curl -s -o nul -w "%%{http_code}" http://localhost:8080/health > health_status.txt
set /p STATUS=<health_status.txt
if not "%STATUS%"=="200" (
    echo Server health check failed
    exit /b 1
)

:: Count test failures
findstr /C:"--- FAIL" test_output.txt > test_failures.txt
for /f %%A in ('type test_failures.txt ^| find /c /v ""') do set FAIL_COUNT=%%A

echo.
echo Test Summary:
echo -------------
echo Failed Tests: %FAIL_COUNT%

:: Display failed tests if any
if %FAIL_COUNT% gtr 0 (
    echo.
    echo Failed Test Details:
    echo -------------------
    type test_failures.txt
)

:: Cleanup temporary files
del test_output.txt health_status.txt test_failures.txt

:: Stop the server
taskkill /F /IM go.exe 2>nul

if %FAIL_COUNT% gtr 0 (
    exit /b 1
) else (
    echo All tests passed successfully!
    exit /b 0
)