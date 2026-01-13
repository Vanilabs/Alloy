.PHONY: build run start clean migrate-create migrate-up migrate-down migrate-status


export POSTGRES_DSN ?= $(shell grep POSTGRES_DSN backend/.env 2>/dev/null | cut -d '=' -f2- | tr -d '"')

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
	cd backend && goose -dir internal/shared/database/migrations postgres "$$POSTGRES_DSN" up

migrate-down:
	cd backend && goose -dir internal/shared/database/migrations postgres "$$POSTGRES_DSN" down

migrate-status:
	cd backend && goose -dir internal/shared/database/migrations postgres "$$POSTGRES_DSN" status

test-backend:
	cd backend && go test ./tests/...

run-backend-stack-only:
	docker compose --profile backend up --build

run-all-stacks:
	docker compose --profile backend --profile frontend up --build