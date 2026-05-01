# Project Planner (ECOMMITRA)

This project features a Vue+Vite frontend and a Golang backend. It uses Supabase as its database.

## Prerequisites

Before running the application, ensure you have the following installed on your machine:
- [Node.js](https://nodejs.org/) (v16 or higher recommended) and npm
- [Golang](https://go.dev/dl/) (v1.20 or higher recommended)
- A [Supabase](https://supabase.com/) account and project.

## Environment Variables

You need to configure the connection to your Supabase instance. If `.env` files already exist, simply update the values.

**Frontend (`frontend/.env`):**
Create or update a `.env` file in the `frontend` folder with your Supabase URL and Anon Key:
```env
VITE_SUPABASE_URL=your_supabase_url
VITE_SUPABASE_ANON_KEY=your_supabase_anon_key
```

**Backend (`backend/.env`):**
Create or update a `.env` file in the `backend` folder with your Supabase URL, Secret/Service Role Key, and Port:
```env
SUPABASE_URL=your_supabase_url
SUPABASE_SECRET_KEY=your_supabase_secret_key
PORT=8080
```

## Database Setup

1. Open your Supabase project dashboard.
2. Navigate to the **SQL Editor**.
3. Copy the contents of the `supabase-setup.sql` file provided in the root directory.
4. Paste it into the SQL Editor and click **Run** to create the necessary tables and set up your database.

## Running the Application

### The Easy Way (Using Start Scripts)

We have provided scripts to easily start both the frontend and backend servers at once.

**For Windows Users:**
Double-click `start.bat` from your file explorer, or run it in your terminal:
```cmd
start.bat
```
This will open two new command prompt windows running the frontend and backend separately.

**For macOS/Linux Users:**
Run the bash script from your terminal:
```bash
chmod +x start.sh
./start.sh
```
This will run the backend in the background and the frontend in the foreground. Pressing `Ctrl+C` will cleanly stop both servers.

### The Manual Way

If you prefer to start the servers manually in separate terminal windows:

**1. Start the Backend (Golang)**
```bash
cd backend
go mod tidy
go run main.go
```
*The backend server usually runs on http://localhost:8080*

**2. Start the Frontend (Vue+Vite)**
```bash
cd frontend
npm install
npm run dev
```
*The frontend server usually runs on http://localhost:5173*

## Accessing the App
- **Frontend App:** http://localhost:5173
- **Backend API:** http://localhost:8080