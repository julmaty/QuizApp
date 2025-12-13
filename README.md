# QuizApp (Go backend + React frontend)

This is a minimal starter scaffold for a simple quiz application.

Quick start (PowerShell):

1. Start the backend:

   cd .\backend; go run .

2. In a separate terminal, start the frontend dev server:

   cd .\frontend; npm install; npm run dev

Open the frontend at http://localhost:5173. The frontend will call the backend at http://localhost:8080/api/quizzes.

Notes
- Backend stores quizzes in memory (for prototyping). Add persistence for production.
- The backend allows CORS from the Vite dev origin.

Run with Docker Compose (build & serve production frontend via nginx)

```powershell
# Build images and start services
docker compose up --build

# Stop services
docker compose down
```

This composes the backend on port 8080 and serves the frontend at http://localhost:5173 (nginx mapped to host port 5173 so the backend's CORS settings continue to work).

Database
- Postgres is included in `docker-compose.yml` as the `db` service. Defaults used by the backend (and the compose file):
   - user: `postgres`
   - password: `password`
   - database: `quizdb`

The backend will create a `quizzes` table on startup and store each quiz as a JSONB record.
