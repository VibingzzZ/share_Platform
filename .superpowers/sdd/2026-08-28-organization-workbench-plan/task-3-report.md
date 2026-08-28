# Task 3 Report: Content and Layout APIs

## Implemented

- Added content models for resources, posts, AI products, dashboards, pagination, list filters, and statuses.
- Added layout defaults and validation. Layouts require each supported module exactly once, restrict hidden entries to supported modules, and only allow `compact`/`comfortable` density and `light`/`dark` theme values.
- Added PostgreSQL repositories using parameterized `pgx` queries. Public list and dashboard queries are constrained to `status = 'published'`; list queries support pagination and the documented resource type/tag and tag filters.
- Added layout JSONB load/default and upsert persistence scoped to the request user ID.
- Added development identity middleware using `X-User-ID` and `X-User-Role`, with explicit identity and administrator guards.
- Registered `/api/v1/dashboard`, public paginated content routes, authenticated layout routes, and administrator CRUD routes for resources, posts, and AI products.
- Wired the server bootstrap to create a PostgreSQL pool and inject content/layout repositories.
- Added HTTP tests for public pagination/filter propagation and the administrator authorization boundary, plus model tests for valid layout customization and unknown module rejection.

## Verification

- `go test ./...` passed.
- `go vet ./...` passed.
- `git diff --check` passed.

## Notes

The existing database migration test skips database integration only when `DATABASE_URL` is unset. The API repositories are injected through `internal/http.Dependencies`, allowing handler tests to run without PostgreSQL while production startup uses the configured PostgreSQL pool.
