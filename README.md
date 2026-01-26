# 👤 User Management

## 🏁 Getting Started
Quick setup guide for local development.

### 1. 🔐 Environment Variables
Before running the application, create a `.env` file from the example:

```bash
cp .env.example .env
```

### 2. 🏗️ Launch Infrastructure
Spin up core dependencies (Postgres, RabbitMQ, etc.) via Docker:

```bash
make infra_up
```

### 3. 🛠️ Initialize & Run
Before starting, ensure your schema is migrated and code is generated:

```bash
# Apply migrations & generate type-safe SQL code
make migrate_up
make sqlc

# Option A: Dev Mode (Live Reload)
air

# Option B: Standard Mode
make dev
```

### 🐳 Full Stack Execution
To run the entire application stack including the API and all dependencies:

```bash
# Development
make dev_up

# Production
make prod_up
```

## 🗄️ Database Migration

This project uses golang-migrate (via `Makefile`) to manage database schema changes in a consistent and reliable way.

```bash
# Create a new migration
make create_migration NAME=users

# Apply all pending migrations
make migrate_up

# Rollback the most recent migration (1 step)
make migrate_down

# Force database to a specific version (useful when migration state is dirty)
make migrate_force VERSION=1

# Migrate database to a specific version
make migrate_goto VERSION=1

# Drop all tables and schema from the database (USE WITH CAUTION)
make migrate_drop
```

## 🧩 SQL Code Generation (sqlc)

Uses sqlc for type-safe Go code generation from SQL.

Installation:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

To generate Go models and query functions from the SQL files, run:
```bash
make sqlc
# or
sqlc generate
```
> Generated code is based on sqlc.yaml and stays in sync with your schema.
