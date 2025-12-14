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

## QuizApp — Editor/Agent Guide

Short summary
- Architecture: Go (Gin) backend (:8080) + React (Vite) frontend (:5173) + Postgres. Docker Compose builds and wires services.

Quick dev commands
- Backend (local): `cd backend && go run .`
- Frontend (dev): `cd frontend && npm install && npm run dev`
- Full stack (docker): `docker compose up --build`

Key files & patterns
- Routes: [backend/main.go](backend/main.go) — Gin router and OPTIONS handling.
- Handlers: [backend/handlers.go](backend/handlers.go) — call `enableCorsHeaders(c)` at top of each handler; check for OPTIONS; use `c.BindJSON` and `c.JSON`.
- Models: [backend/models.go](backend/models.go) — `Quiz`, `Question`, `Option`, `Response`. Normalized DB tables; `Quiz` JSON field legacy handling exists.
- DB init: [backend/db.go](backend/db.go) — reads `DATABASE_URL` (defaults to the compose `db`), uses GORM `AutoMigrate`, and retries DB open with exponential backoff.
- Frontend fetch: [frontend/src/App.jsx](frontend/src/App.jsx) — hardcoded `http://localhost:8080/api/quizzes`; adjust when backend URL changes.

Project-specific conventions (do these exactly)
- CORS: Single helper `enableCorsHeaders(c *gin.Context)` in `handlers.go`; every handler must call it and handle `OPTIONS` requests.
- Preloading relations: handlers always use `db.Preload("Questions.Options")` before reads; preserve this pattern to avoid missing nested data.
- Legacy JSONB rows: `listQuizzes` merges legacy `data` JSONB into normalized `Questions` when empty — keep this compatibility logic when changing models.
- ID/CreatedAt behavior: POST handler (`createQuiz`) auto-generates `Quiz.ID` when empty and sets `CreatedAt` when zero. Preserve or intentionally modify both when changing creation flows.
- Schema changes: there is no migration tool — the code relies on `db.AutoMigrate(...)` in `initDB()`. Update `models.go` and expect AutoMigrate to run on startup.

Adding endpoints or schema fields
- Register new routes in [backend/main.go](backend/main.go) and implement handlers in [backend/handlers.go](backend/handlers.go).
- Add model fields in [backend/models.go](backend/models.go) and run the app (or `docker compose up`) to let `AutoMigrate` apply the simple schema changes. For complex migrations, plan manual SQLs (no migration framework present).

Docker / compose notes
- `docker-compose.yml` exposes backend `8080:8080` and frontend `5173:80` (nginx serves built frontend). Postgres uses `POSTGRES_*` envs; backend defaults its DSN to `postgres://postgres:password@db:5432/quizdb?sslmode=disable` if `DATABASE_URL` is not set.

Libraries & runtime
- Backend: Gin (HTTP), GORM (ORM), Postgres driver; vendor directory present.
- Frontend: React + Vite; production frontend image built into nginx via `frontend/Dockerfile`.

Testing & gaps
- There are no automated tests in the repo. Validate changes manually using the dev commands above or with `docker compose up`.

What to watch for when editing
- Do not remove `enableCorsHeaders` calls or OPTIONS handling — the frontend dev server expects CORS set to `http://localhost:5173`.
- Keep `Preload("Questions.Options")` so nested options load correctly.
- When changing JSON shape stored in DB (legacy `data` JSONB vs normalized tables), update the compatibility code in `listQuizzes`.

If anything is unclear or you'd like the guide expanded (examples for adding a route, sample request/response, or migration steps), tell me which part to expand.

