# Project Planner (ECOMMITRA)

[![Go Backend](https://img.shields.io/badge/Backend-Go_1.20+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Vue Frontend](https://img.shields.io/badge/Frontend-Vue_3_/_Vite-4FC08D?style=for-the-badge&logo=vuedotjs&logoColor=white)](https://vuejs.org/)
[![Database](https://img.shields.io/badge/Database-Supabase_Postgres-3ECF8E?style=for-the-badge&logo=supabase&logoColor=white)](https://supabase.com/)

A high-performance, interactive project management and architectural planning dashboard built with a **Golang** backend and **Vue 3** frontend. This platform features real-time data visualization, AI-ready logic, and automated schema management.

---

## 🚀 Quick Start (Automated)

We provide one-click scripts to launch both the frontend and backend simultaneously.

### **Windows**
Double-click `start.bat` in the root directory. This will:
1. Open a terminal for the **Backend** (port 8080).
2. Open a terminal for the **Frontend** (port 5173).

### **macOS / Linux**
Run the bash script:
```bash
chmod +x start.sh
./start.sh
```

---

## ⚙️ Environment Configuration

Before launching, ensure your `.env` files are correctly configured.

### **Backend** (`/backend/.env`)
The backend uses **GORM** to manage the PostgreSQL connection.
```env
PORT=8080
DATABASE_URL=postgresql://postgres:your_password@db.your_project.supabase.co:5432/postgres
JWT_SECRET=your_random_string_here
ADMIN_EMAIL=your_admin_email@example.com
ADMIN_ROLE=BABA
```

### **Frontend** (`/frontend/.env`)
```env
VITE_SUPABASE_URL=https://your_project.supabase.co
VITE_SUPABASE_ANON_KEY=your_anon_key
VITE_ADMIN_ROLE=BABA
```

---

## 🗄️ Database Strategy: "Schema-on-Boot"

This project utilizes **GORM Auto-Migrations**. You do **not** need to manually run SQL files or manage migrations through PGAdmin.

1. **Automatic Schema Sync:** When the backend starts, it automatically inspects the `models/` directory and updates the Supabase tables to match.
2. **Instant Seeding:** If the database is empty, the system automatically seeds the **Master ECOMMITRA Architecture Manifest**, including all core sections (Core Engine, AI Assistant, Inventory Matrix, etc.).
3. **Manual Override (Optional):** If you need to reset the database manually, you can use the `supabase-setup.sql` provided in the root, but this is rarely necessary.

---

## 🛠️ Manual Execution

If you prefer to run services individually:

### **1. Backend (Go)**
```bash
cd backend
go run main.go
```
*Server starts at `http://localhost:8080`*

### **2. Frontend (Vue 3)**
```bash
cd frontend
npm install
npm run dev
```
*Server starts at `http://localhost:5173`*

---

## 🌐 Production Deployment (PM2)

For Linux servers (AWS, EC2, Bitnami), use **PM2** to manage the process lifecycle.

### **1. Build the Frontend**
```bash
cd frontend
npm install
npm run build
```

### **2. Compile the Backend (Recommended)**
```bash
cd backend
go build -o ecommitra-api main.go
```

### **3. Start with PM2**
Use the provided `ecosystem.config.js` to manage both services:
```bash
# Install PM2 if needed: sudo npm install -g pm2
pm2 start ecosystem.config.js
pm2 save
pm2 startup
```

---

## 🏗️ Technical Architecture
- **State Management:** Reactive Vue Composition API.
- **Backend:** Standard Go `net/http` with high-performance routing.
- **ORM:** GORM (PostgreSQL) with connection pooling.
- **Middleware:** Custom CORS and JWT validation layers.
- **DevOps:** PM2 process management with auto-restart and log monitoring.