# iam box

**Single entrypoint for application access rights management**

Lightweight, self-contained authorization service with REST API and CLI interfaces.

## Features

- ✅ Role-Based Access Control (RBAC)
- ✅ Wildcard permissions (e.g., "user can delete ANY invoice")
- ✅ Audit logging of every decision
- ✅ REST API + CLI interface (same binary)
- ✅ PostgreSQL storage
- ✅ Docker-ready

## Quick Start

### 1. Run with Docker Compose

```bash
git clone https://github.com/awayneff/iam-box
cd iam-box
docker-compose up -d
```

### 2. Or build and run binary

```bash
cd iam-box/go
go build -o iam-box main.go
./iam-box server --port 8080
```

## Usage

### Start API Server

```bash
./iam-box server --port 8080 --db-host localhost
```

### CLI Commands

```bash
# Grant permission
./iam-box grant alice delete invoice 123

# Check permission
./iam-box check alice delete invoice 123

# Revoke permission
./iam-box revoke alice delete invoice 123

# List all permissions with pagination
./iam-box list

# List user permissions
./iam-box list-user alice

# [not implemented] View audit log
./iam-box audit alice
```

### [not implemented] Interactive Mode

```bash
./iam-box
iam> grant alice delete invoice 123
✅ Granted
iam> check alice delete invoice 123
✅ ALLOWED
iam> exit
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/permissions/grant` | Grant permission |
| POST | `/api/v1/permissions/can` | Check permission |
| DELETE | `/api/v1/permissions/revoke` | Revoke permission |
| GET | `/api/v1/permissions/{user_id}` | List user permissions |
| GET | `/api/v1/decisions/` | List audit log |
| GET | `/_core/health` | Health check |

## Configuration

### Flags

```bash
--db-host     Database host (default: localhost)
--db-user     Database user (default: iam_user)
--db-password Database password (default: iam_password)
--db-name     Database name (default: iam_db)
--db-port     Database port (default: 5432)
--port        API port (default: 8080)
```

### [not implemented] Environment Variables

```bash
export IAM_DB_HOST=postgres.example.com
export IAM_DB_PASSWORD=secret
export IAM_PORT=9090
./iam-box server
```

## Database Schema

```sql
-- Permissions table
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    granted_at TIMESTAMP DEFAULT NOW(),
    created_by TEXT
);

-- Audit log
CREATE TABLE decisions (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    allowed BOOLEAN NOT NULL,
    reason TEXT,
    timestamp TIMESTAMP DEFAULT NOW()
);
```

## Examples

### Grant and Check via API

```bash
# Grant
curl -X POST http://localhost:8080/api/v1/permissions/grant \
  -H "Content-Type: application/json" \
  -d '{"user_id":"alice","action":"delete","resource_type":"invoice","resource_id":"123"}'

# Check
curl -X POST http://localhost:8080/api/v1/permissions/can \
  -H "Content-Type: application/json" \
  -d '{"user_id":"alice","action":"delete","resource_type":"invoice","resource_id":"123"}'

# Response: {"allowed":true}
```

### Wildcard Permissions

```bash
# Grant permission for ALL invoices
./iam-box grant alice delete invoice NULL
# or
./iam-box grant alice delete invoice

# Check specific invoice (returns true)
./iam-box check alice delete invoice 123
```

## Tech Stack

- **Go** - Core language
- **Cobra** - CLI framework
- **Chi** - HTTP router
- **GORM** - ORM
- **PostgreSQL** - Database

## Note

This project is in early stage of development