# Backend

This directory contains the backend/server-side code for the Alloy project.

## Setup Guide

**Note:** All commands can be run from the project root directory. There is no need to `cd` into the `backend` folder.

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL 16 or higher
- Redis 7 or higher
- Docker and Docker Compose (for containerized setup)
- direnv (for environment variable management)

### Step 1: Install direnv

direnv allows automatic loading of environment variables from `.envrc` files.

**macOS:**

```bash
brew install direnv
```

**Linux:**

```bash
# For Ubuntu/Debian
sudo apt-get install direnv

# For Fedora
sudo dnf install direnv
```

**Add to your shell configuration** (add to `~/.zshrc` or `~/.bashrc`):

```bash
eval "$(direnv hook zsh)"  # For zsh
# or
eval "$(direnv hook bash)" # For bash
```

Reload your shell:

```bash
source ~/.zshrc  # or source ~/.bashrc
```

### Step 2: Set Up .envrc File

Create the `.envrc` file from the example:

```bash
# From project root
cp .envrc.example .envrc
```

**Important:** You'll need to update the `GOOSE_DBSTRING` in `.envrc` to match your PostgreSQL credentials after setting up your `.env` file (see Step 4).

### Step 3: Allow direnv

Navigate to the project root and allow direnv to load the `.envrc` file:

```bash
# From project root
direnv allow
```

This will load the environment variables defined in `.envrc` (including `GOOSE_DRIVER` and `GOOSE_DBSTRING`).

### Step 4: Install Goose (Database Migration Tool)

Goose is used for database migrations. Install it globally:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Verify installation:

```bash
goose --version
```

### Step 5: Set Up Environment Variables

Create a `.env` file in the `backend` directory:

```bash
# From project root
cp backend/.env.example backend/.env  # If you have an example file
# Or create backend/.env manually
```

**Important:** After setting up your PostgreSQL credentials in `backend/.env`, make sure to update the `GOOSE_DBSTRING` in `.envrc` to match. The `GOOSE_DBSTRING` must correspond with the PostgreSQL credentials in your `.env` file.

For example, if your `.env` has:

```env
POSTGRES_USER=alloy
POSTGRES_PASSWORD=alloy_secret
POSTGRES_DB=alloydb
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

Then your `.envrc` should have:

```bash
export GOOSE_DBSTRING="postgres://alloy:alloy_secret@localhost:5432/alloydb?sslmode=disable"
```

Required environment variables (add to `backend/.env`):

```env
# Application
APP_MODE=dev
PORT=8082

# PostgreSQL (use 'postgres' as hostname for Docker, 'localhost' for local)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=alloy
POSTGRES_PASSWORD=alloy_secret
POSTGRES_DB=alloydb
POSTGRES_SCHEMA=public
POSTGRES_DSN=host=localhost user=alloy password=alloy_secret dbname=alloydb port=5432 sslmode=disable search_path=public

# Redis (use 'redis' as hostname for Docker, 'localhost' for local)
REDIS_ADDR=localhost:6379

# Cassandra (use 'cassandra' as hostname for Docker, 'localhost' for local)
CASSANDRA_HOST=localhost
CASSANDRA_PORT=9042
CASSANDRA_KEYSPACE=chat

# JWT
JWT_SECRET=your-secret-key-here-change-in-production
JWT_EXPIRY=1
JWT_REFRESH_EXPIRY=7

