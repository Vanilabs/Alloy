.PHONY: build run start clean migrate-create migrate-up migrate-down migrate-status


export POSTGRES_DSN ?= $(shell grep POSTGRES_DSN backend/.env 2>/dev/null | cut -d '=' -f2- | tr -d '"')
export GOOSE_DRIVER ?= postgres
export GOOSE_DBSTRING ?= $(shell grep GOOSE_DBSTRING .envrc 2>/dev/null | cut -d '=' -f2- | tr -d '"' || echo "postgres://alloy:alloy_secret@localhost:5432/alloydb?sslmode=disable")

build:
	cd backend && go build -o bin/alloy ./cmd

run:
	cd backend && go run ./cmd

start: build
	cd backend && ./bin/alloy

clean:
	rm -rf backend/bin

migrate-create:
	cd backend && goose -dir internal/shared/database/migrations create $(name) sql

migrate-up:
	cd backend && GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir internal/shared/database/migrations up

migrate-down:
	cd backend && GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir internal/shared/database/migrations down

migrate-status:
	cd backend && GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir internal/shared/database/migrations status

test-backend:
	cd backend && go test ./tests/...

run-backend-stack-only:
	docker compose --profile alloy-api up --build

run-all-stacks:
	docker compose --profile alloy-api --profile frontend up --build