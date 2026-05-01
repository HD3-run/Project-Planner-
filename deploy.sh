#!/bin/bash
echo "🚀 Deploying latest updates..."

# 1. Pull latest code
git fetch origin main
git reset --hard origin/main
git clean -fd

# 2. Update Backend
echo "📦 Updating Backend..."
cd backend
go mod tidy
cd ..

# 3. Update Frontend
echo "🎨 Building Frontend..."
cd frontend
npm install
npm run build
cd ..

# 4. Restart PM2
echo "🔄 Restarting servers..."
pm2 restart all

echo "✅ Deployment complete!"
