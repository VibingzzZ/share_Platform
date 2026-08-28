# Organization Workbench Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a runnable internal organization workbench with Go/PostgreSQL APIs, seeded content, and a customizable frontend dashboard.

**Architecture:** A Go HTTP service owns migrations, seed data, content APIs, and layout persistence. A small frontend application consumes the API and provides dashboard modules plus a DIY drawer; a development fallback keeps the page usable when the API is unavailable.

**Tech Stack:** Go 1.22+, chi router, pgx, PostgreSQL, vanilla TypeScript/CSS/HTML, Docker Compose, Go testing.

**Spec:** `docs/superpowers/specs/2026-08-28-organization-workbench-design.md`

## Global Constraints

- Use PostgreSQL; do not introduce SQLite as a runtime dependency.
- Keep resource files under `resources/` and store relative paths in PostgreSQL.
- Expose versioned REST endpoints under `/api/v1` with JSON errors and paginated list envelopes.
- Layout JSON must validate module keys, order, visibility, density, and theme.
- Ordinary members read published content and save their own layout; admins manage content.

### Task 1: Repository and Service Scaffolding

**Files:**
- Create: `go.mod`, `cmd/server/main.go`, `internal/http/router.go`, `internal/config/config.go`
- Create: `docker-compose.yml`, `.env.example`, `Makefile`
- Test: `internal/http/router_test.go`

**Interfaces:** `NewRouter(cfg Config, deps Dependencies) http.Handler`; `GET /api/v1/health` returns `{status:"ok"}`.

- [ ] Write a router test that requests `/api/v1/health` and asserts HTTP 200 JSON.
- [ ] Run `go test ./internal/http/...` and verify the test fails because the module is absent.
- [ ] Add Go module dependencies, config loading, server bootstrap, router, health handler, and graceful shutdown.
- [ ] Add Docker Compose PostgreSQL service and Make targets for `dev`, `test`, and `migrate`.
- [ ] Run `go test ./...` and `go vet ./...`.

### Task 2: Database Migrations and Seed Data

**Files:**
- Create: `db/migrations/001_init.sql`, `db/migrations/002_indexes.sql`, `db/seed/001_demo.sql`
- Create: `internal/db/migrate.go`, `internal/db/seed.go`
- Create: `resources/README.md`
- Test: `internal/db/migrate_test.go`

**Interfaces:** migration runner applies ordered SQL files; seed is idempotent and creates demo member, resources, posts, and Dify product.

- [ ] Write migration tests against `DATABASE_URL`, skipping only when it is unset.
- [ ] Add tables `users`, `resources`, `posts`, `ai_products`, and `user_layouts` with UUID keys, timestamps, publish status, and JSONB layout.
- [ ] Add indexes for status/date and GIN indexes for tags where applicable.
- [ ] Add idempotent seed records and a sample resource-folder contract.
- [ ] Run migration and seed tests with `docker compose up -d db`.

### Task 3: Content and Layout APIs

**Files:**
- Create: `internal/model/models.go`, `internal/repository/content.go`, `internal/repository/layout.go`
- Create: `internal/http/handlers_content.go`, `internal/http/handlers_layout.go`, `internal/http/middleware.go`
- Modify: `internal/http/router.go`
- Test: `internal/http/handlers_test.go`, `internal/model/layout_test.go`

**Interfaces:** `GET /api/v1/dashboard`; paginated `GET /api/v1/resources`, `/posts`, `/ai-products`; `GET/PUT /api/v1/layout`; admin CRUD endpoints; `X-User-ID` development identity and `X-User-Role` role adapter.

- [ ] Write failing tests for pagination/filtering, published-only reads, admin authorization, and rejecting unknown layout modules.
- [ ] Implement model validation and repository queries using pgx parameterized SQL.
- [ ] Implement handlers with consistent `{items,page,pageSize,total}` envelopes and `{code,message}` errors.
- [ ] Implement layout upsert with defaults and JSONB validation.
- [ ] Run `go test ./...` and `go vet ./...`.

### Task 4: Frontend Workbench and DIY Layout

**Files:**
- Create: `web/index.html`, `web/styles.css`, `web/app.ts`, `web/api.ts`, `web/types.ts`
- Modify: `cmd/server/main.go` to serve `web/` and static resources
- Test: `web/app.test.ts` (or browser smoke script)

**Interfaces:** `loadDashboard()`, `loadLayout()`, `saveLayout(layout)`, and UI state for module order/visibility, density, and theme.

- [ ] Add a browser smoke test covering initial dashboard render and DIY save/restore.
- [ ] Build an internal-workbench layout with top bar, overview, resource list, development log, AI lab, and DIY drawer.
- [ ] Connect API calls and use seeded fallback data for offline development.
- [ ] Persist layout through `PUT /api/v1/layout`, with localStorage fallback and responsive mobile styles.
- [ ] Verify keyboard focus, readable contrast, and no horizontal overflow at desktop/mobile widths.

### Task 5: Documentation and Verification

**Files:**
- Modify: `README.md`
- Create: `docs/api.md`

- [ ] Document prerequisites, environment variables, Docker startup, migrations, seed data, and API examples.
- [ ] Run `docker compose config`, `go test ./...`, `go vet ./...`, and `git diff --check`.
- [ ] Start the service, request health/dashboard endpoints, and perform the frontend smoke check.
