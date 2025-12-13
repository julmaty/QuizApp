## Repo snapshot

This repository currently has no code or discoverable project files. Treat this file as a living guide that should be updated when sources are added. The goal: give an AI coding agent immediate, repository-specific actions to become productive.

## What to do first (automated discovery checklist)
- Look for these top-level artifacts in this order and report results back: `package.json`, `pyproject.toml`, `requirements.txt`, `setup.py`, `go.mod`, `pom.xml`, `build.gradle`, `Cargo.toml`, `Dockerfile`, `Makefile`, `README.md`, `src/`, `cmd/`, `internal/`, `pkg/`, `services/`, `tests/`.
- If none of the above exist, ask the human: "What language/framework should I assume?" and request any README or target files to proceed.

## How to infer the "big picture" once code exists
- Identify the entry points: search for `main(` (Go, Python, C/C++), `if __name__ == "__main__"` (Python), `package main` (Go), `public static void main` (Java), `src/index.js`/`server.js`/`app.ts` (Node). These files often show service boundaries and startup arguments.
- Look for service wiring files: `docker-compose.yml`, `k8s/` manifests, `deploy/`, `infrastructure/`. They reveal multi-service topology and env var contracts.
- Map data flow by locating API surface files (controllers/handlers/routes) and persistence layers (migrations/, models/, db/). Example patterns to note: `controllers/*` calling `services/*` calling `repositories/*`.

## Developer workflows (what to try first)
- Build: attempt language-specific quick checks in this order if files are found:
  - Node: `npm ci && npm run build` then `npm test`.
  - Python: create venv, `pip install -r requirements.txt`, run `pytest -q`.
  - Go: `go test ./...` then `go build ./...`.
  - Java/Gradle/Maven: run `mvn -q -DskipTests package` or `./gradlew build`.
  - Docker: `docker build -t repo-temp .` if `Dockerfile` present.
- Tests: run the project's test command from `package.json`, `pyproject.toml`, or CI config. If none, run sensible defaults above and report failures and missing config.
- Debugging: prefer adding lightweight logging near suspected entry points. If the project uses containers, suggest `docker run -e LOG_LEVEL=debug` or `docker-compose up`.

## Project-specific conventions & patterns (how to write code here)
- If present, follow the repository's top-level folder layout. Prefer minimal invasive changes: add files under `src/` or `services/<name>/` and update corresponding build configs.
- Tests are expected alongside code (same package/dir) unless a `tests/` folder exists. When adding tests, prefer the project's test runner and mirror its naming patterns (e.g., `*_test.go`, `test_*.py`, `.spec.ts`).

## Integration points & external deps
- If CI/CD config exists (e.g., `.github/workflows/`, `.gitlab-ci.yml`), read it before making changes — it documents build/test commands and environment variables used in CI.
- Look for `.env`, `.env.example`, or secrets references in workflows to infer required env vars. Do NOT attempt to read secrets; instead, ask the user to provide sanitized values or a local dev env.

## When you add/modify files: what to include in PRs
- Keep changes small and focused. Each PR should include:
  - A short description of the change and why (1–2 lines).
  - Commands to reproduce locally (build/test commands used).
  - A short note if new env vars or services are required.

## Examples (what I'll call out automatically)
- "Found Node project: `package.json` with scripts `start`, `test` — will run `npm ci` then `npm test`."
- "Found Dockerfile + k8s manifests — will map image names and `envFrom` to identify runtime config."

## If you (human) are reading this
- If you want a tuned instruction file for a specific language or architecture, add a README or a sample project file (one of the artifacts in the discovery checklist) and ask me to regenerate this guide — I'll merge and preserve any existing text in this file.

---
Please tell me which language/framework you'll use or add a README/code and I'll update this file to include repository-specific commands and examples. Do you want me to initialize a simple Node or Python starter layout now?

## Project: QuizApp (Go backend + React Vite frontend)

QuizApp is a full-stack quiz application with a persistent Postgres backend and a React frontend served via Vite/nginx.

**GitHub Repository**: https://github.com/julmaty/QuizApp

Key files added by the scaffolder:
- `backend/main.go`, `backend/handlers.go`, `backend/models.go` — Go HTTP server, in-memory store, endpoints:
  - GET /api/quizzes — list quizzes
  - POST /api/quizzes — create quiz
  - GET /health — health check
- `frontend/` — Vite + React app with `src/` components: `CreateQuiz.jsx`, `QuizList.jsx`, `App.jsx`.

Run & dev commands (PowerShell)
```powershell
# Run backend (Go)
cd .\backend; go run .

# Run frontend dev server (Vite)
cd ..\frontend; npm install; npm run dev
```

Build for production (frontend)
```powershell
# from repo root
cd .\frontend; npm run build
```

Notes on architecture & patterns
- Backend is intentionally minimal: an in-memory slice protected by a mutex. It's suitable for prototyping and local dev only. Replace with a DB and repository layer for persistence.
- Frontend uses Vite dev server (http://localhost:5173) and talks to backend at http://localhost:8080. CORS is allowed from the Vite origin in the backend handlers.

What to change when the project grows
- Replace `store` in `backend/handlers.go` with a repository interface and implementations (in-memory, postgres). Add models and migrations under `migrations/`.
- Add an API client module in the frontend (`src/api`) to centralize fetch calls and error handling.

If you'd like, I can:
- Replace the in-memory store with SQLite + `database/sql` and add migration tooling.
- Add Dockerfiles and a `docker-compose.yml` to run frontend + backend together.

