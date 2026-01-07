# User Management

## 🚀 Development Setup

### 🔄 Using Air for Live Reload

This project utilizes **Air** for automatic live reloading during development, significantly speeding up the development feedback loop.

1. **Install Air:**

```bash
go install github.com/air-verse/air@latest
```
> **Note:** Make sure your `$GOPATH/bin` directory is included in your system's `$PATH` environment variable so the `air` command is accessible globally.

2. **Run the application with Air:**
*If this is your first time, initialize the configuration file:*
```bash
air init
```

Then run:
```bash
air
```
>Your application will compile and start. Air will automatically monitor your Go source files and restart the application whenever a change is detected and saved.

## 🔐 Environment Variables

Before running the application, create a `.env` file from the example:

```bash
cp .env.example .env
```

## 🗄️ Database Migration

This project uses golang-migrate (via `Makefile`) to manage database schema changes in a consistent and reliable way.

```bash
# Create a new migration
make create_migration NAME=users

# Apply all pending migrations
make migrate_up

# Rollback the most recent migration
make migrate_down

# Force database to a specific version (useful when migration state is dirty)
make migrate_force VERSION=1

# Migrate database to a specific version
make migrate_goto VERSION=1

# Drop all tables and schema from the database (USE WITH CAUTION)
make migrate_drop
```

## 🧩 SQL Code Generation (sqlc)

This project uses sqlc to generate type-safe Go code from SQL queries, helping reduce boilerplate and prevent runtime SQL errors.

Before running the application, make sure `sqlc` is installed:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

To generate Go models and query functions from the SQL files, run:
```bash
# Option 1: Use make to generate the code:
make sqlc

# Option 2: Alternatively, run the command directly:
sqlc generate
```
> This command reads the SQL definitions and configuration in sqlc.yaml and generates the corresponding Go code, keeping your database layer strongly typed and in sync with your schema.
