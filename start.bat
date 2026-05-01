@echo off
echo Installing frontend dependencies (if they aren't already installed)...
cd frontend
call npm install
cd ..

echo.
echo Starting the Golang Backend Server...
start "ECOMMITRA Backend" cmd /c "cd backend && title ECOMMITRA Backend && color 0A && go mod tidy && go run main.go"

echo Starting the Vue+Vite Frontend Server...
start "ECOMMITRA Frontend" cmd /c "cd frontend && title ECOMMITRA Frontend && color 0B && npm run dev"

echo.
echo ========================================================
echo ✅ Both servers are starting up in separate windows!
echo.
echo The Frontend will usually run on: http://localhost:5173
echo The Backend will usually run on:  http://localhost:8080
echo ========================================================
echo.
pause
