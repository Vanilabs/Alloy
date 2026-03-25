.PHONY: generate-openapi install-be install-fe test-backend migrate-create migrate-up migrate-down migrate-status run-frontend-local run-backend-local run-all-local run-all-staging run-postgres-only remote-deploy

REMOTE_USER=root
REMOTE_HOST=89.167.124.194
REMOTE_DIR=~/alloy
BACKEND_ENV=backend/.env

remote-deploy: generate-openapi
	@echo "🚀 Syncing code to remote server..."
	rsync -avzP --exclude='.git' --exclude='node_modules' --exclude='.next' --exclude='.gitignore' ./ $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_DIR)/
	
	@echo "🔄 Rebuilding and restarting containers on remote..."
	ssh $(REMOTE_USER)@$(REMOTE_HOST) "cd $(REMOTE_DIR) && make run-all-staging"
	
	@echo "✅ Deployment complete!"

export POSTGRES_DSN ?= $(shell grep POSTGRES_DSN backend/.env 2>/dev/null | cut -d '=' -f2- | tr -d '"')
export GOOSE_DRIVER ?= postgres
export GOOSE_DBSTRING ?= $(shell grep GOOSE_DBSTRING .envrc 2>/dev/null | cut -d '=' -f2- | tr -d '"' || echo "postgres://alloy:alloy_secret@localhost:5432/alloydb?sslmode=disable")

generate-openapi:
	cd backend && bru2openapi -d ./docs -o ./internal/shared/apidocs/openapi.yaml -exclude tmp,examples
	perl -0pi -e 's#/\{\{[^}]+\}\}/api#/api#g' backend/internal/shared/apidocs/openapi.yaml

migrate-create:
	cd backend && goose -dir internal/shared/database/migrations create $(name) sql

migrate-up:
	cd backend && GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir internal/shared/database/migrations up

migrate-down:
	cd backend && GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir internal/shared/database/migrations down

migrate-status:
	cd backend && GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir internal/shared/database/migrations status

run-backend-local:
	docker compose -f docker-compose.yml -f docker-compose-dev.yml --profile alloy-api up --build

run-frontend-local:
	docker compose -f docker-compose.yml -f docker-compose-dev.yml --profile alloy-ui up --build

run-all-local:
	docker compose -f docker-compose.yml -f docker-compose-dev.yml --profile alloy-api --profile alloy-ui up --build

run-all-staging:
	docker compose -f docker-compose.yml -f docker-compose-staging.yml up -d --build
