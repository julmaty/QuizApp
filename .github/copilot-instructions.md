## QuizApp: Go backend + React Vite frontend

A full-stack quiz application: Go REST API + Postgres + React frontend.
**GitHub Repository**: https://github.com/julmaty/QuizApp

### Architecture overview

**Data flow**: Frontend (React) → Backend (Go HTTP mux + Postgres) → Postgres
- Frontend hardcodes backend URL: `http://localhost:8080/api/quizzes` ([frontend/src/App.jsx](../frontend/src/App.jsx))
- Backend serves on `:8080` ([backend/main.go](../backend/main.go))
- Postgres is the source of truth for quizzes; JSONB column stores serialized `Quiz` structs
- CORS is hardcoded to allow only Vite dev origin: `http://localhost:5173` ([backend/handlers.go](../backend/handlers.go))

**Key data structures** ([backend/models.go](../backend/models.go)):
- `Quiz`: ID (auto-generated from title + timestamp), Title, Questions (slice), CreatedAt
- `Question`: Text, Options (slice), Multiple (boolean), Answers (indices)

**API endpoints** (all return JSON; POST/GET support OPTIONS for CORS):
- `GET /api/quizzes` → list all quizzes (ordered by created_at DESC)
- `POST /api/quizzes` → create quiz (auto-generates ID if empty, sets CreatedAt if zero)
- `GET /api/quizzes/{id}` → fetch single quiz by ID (see [backend/handlers.go](../backend/handlers.go) for getQuiz)
- `GET /health` → returns "ok"

### Dev commands (PowerShell from repo root)

**Local dev** (backend + frontend on separate ports):
```powershell
# Terminal 1: Start backend (listening on :8080)
cd .\backend; go run .

# Terminal 2: Start frontend Vite dev server (on :5173)
cd .\frontend; npm install; npm run dev
```

**Production with Docker Compose**:
```powershell
docker compose up --build
# Services: backend (:8080), frontend nginx (:5173), postgres (internal)
docker compose down
```

**Frontend build only** (outputs to `frontend/dist`):
```powershell
cd .\frontend; npm run build
```

### Backend patterns

**Database initialization** ([backend/db.go](../backend/db.go)):
- Reads `DATABASE_URL` env var; defaults to `postgres://postgres:password@db:5432/quizdb?sslmode=disable`
- Auto-creates `quizzes` table (id TEXT PK, title, data JSONB, created_at)
- Implements exponential backoff (500ms → 5s) to wait for Postgres container readiness

**HTTP handler pattern** ([backend/handlers.go](../backend/handlers.go)):
- All handlers call `enableCors(w)` before processing; each handler checks for OPTIONS method
- POST handlers auto-generate Quiz IDs if empty: `{Title}-{timestamp}`
- All POST/GET responses are JSON-encoded; errors use `http.Error()`
- Database errors are logged via `log.Printf()` before returning HTTP errors

**Main.go wiring** ([backend/main.go](../backend/main.go)):
- Calls `initDB()` on startup to connect and auto-create schema
- Uses `http.NewServeMux()` with explicit route registration (not a third-party router)
- Wraps mux with a `logger` middleware before starting server

### Frontend patterns

**React component structure** ([frontend/src](../frontend/src)):
- `App.jsx` is entry point; manages quiz list state and fetches on mount
- `CreateQuiz.jsx` → form to create quiz; calls `onCreated` callback after success
- `QuizList.jsx` → displays quizzes; likely renders `QuizItem` or similar
- Uses direct `fetch()` to call backend; no API client abstraction yet

**Vite config** ([frontend/package.json](../frontend/package.json)):
- Scripts: `dev` (start dev server), `build` (vite build), `preview` (preview prod)

### Docker Compose topology

Service dependencies and ports ([docker-compose.yml](../docker-compose.yml)):
- `db` (Postgres 15): internal only; health check ensures readiness before backend starts
- `backend`: depends on db; port 8080 → host 8080
- `frontend`: depends on backend; port 80 (nginx) → host 5173; nginx serves prod build
- Volume `db-data` persists Postgres data

Backend env inside compose: `DATABASE_URL` not set, so uses default pointing to `db:5432`
Frontend nginx is built from [frontend/Dockerfile](../frontend/Dockerfile); should listen on :80 internally

### Conventions & change patterns

**When adding features**:
- New HTTP endpoints go in [backend/main.go](../backend/main.go) (route registration) + [backend/handlers.go](../backend/handlers.go) (implementation)
- New Quiz fields require updates to [backend/models.go](../backend/models.go) AND the table schema in [backend/db.go](../backend/db.go)
- Frontend components should call `http://localhost:8080/api/quizzes` directly; no abstraction layer yet
- Ensure POST handlers update both `QuizCreatedAt` **and** generate IDs if missing (see [backend/handlers.go](../backend/handlers.go))

**CORS gotchas**:
- Every handler must include `enableCors(w)` call
- Hardcoded to `http://localhost:5173`; must change when frontend URL changes
- Frontend Vite dev server runs on :5173; production nginx also mapped to :5173 in compose

**Database schema**:
- Quizzes stored as JSONB; no normalization of Questions
- Schema auto-created on backend startup; no migrations tool yet

