# Go Task API

A RESTful Task Manager API built with Go, PostgreSQL and Clean Architecture principles.

## Tech Stack

- **Go 1.26** — main language
- **PostgreSQL 16** — database
- **chi** — HTTP router
- **pgx** — PostgreSQL driver
- **golang-jwt** — JWT authentication
- **golang-migrate** — database migrations
- **slog** — structured logging
- **Docker** — local database


## Architecture
```text
go-task-api/
├── cmd/api/          # Entry point
├── internal/
│   ├── domain/       # Models & interfaces
│   ├── repository/   # SQL queries
│   ├── usecase/      # Business logic
│   └── handler/      # HTTP handlers & middleware
└── migrations/       # SQL migrations
```

## Prerequisites

- Go 1.26+
- Docker & Docker Compose

## Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/Cxxlkid/go-task-api.git
cd go-task-api
```

### 2. Configure environment
```bash
cp .env.example .env
# Edit .env with your values
```

### 3. Start the database
```bash
docker-compose up -d
```

### 4. Run migrations
```bash
docker run --rm -v "${PWD}/migrations:/migrations" --network go-task-api_default migrate/migrate \
  -path=/migrations \
  -database "postgres://postgres:postgres@go-task-api-db:5432/taskapi?sslmode=disable" up
```

### 5. Start the server
```bash
go run cmd/api/main.go
```

Server runs on `http://localhost:8080`

## API Endpoints

### Auth
| Method | Endpoint | Description | Auth (Jwt Needed) |
|--------|----------|-------------|------|
| POST | `/auth/register` | Register a new user | ❌ |
| POST | `/auth/login` | Login and get JWT token | ❌ |

### Users
| Method | Endpoint | Description | Auth (Jwt Needed) |
|--------|----------|-------------|------|
| GET | `/users/me` | Get current user | ✅ |

### Tasks
| Method | Endpoint | Description | Auth (Jwt Needed) |
|--------|----------|-------------|------|
| POST | `/tasks` | Create a task | ✅ |
| GET | `/tasks` | List tasks | ✅ |
| GET | `/tasks/{id}` | Get a task | ✅ |
| PUT | `/tasks/{id}` | Update a task | ✅ |
| DELETE | `/tasks/{id}` | Delete a task | ✅ |

### Health
| Method | Endpoint | Description | Auth (Jwt Needed) |
|--------|----------|-------------|------|
| GET | `/health` | Health check | ❌ |

### Query Parameters for GET /tasks
| Param | Description | Example |
|-------|-------------|---------|
| `status` | Filter by status | `?status=todo` |
| `assigned_to` | Filter by user ID | `?assigned_to=uuid` |
| `sort` | Sort by field | `?sort=due_date` |

## Status Values

| Value | Description |
|-------|-------------|
| `todo` | Task not started |
| `in_progress` | Task in progress |
| `done` | Task completed |