module.exports = {
  apps: [
    {
      name: 'ecommitra-backend',
      cwd: './backend',
      // Note: For actual production, it's highly recommended to compile your Go app first!
      // 1. Run: `cd backend && go build -o ecommitra-api main.go`
      // 2. Change script to: './ecommitra-api'
      // 3. Remove 'interpreter' and 'interpreter_args'
      script: 'main.go',
      interpreter: 'go',
      interpreter_args: 'run',
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      env: {
        PORT: 8080,
      }
    },
    {
      name: 'ecommitra-frontend',
      cwd: './frontend',
      script: 'npm',
      args: 'run preview', // Make sure to run `npm run build` in the frontend directory before starting PM2!
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      env: {
        PORT: 5173,
      }
    }
  ]
};
