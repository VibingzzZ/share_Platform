.PHONY: dev test migrate

dev:
	go run ./cmd/server

test:
	go test ./...

migrate:
	go run ./cmd/migrate
