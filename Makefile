.PHONY: dev test migrate web-test web-build

dev:
	go run ./cmd/server

test:
	go test ./...

web-test:
	cd web && npm test -- --run

web-build:
	cd web && npm run build

migrate:
	go run ./cmd/migrate
