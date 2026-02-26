.PHONY: build run start clean migrate-create migrate-up migrate-down migrate-status remote-deploy run-staging-stacks

REMOTE_USER=root
REMOTE_HOST=89.167.124.194
REMOTE_DIR=~/alloy

remote-deploy:
	@echo "🚀 Syncing code to remote server..."
	rsync -avzP --exclude='.git' --exclude='node_modules' --exclude='.next' --exclude='.gitignore' ./ $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_DIR)/
	
	@echo "🔄 Rebuilding and restarting containers on remote..."
	ssh $(REMOTE_USER)@$(REMOTE_HOST) "cd $(REMOTE_DIR) && make run-staging-stacks"
	
	@echo "✅ Deployment complete!"

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

run-local-backend-only:
	docker compose -f docker-compose.yml -f docker-compose-dev.yml --profile alloy-api up --build

run-all-local-stacks:
	docker compose -f docker-compose.yml -f docker-compose-dev.yml --profile alloy-api --profile alloy-ui up --build

install-be:
	cd ./backend && go mod download

install-fe:
	cd ./frontend && pnpm install

run-local-frontend-only:
	docker compose -f docker-compose.yml -f docker-compose-dev.yml --profile alloy-ui up --build


run-staging-stacks:
	docker compose -f docker-compose.yml -f docker-compose-staging.yml up -d --build