# CORS
ORIGINS=http://localhost:3000,http://localhost:8080
```

**Note:** When using Docker Compose, use service names (`postgres`, `redis`, `cassandra`) instead of `localhost` for hostnames.

### Step 6: Set Up Database

#### Option A: Using Docker Compose (Recommended)

Start all services including databases:

```bash
# From project root
make run-backend-stack-only
# or
docker compose --profile alloy up
```

This will start:

- PostgreSQL on port 5432
- Redis on port 6379
- Cassandra on port 9042
- Run database migrations automatically
- Start the backend with hot reload

#### Option B: Local Database Setup

1. **Install and start PostgreSQL:**

   ```bash
   # macOS
   brew install postgresql@16
   brew services start postgresql@16

   # Create database and user
   createdb alloydb
   createuser alloy
   ```

2. **Install and start Redis:**

   ```bash
   # macOS
   brew install redis
   brew services start redis
   ```

3. **Install and start Cassandra:**
   ```bash
   # macOS
   brew install cassandra
   brew services start cassandra
   ```

### Step 7: Run Database Migrations

After setting up the database, run migrations:

```bash
# From project root
make migrate-up
```

Check migration status:

```bash
make migrate-status
```

Rollback last migration:

```bash
make migrate-down
```

### Step 8: Install Dependencies

Install Go dependencies:

```bash
# From project root
cd backend && go mod download && cd ..
```

### Step 9: Run the Backend

#### Option A: Using Make (Recommended)

**Development mode (with hot reload using Air):**

```bash
# From project root
make run-backend-stack-only
```

**Local development (without Docker):**

```bash
# From project root
make run
```

**Build and run:**

```bash
# From project root
make build
make start
```

#### Option B: Using Docker Compose

```bash
# From project root
docker compose --profile alloy up --build
```

#### Option C: Direct Go Command

```bash
# From project root
cd backend && go run ./cmd
```

The backend will be available at `http://localhost:8082` (or the port specified in `PORT` env var).

## Available Make Commands

All commands are run from the project root:

- `make build` - Build the backend binary
- `make run` - Run the backend in development mode
- `make start` - Build and start the backend
- `make clean` - Remove build artifacts
- `make migrate-create name=<name>` - Create a new migration file
- `make migrate-up` - Run pending migrations
- `make migrate-down` - Rollback last migration
- `make migrate-status` - Check migration status
- `make test-backend` - Run backend tests
- `make run-backend-stack-only` - Start backend stack with Docker Compose
- `make run-all-stacks` - Start all stacks (backend + frontend)

## Project Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── app/                 # Application bootstrap
│   ├── modules/             # Feature modules (auth, users, messaging)
│   └── shared/              # Shared utilities (config, database, router)
├── email_templates/         # HTML email templates
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── .air.toml                # Air hot reload configuration
└── Dockerfile               # Docker image definition
```

## Development

### Hot Reload

The backend uses [Air](https://github.com/air-verse/air) for hot reload during development. When running with Docker Compose, changes to Go files, templates, and SQL files will automatically trigger a rebuild and restart.

### Creating Migrations

Create a new migration:

```bash
# From project root
make migrate-create name=add_new_table
```

This creates a new migration file in `backend/internal/shared/database/migrations/`.

### Testing

Run tests:

```bash
# From project root
make test-backend
```

## Troubleshooting

### Migration Issues

If migrations fail:

1. Check that PostgreSQL is running
2. **Important:** Verify that `GOOSE_DBSTRING` in `.envrc` matches your PostgreSQL credentials in `backend/.env`. The connection string must correspond with:
   - `POSTGRES_USER`
   - `POSTGRES_PASSWORD`
   - `POSTGRES_DB`
   - `POSTGRES_HOST`
   - `POSTGRES_PORT`
3. Ensure the database and user exist
4. Check migration status: `make migrate-status`
5. After updating `.envrc`, run `direnv allow` to reload the environment variables

### Connection Issues

- **PostgreSQL:** Ensure the database is running and credentials match
- **Redis:** Check that Redis is accessible on the configured port
- **Cassandra:** Verify Cassandra is running and the keyspace exists

### Port Conflicts

If port 8082 is already in use:

1. Change `PORT` in `backend/.env`
2. Update port mapping in `docker-compose.yml` if using Docker

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Goose Migration Tool](https://github.com/pressly/goose)
- [Air Hot Reload](https://github.com/air-verse/air)
- [Fiber Web Framework](https://docs.gofiber.io/)
