#!/bin/bash

# Exit script if a command fails
set -e

echo "📦 Installing frontend dependencies (if they aren't already installed)..."
cd frontend
npm install
cd ..

echo ""
echo "🚀 Starting ECOMMITRA servers..."

# Start the Golang backend in the background
cd backend
go mod tidy
go run main.go &
BACKEND_PID=$!
echo "✅ Backend started (PID: $BACKEND_PID)"
cd ..

echo "✅ Frontend starting..."
echo "========================================================"
echo "Press Ctrl+C to stop both servers."
echo "========================================================"
echo ""

# Trap Ctrl+C (SIGINT) so that when you stop the script, it automatically kills the backend too
trap 'echo -e "\n🛑 Stopping servers..."; kill $BACKEND_PID; exit' SIGINT EXIT

# Start the Vue+Vite frontend in the foreground
cd frontend
npm run dev
