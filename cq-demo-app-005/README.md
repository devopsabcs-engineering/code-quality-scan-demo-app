# cq-demo-app-005 — Go / net/http Task & Project Management API

A deliberately flawed Go REST API for task and project management, built with `net/http`. This demo app contains **intentional code quality violations** for use with the Code Quality Scanner workshop.

## Intentional Violations

| Category | File(s) | Description |
|----------|---------|-------------|
| **High Complexity** | `internal/handlers/task_handler.go`, `internal/services/task_service.go` | `HandleProcessTask()` and `ProcessTask()` have CCN > 15 with deeply nested switch/case + if/else chains |
| **High Complexity** | `internal/handlers/project_handler.go`, `internal/services/project_service.go` | `HandleUpdateProjectStatus()` and `ProcessProjectStatus()` duplicate the same complex pattern |
| **Code Duplication** | `internal/handlers/task_handler.go` ↔ `project_handler.go` | Filtering logic, validation, and status transition patterns are duplicated |
| **Code Duplication** | `internal/services/task_service.go` ↔ `project_service.go` | State machine logic duplicated across services |
| **Code Duplication** | `internal/utils/formatters.go` | `FormatTaskSummary` and `FormatProjectSummary` duplicate label formatting |
| **Lint: Unused Vars** | `internal/utils/validators.go` | `debugMode`, `maxRetries`, `defaultTimeout`, `internalVersion` declared but unused |
| **Lint: Unused Funcs** | `internal/utils/validators.go` | `validateLength()`, `validateBudget()` defined but never called |
| **Lint: Error Handling** | `internal/utils/validators.go` | `json.Unmarshal()` and `regexp.Compile()` errors ignored |
| **Lint: Unreachable Code** | `internal/utils/validators.go` | Code after `return` in `ValidateTaskTitle()` and `ValidateProjectName()` |
| **Lint: Error Handling** | `internal/utils/formatters.go` | Errors returned without wrapping (`%w` not used) |
| **Lint: Magic Strings** | `internal/utils/formatters.go` | Status and priority labels hardcoded throughout |
| **Low Test Coverage** | `tests/handlers/task_handler_test.go` | Only GET endpoints tested; services, utils, and POST/PUT/DELETE untested |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Health check with service info |
| `GET` | `/health` | Simple health probe |
| `GET` | `/api/tasks` | List all tasks (optional `?status=` and `?priority=` filters) |
| `POST` | `/api/tasks` | Create a new task |
| `GET` | `/api/tasks/{id}` | Get task by ID |
| `PUT` | `/api/tasks/{id}` | Update a task |
| `DELETE` | `/api/tasks/{id}` | Delete a task |
| `POST` | `/api/tasks/process` | Process task state transition |
| `GET` | `/api/projects` | List all projects |
| `POST` | `/api/projects` | Create a new project |
| `GET` | `/api/projects/{id}` | Get project by ID |
| `PUT` | `/api/projects/{id}` | Update a project |
| `DELETE` | `/api/projects/{id}` | Delete a project |
| `POST` | `/api/projects/status` | Process project state transition |

## Run Locally

### Using Docker (recommended)

```bash
docker build -t cq-demo-app-005 .
docker run -p 8080:8080 cq-demo-app-005
```

### Using Go directly

```bash
go run ./cmd/server
```

### Verify it works

```bash
curl http://localhost:8080/
curl http://localhost:8080/api/tasks
curl http://localhost:8080/api/projects
```

## Run Tests

```bash
go test ./tests/... -v
```

### Run with coverage

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Run Linting

```bash
golangci-lint run ./...
```

## Technology Stack

- **Language**: Go 1.22
- **Framework**: `net/http` (standard library)
- **Linter**: golangci-lint
- **Container**: Multi-stage Docker build with Alpine

## Project Structure

```
cq-demo-app-005/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── task_handler.go         # Task HTTP handlers (COMPLEXITY)
│   │   └── project_handler.go      # Project HTTP handlers (DUPLICATION)
│   ├── services/
│   │   ├── task_service.go         # Task business logic (COMPLEXITY + DUPLICATION)
│   │   └── project_service.go      # Project business logic (DUPLICATION)
│   ├── models/
│   │   ├── task.go                 # Task data models
│   │   └── project.go             # Project data models
│   └── utils/
│       ├── validators.go           # Validation utilities (LINT VIOLATIONS)
│       └── formatters.go           # Formatting utilities (LINT VIOLATIONS)
├── tests/handlers/
│   └── task_handler_test.go        # Minimal tests (LOW COVERAGE)
├── infra/main.bicep                # Azure infrastructure (ACR + Web App for Containers)
├── .golangci.yml                   # Linter configuration
├── Dockerfile                      # Multi-stage Go build
├── go.mod                          # Go module definition
└── go.sum                          # Dependency checksums
```
