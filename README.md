# Sprint Backlog - Backend

REST API backend for Sprint Backlog Management System. Built with Go and Gin framework.

## Tech Stack

| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go | 1.24 |
| Framework | Gin | 1.11 |
| Database | PostgreSQL | 15+ |
| ORM | GORM | 1.31 |
| Authentication | JWT + Google OAuth | - |
| API Documentation | Swagger (swaggo) | 1.16 |
| UUID | google/uuid | 1.6 |

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database connection & migrations
│   ├── dto/
│   │   ├── request/          # Request DTOs
│   │   └── response/         # Response DTOs
│   ├── handler/              # HTTP handlers (controllers)
│   ├── middleware/           # Auth, CORS, Logger middleware
│   ├── models/               # GORM models (entities)
│   ├── repository/           # Data access layer
│   ├── router/               # Route definitions
│   ├── service/              # Business logic layer
│   └── utils/                # Utility functions
├── docs/                     # Swagger documentation (auto-generated)
├── pkg/
│   └── constants/            # Application constants
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 15+
- Make (optional)

### Environment Variables

Create a `.env` file in the backend root folder:

```env
# Server
PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sprint_backlog

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### Running the Application

#### Using Make

```bash
# Download dependencies
make deps

# Run the application
make run

# Build the application
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Generate Swagger docs
make swagger

# Clean build artifacts
make clean
```

#### Using Go directly

```bash
# Download dependencies
go mod download

# Run the application
go run cmd/api/main.go

# Build the application
go build -o bin/api cmd/api/main.go

# Run tests
go test -v ./...
```

#### Using Docker

```bash
# Build Docker image
make docker-build
# or
docker build -t sprint-backlog-api .

# Run with Docker Compose (from root directory)
docker-compose up -d
```

## API Documentation

Swagger UI is available at `/swagger/index.html` when the application is running.

### Main Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/google/verify` | Verify Google OAuth code |
| GET | `/api/auth/me` | Get current user info |
| GET | `/api/projects` | Get all projects |
| POST | `/api/projects` | Create a project |
| GET | `/api/projects/:id` | Get project by ID |
| PUT | `/api/projects/:id` | Update project |
| DELETE | `/api/projects/:id` | Delete project |
| GET | `/api/items` | Get backlog items |
| POST | `/api/items` | Create backlog item |
| GET | `/api/items/:id` | Get item by ID |
| PUT | `/api/items/:id` | Update item |
| DELETE | `/api/items/:id` | Delete item |
| GET | `/api/sprints` | Get sprints |
| POST | `/api/sprints` | Create sprint |
| GET | `/api/sprints/:id` | Get sprint by ID |
| PUT | `/api/sprints/:id` | Update sprint |
| POST | `/api/sprints/:id/start` | Start a sprint |
| POST | `/api/sprints/:id/complete` | Complete a sprint |
| POST | `/api/sprints/:id/cancel` | Cancel a sprint |

## Models

### User
- ID, Email, Name, Picture, GoogleID
- Timestamps (CreatedAt, UpdatedAt)

### Project
- ID, Name, Description, Key
- OwnerID (User reference)
- Timestamps

### BacklogItem
- ID, Title, Description, Type, Status, Priority
- StoryPoints, Labels
- ProjectID, SprintID, AssigneeID
- Timestamps

### Sprint
- ID, Name, Goal, Status
- StartDate, EndDate
- ProjectID
- Timestamps

## Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
# Output: coverage.html
```

## Development

### Generate Swagger Documentation

```bash
# Install swag CLI (if not installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
make swagger
# or
swag init -g cmd/api/main.go -o docs
```

### Code Structure

The application follows clean architecture pattern:

1. **Handler** - Receives HTTP requests, validates input, calls service
2. **Service** - Business logic, orchestration
3. **Repository** - Data access, database operations
4. **Model** - Entity definitions (GORM models)
5. **DTO** - Data Transfer Objects for request/response
