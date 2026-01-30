# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Alloy is a full-stack internal software solution with multi-tenancy support, real-time messaging, and RBAC.

**Tech Stack:**
- Backend: Go 1.25.1 with Fiber v2, GORM, PostgreSQL, Redis, Cassandra
- Frontend: Next.js 16, React 19, TypeScript, TailwindCSS 4, Redux Toolkit
- Mobile: Flutter (in `mobile/`)

## Development Commands

### Backend
```bash
make run                    # Run with Air hot reload
make build                  # Build Go binary
make test-backend           # Run tests
```

### Frontend
```bash
pnpm install                # Install dependencies (run from frontend/)
pnpm dev                    # Development server
pnpm build                  # Production build
pnpm lint                   # Run ESLint
```

### Database Migrations
```bash
make migrate-create name=<name>  # Create new migration
make migrate-up                  # Apply migrations
make migrate-down                # Rollback last migration
make migrate-status              # Check migration status
```

### Docker (Full Stack)
```bash
make run-backend-stack-only     # Backend + PostgreSQL + Redis + Cassandra
make run-all-stacks             # Full stack including frontend
```

### Environment Setup
Requires `direnv` for database connection strings. The `.envrc` file sets `GOOSE_DBSTRING` and `GOOSE_DRIVER` for migrations.

## Architecture

### Backend Module Pattern
Each feature module in `backend/internal/modules/` follows this structure:
```
module/
├── module.go       # Factory: initializes Repository → Service → Handler
├── handler.go      # HTTP handlers (Fiber routes)
├── service.go      # Business logic
├── repository.go   # Database operations (GORM)
├── dto.go          # Request/response DTOs
└── errors.go       # Module-specific errors
```

All modules implement `IHandler` interface and are registered in `internal/app/bootstrap.go`. Cross-module communication uses the `ModuleServices` registry.

### Database Strategy
- **PostgreSQL**: Relational data (users, auth, organizations, invitations)
- **Cassandra**: High-volume time-series data (chat messages in `chat` keyspace)
- **Redis**: Session cache, WebSocket connection tracking

### Frontend Monorepo Structure
The frontend uses pnpm workspaces with packages in `frontend/packages/`:
- `ui/` - Radix UI component library (50+ components)
- `store/` - Redux Toolkit store and slices
- `plugins/` - Plugin registry system
- `permissions/` - RBAC permission helpers
- `types/` - Shared TypeScript types
- `feature-flags/` - Runtime feature management

### Route Groups (Next.js App Router)
- `(public)/` - No auth required (login)
- `(system)/` - Authenticated, no tenant context (select-tenant)
- `(tenant)/` - Full app with tenant context

### State Management
Redux slices handle global state: `auth`, `tenant`, `users`, `chat`, `permissions`, `features`, `ui`. The middleware in `frontend/middleware.ts` handles auth routing.

## Key Patterns

### Adding a New Backend Module
1. Create module directory in `backend/internal/modules/`
2. Implement: `module.go`, `handler.go`, `service.go`, `repository.go`, `dto.go`
3. Handler must implement `IHandler` interface with `RegisterRoutes(router fiber.Router)`
4. Register in `bootstrap.go`

### API Routes
All API routes are prefixed with `/api`. Auth middleware is applied via `middlewares.AuthMiddleware()`.

### Real-time Communication
WebSocket connections use Socket.io. Backend socket management is in `backend/internal/shared/socket/`. Frontend uses `socket.io-client`.

### Authentication Flow
- Magic link authentication via email (Mailjet)
- JWT tokens with refresh mechanism
- Middleware extracts user from token and sets `c.Locals("user", claims)`

## Code Style

- Do not add unnecessary comments to code. Code should be self-documenting.

## Important Files

- `backend/cmd/main.go` - Application entry point
- `backend/internal/app/bootstrap.go` - Module initialization
- `backend/internal/shared/config/config.go` - Environment configuration
- `frontend/middleware.ts` - Auth routing logic
- `frontend/packages/store/slices/` - Redux state definitions
- `docker-compose.yaml` - Service definitions
